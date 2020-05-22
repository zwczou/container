package container

import (
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/sirupsen/logrus"
)

type Provider interface {
	Name() string
	Load(*Container) error
	Exit()
}

// 自带pub/sub用于在扩展之间通信
type Container struct {
	sync.RWMutex
	*pubsub.PubSub
	Metadata
	existNames map[string]bool
	extensions []Provider
}

func New() *Container {
	c := &Container{
		PubSub:     pubsub.New(0),
		existNames: make(map[string]bool),
	}
	return c
}

// 前置注册扩展
func (c *Container) Front(exts ...Provider) *Container {
	c.Lock()
	for _, ext := range exts {
		if _, ok := c.existNames[ext.Name()]; !ok {
			c.existNames[ext.Name()] = true
			c.extensions = append([]Provider{ext}, c.extensions...)
		}
	}
	c.Unlock()
	return c
}

// 注册扩展
func (c *Container) Push(exts ...Provider) *Container {
	c.Lock()
	for _, ext := range exts {
		if _, ok := c.existNames[ext.Name()]; !ok {
			c.existNames[ext.Name()] = true
			c.extensions = append(c.extensions, ext)
		}
	}
	c.Unlock()
	return c
}

// 加载扩展
func (c *Container) Load() error {
	c.RLock()
	defer c.RUnlock()

	for _, ext := range c.extensions {
		start := time.Now()
		err := ext.Load(c)
		if err != nil {
			return err
		}
		logrus.WithField("extension", ext.Name()).WithField("spent", time.Since(start)).Info("extension loading")
	}
	return nil
}

// 导出扩展
func (c *Container) All() []Provider {
	c.RLock()
	defer c.RUnlock()
	return c.extensions
}

// 注销扩展
func (c *Container) Exit() {
	c.RLock()

	for _, ext := range c.extensions {
		start := time.Now()
		ext.Exit()
		logrus.WithField("extension", ext.Name()).WithField("spent", time.Since(start)).Info("extension exiting")
	}
	c.RUnlock()

	if c.PubSub != nil {
		c.PubSub.Shutdown()
	}
}
