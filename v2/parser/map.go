package parser

import "sync"

type threadSafeMap[k comparable, v any] struct {
	sync.RWMutex
	internal map[k]v
}

func (c *threadSafeMap[k, v]) Add(key k, value v) {
	c.Lock()
	c.internal[key] = value
	c.Unlock()
}

func (c *threadSafeMap[k, v]) LoadAndDelete(key k) (value v, ok bool) {
	c.RLock()
	value, ok = c.internal[key]
	c.RUnlock()
	if ok {
		c.Delete(key)
	}
	return
}

func (c *threadSafeMap[k, v]) Load(key k) (value v, ok bool) {
	c.RLock()
	value, ok = c.internal[key]
	c.RUnlock()
	return
}

func (c *threadSafeMap[k, v]) Delete(key k) {
	c.Lock()
	delete(c.internal, key)
	c.Unlock()
}
