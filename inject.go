package container

import (
	"fmt"
	"reflect"
)

func (c *defaultContainer) Set(vals ...any) {
	c.valueLock.Lock()
	for _, val := range vals {
		c.values[reflect.TypeOf(val)] = reflect.ValueOf(val)
	}
	c.valueLock.Unlock()
}

func (c *defaultContainer) getValue(t reflect.Type) reflect.Value {
	c.valueLock.RLock()
	defer c.valueLock.RUnlock()

	val := c.values[t]
	if val.IsValid() {
		return val
	}

	if t.Kind() == reflect.Interface {
		for k, v := range c.values {
			if k.Implements(t) {
				val = v
				break
			}
		}
	}
	return val
}

func (c *defaultContainer) Get(vals ...any) error {
	for _, val := range vals {
		v := reflect.ValueOf(val)
		if value := c.getValue(v.Type()); value.IsValid() {
			v.Set(value)
			continue
		}

		isSet := false
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
			if value := c.getValue(v.Type()); value.IsValid() {
				v.Set(value)
				isSet = true
				break
			}
		}

		if !isSet {
			return fmt.Errorf("Value not found for Type: %v", v.Type())
		}
	}
	return nil
}

func (c *defaultContainer) MustGet(vals ...any) {
	err := c.Get(vals...)
	if err != nil {
		panic(err)
	}
}
