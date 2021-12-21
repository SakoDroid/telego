package parser

import (
	"encoding/json"

	errs "github.com/SakoDroid/telebot/Errors"
	objs "github.com/SakoDroid/telebot/objects"
)

func ParseUpdate(body []byte, uc *chan *objs.Update, pu *chan *objs.Update) (int, error) {
	def := &objs.DefaultResult{}
	err2 := json.Unmarshal(body, def)
	if err2 != nil {
		return 0, err2
	}
	if !def.Ok {
		return 0, &errs.MethodNotSentError{Method: "getUpdates", Reason: "server returned false for \"ok\" field."}
	}
	ur := &objs.UpdateResult{}
	err := json.Unmarshal(body, ur)
	if err != nil {
		return 0, err
	}
	return parse(ur, uc, pu)
}

func parse(ur *objs.UpdateResult, uc *chan *objs.Update, pu *chan *objs.Update) (int, error) {
	lastOffset := 0
	for _, val := range ur.Result {
		if val.Update_id > lastOffset {
			lastOffset = val.Update_id
		}
		if val.Poll.Id != "" {
			*pu <- val
		} else {
			*uc <- val
		}
	}
	return lastOffset, nil
}
