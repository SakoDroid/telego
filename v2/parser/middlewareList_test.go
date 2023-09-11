package parser

import (
	"testing"

	objs "github.com/SakoDroid/telego/v2/objects"
)

func TestMiddleList(t *testing.T) {
	list := &middlewareLinkedList{}

	fakeUpdate := &objs.Update{
		Update_id: 1234,
		Message: &objs.Message{
			MessageId: 1234,
		},
		Poll: &objs.Poll{
			Id: "1234",
		},
	}

	list.addToBegin(func(update *objs.Update, next func()) {
		update.Update_id = 4567
		next()
	})
	list.addToBegin(func(update *objs.Update, next func()) {
		update.Message.MessageId = 4567
		next()
	})
	list.addToBegin(func(update *objs.Update, next func()) {
		update.Poll.Id = "4567"
		next()
	})

	list.executeChain(fakeUpdate)

	if fakeUpdate.Update_id != 4567 || fakeUpdate.Message.MessageId != 4567 || fakeUpdate.Poll.Id != "4567" {
		t.FailNow()
	}
}
