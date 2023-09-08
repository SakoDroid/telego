package parser

import (
	"sync"

	objs "github.com/SakoDroid/telego/objects"
)

type chatRequestHandler struct {
	requestId int
	function  *func(*objs.Update)
}

type chatRequestHandlerMap struct {
	sync.RWMutex
	internal map[int]*chatRequestHandler
}

func (c *chatRequestHandlerMap) Add(key int, value *chatRequestHandler) {
	c.Lock()
	c.internal[key] = value
	c.Unlock()
}

func (c *chatRequestHandlerMap) Load(key int) (value *chatRequestHandler, ok bool) {
	c.RLock()
	value, ok = c.internal[key]
	c.RUnlock()
	if ok {
		c.Delete(key)
	}
	return
}

func (c *chatRequestHandlerMap) Delete(key int) {
	c.Lock()
	delete(c.internal, key)
	c.Unlock()
}
