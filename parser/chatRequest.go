package parser

import (
	objs "github.com/SakoDroid/telego/objects"
)

type chatRequestHandler struct {
	requestId int
	function  *func(*objs.Update)
}
