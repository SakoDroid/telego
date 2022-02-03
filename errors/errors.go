package errors

import (
	"strconv"

	objs "github.com/SakoDroid/telego/objects"
)

//MethodNotSentError is returned when API server responds with any code other than 200
type MethodNotSentError struct {
	Method, Reason string
	FailureResult  *objs.FailureResult
}

func (mnse *MethodNotSentError) Error() string {
	out := "Unable to send " + mnse.Method + ". " + mnse.Reason
	if mnse.FailureResult != nil {
		out += "\nError code : " + strconv.Itoa(mnse.FailureResult.ErrorCode) + ", Description : " + mnse.FailureResult.Description
	}
	return out
}

//BotInterfaceAlreadyCreated indicates that the bai is already created.
type BotInterfaceAlreadyCreated struct {
}

func (biac *BotInterfaceAlreadyCreated) Error() string {
	return "only one bot interface can be created"
}

//UpdateRoutineAlreadyStarted indicates that the bot hase been started.
type UpdateRoutineAlreadyStarted struct {
}

func (uras *UpdateRoutineAlreadyStarted) Error() string {
	return "StartUpdateRoutine has already been called."
}

//UpdateNotOk indicates that server returned "ok : false" in response.
type UpdateNotOk struct {
	Offset int
}

func (uno *UpdateNotOk) Error() string {
	return "getUpdates returned \"ok\" : false. Offset used in request : " + strconv.Itoa(uno.Offset)
}

//RequiredArgumentError indicates that a required argument is missing or has a problem.
type RequiredArgumentError struct {
	ArgName, MethodName string
}

func (ram *RequiredArgumentError) Error() string {
	return "Required argument \"" + ram.ArgName + " missing or has wrong value in " + ram.MethodName + " method."
}

//ChatIdProblem indicates a problem in the chat id.
type ChatIdProblem struct {
}

func (cip *ChatIdProblem) Error() string {
	return "Cannot have both chatIdInt and chatIdString at the same time. Only one of them is allowed."
}

//MediaGroupFullError indicates that the media group is full.
type MediaGroupFullError struct {
}

func (mgfe *MediaGroupFullError) Error() string {
	return "the media group is full. (10 files)"
}

//LiveLocationNotStarted indicates that the live location has not been started yet.
type LiveLocationNotStarted struct {
}

func (llns *LiveLocationNotStarted) Error() string {
	return "live location has not been started (sent)."
}
