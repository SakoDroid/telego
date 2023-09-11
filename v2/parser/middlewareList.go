package parser

import (
	"sync"

	objs "github.com/SakoDroid/telego/v2/objects"
)

var middlewares = &middlewareLinkedList{}

type middlewareListMember struct {
	prev, next *middlewareListMember
	internal   func(*objs.Update, func())
}

func (l *middlewareListMember) execute(up *objs.Update) {
	l.internal(up, func() {
		if l.next != nil {
			l.next.execute(up)
		}
	})
}

type middlewareLinkedList struct {
	sync.RWMutex
	first *middlewareListMember
}

func (l *middlewareLinkedList) addToBegin(elm func(update *objs.Update, next func())) {
	l.Lock()
	mem := &middlewareListMember{
		prev:     nil,
		next:     l.first,
		internal: elm,
	}
	if l.first != nil {
		l.first.prev = mem
	}
	l.first = mem
	l.Unlock()
}

func (l *middlewareLinkedList) executeChain(up *objs.Update) {
	l.first.execute(up)
}

func AddMiddleWare(middleware func(update *objs.Update, next func())) {
	middlewares.addToBegin(middleware)
}
