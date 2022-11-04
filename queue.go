package container

import (
	"reflect"

	"github.com/korovkin/limiter"
	"github.com/sirupsen/logrus"
)

type Queuer interface {
	Sub()
	Unsub()
	Listen(fn any)
}

type queueConfig struct {
	concurrency int
	size        int
}

func (c *queueConfig) Default() {
	c.concurrency = 1
	c.size = 200
}

type QueueOption func(*queueConfig)

func WithConcurrency(concurrency int) QueueOption {
	return func(cfg *queueConfig) {
		cfg.concurrency = concurrency
	}
}

func WithSize(size int) QueueOption {
	return func(cfg *queueConfig) {
		cfg.size = size
	}
}

type defaultQueue struct {
	name        string
	app         *defaultContainer
	concurrency int
	msgChan     chan interface{}
}

func (c *defaultContainer) Queue(name string, opts ...QueueOption) Queuer {
	var cfg queueConfig
	cfg.Default()
	for _, opt := range opts {
		opt(&cfg)
	}

	queue := &defaultQueue{
		app:         c,
		name:        name,
		concurrency: cfg.concurrency,
		msgChan:     make(chan interface{}, cfg.size),
	}

	c.valueLock.Lock()
	c.queues = append(c.queues, queue)
	c.valueLock.Unlock()
	return queue
}

func (q *defaultQueue) Sub() {
	q.app.pubSub.AddSub(q.msgChan, q.name)
}

func (q *defaultQueue) Unsub() {
	q.app.pubSub.Unsub(q.msgChan, q.name)
}

func (q *defaultQueue) Listen(fn interface{}) {
	q.Sub()
	logrus.WithField("queue", q.name).Info("queue init")

	funcValue := reflect.ValueOf(fn)
	funcType := funcValue.Type()
	queueType := reflect.TypeOf(q)
	var limit *limiter.ConcurrencyLimiter
	if q.concurrency > 1 {
		limit = limiter.NewConcurrencyLimiter(q.concurrency)
	}

	procFunc := func(vs ...interface{}) func() {
		return func() {
			var values []reflect.Value
			for _, val := range vs {
				values = append(values, reflect.ValueOf(val))
			}

			// 如果需要传入Queue自动填充
			if funcType.NumIn() > 0 && funcType.In(0) == queueType {
				if len(values) == 0 || values[0].Type() != queueType {
					values = append([]reflect.Value{reflect.ValueOf(q)}, values...)
				}
			}

			funcValue.Call(values)
		}
	}

	for {
		select {
		case v, ok := <-q.msgChan:
			if !ok {
				// 如果队列退出，直接返回
				logrus.WithField("queue", q.name).Info("queue exiting")
				return
			}

			params := v.([]interface{})
			if q.concurrency > 1 {
				_, _ = limit.Execute(procFunc(params...))
			} else {
				procFunc(params...)()
			}
		case <-q.app.exitChan:
			logrus.WithField("queue", q.name).Info("app exiting")
			go q.Unsub()
			if q.concurrency > 1 {
				limit.WaitAndClose()
			}
			return
		}
	}
}
