package Errors

type MethodNotSentError struct {
	Method, Reason string
}

func (mnse *MethodNotSentError) Error() string {
	return "Unable to send " + mnse.Method + ". " + mnse.Reason
}
