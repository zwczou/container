package container

import (
	"sync"
	"time"

	"github.com/codegangsta/inject"
	"github.com/cskr/pubsub"
	"github.com/sirupsen/logrus"
)

type Provider interface {
	Name() string
	Load(*Container) error
	Exit()
}

// 自带依赖注射功能
// 自带pub/sub用于在扩展之间通信
type Container struct {
	sync.RWMutex
	inject.Injector
	*pubsub.PubSub
	names      []string
	extensions map[string]Provider
}

func New() *Container {
	c := &Container{
		Injector:   inject.New(),
		PubSub:     pubsub.New(0),
		extensions: make(map[string]Provider),
	}
	return c
}

// 前置注册扩展
func (c *Container) Pre(ext Provider) *Container {
	c.Lock()
	if _, ok := c.extensions[ext.Name()]; !ok {
		c.extensions[ext.Name()] = ext
		c.names = append([]string{ext.Name()}, c.names...)
	}
	c.Unlock()
	return c
}

// 注册扩展
func (c *Container) Use(ext Provider) *Container {
	c.Lock()
	if _, ok := c.extensions[ext.Name()]; !ok {
		c.extensions[ext.Name()] = ext
		c.names = append(c.names, ext.Name())
	}
	c.Unlock()
	return c
}

// 加载扩展
func (c *Container) Load() error {
	c.RLock()
	defer c.RUnlock()

	for _, name := range c.names {
		start := time.Now()
		ext := c.extensions[name]
		err := ext.Load(c)
		if err != nil {
			return err
		}
		logrus.WithField("extension", name).WithField("spent", time.Since(start)).Info("extension loading")
	}
	return nil
}

// 注销扩展
func (c *Container) Exit() {
	c.RLock()
	for _, name := range c.names {
		start := time.Now()
		ext := c.extensions[name]
		ext.Exit()
		logrus.WithField("extension", name).WithField("spent", time.Since(start)).Info("extension exiting")
	}
	c.RUnlock()

	if c.PubSub != nil {
		c.PubSub.Shutdown()
	}
}
