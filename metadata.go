// 代码来源于gin.Context
package container

import (
	"sync"
	"time"
)

type Metadata struct {
	mu   sync.RWMutex
	Keys map[string]interface{}
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  m.Keys if it was not used previously.
func (m *Metadata) Set(key string, value interface{}) {
	m.mu.Lock()
	if m.Keys == nil {
		m.Keys = make(map[string]interface{})
	}

	m.Keys[key] = value
	m.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (m *Metadata) Get(key string) (value interface{}, exists bool) {
	m.mu.RLock()
	value, exists = m.Keys[key]
	m.mu.RUnlock()
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (m *Metadata) MustGet(key string) interface{} {
	if value, exists := m.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (m *Metadata) GetString(key string) (s string) {
	if val, ok := m.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (m *Metadata) GetBool(key string) (b bool) {
	if val, ok := m.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (m *Metadata) GetInt(key string) (i int) {
	if val, ok := m.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (m *Metadata) GetInt64(key string) (i64 int64) {
	if val, ok := m.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (m *Metadata) GetFloat64(key string) (f64 float64) {
	if val, ok := m.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (m *Metadata) GetTime(key string) (t time.Time) {
	if val, ok := m.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (m *Metadata) GetDuration(key string) (d time.Duration) {
	if val, ok := m.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (m *Metadata) GetStringSlice(key string) (ss []string) {
	if val, ok := m.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (m *Metadata) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := m.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (m *Metadata) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := m.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (m *Metadata) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := m.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
