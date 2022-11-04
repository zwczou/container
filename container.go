package container

import (
	"reflect"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/sirupsen/logrus"
)

type defaultContainer struct {
	lock       sync.RWMutex
	valueLock  sync.RWMutex
	pubSub     *pubsub.PubSub
	values     map[reflect.Type]reflect.Value
	existNames map[string]bool
	extensions []Provider
	queues     []Queuer
	exitChan   chan struct{}
}

func New() Container {
	return &defaultContainer{
		pubSub:     pubsub.New(0),
		values:     make(map[reflect.Type]reflect.Value),
		existNames: make(map[string]bool),
		exitChan:   make(chan struct{}),
	}
}

func (c *defaultContainer) Add(ext Provider) {
	c.lock.Lock()
	if _, ok := c.existNames[ext.Name()]; !ok {
		c.existNames[ext.Name()] = true
		c.extensions = append(c.extensions, ext)
	}
	c.lock.Unlock()
}

func (c *defaultContainer) Load() error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, ext := range c.extensions {
		start := time.Now()
		err := ext.Load(c)
		if err != nil {
			return err
		}

		fields := logrus.Fields{
			"extension": ext.Name(),
			"spent":     time.Since(start),
		}
		logrus.WithFields(fields).Info("extension load")
	}
	return nil
}

func (c *defaultContainer) Exit() {
	c.lock.RLock()
	for i := len(c.extensions) - 1; i >= 0; i-- {
		ext := c.extensions[i]
		start := time.Now()
		ext.Exit()

		fields := logrus.Fields{
			"extension": ext.Name(),
			"spent":     time.Since(start),
		}
		logrus.WithFields(fields).Info("extension exit")
	}
	c.lock.RUnlock()

	close(c.exitChan)

	if c.pubSub != nil {
		c.pubSub.Shutdown()
	}

	for _, queue := range c.queues {
		queue.Unsub()
	}
}

func (c *defaultContainer) Pub(name string, params ...any) {
	c.pubSub.Pub(params, name)
}

func (c *defaultContainer) TryPub(name string, params ...any) {
	c.pubSub.TryPub(params, name)
}
