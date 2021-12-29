package gorqlite

import (
	"fmt"
	"runtime/debug"
)

type Error struct {
	Inner      string
	StackTrace string
	Message    string
}

func NewError(messagef string, msgArgs ...interface{}) *Error {
	return &Error{
		Inner:      "",
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
	}
}

func WrapError(err error, messagef string, msgArgs ...interface{}) *Error {
	return &Error{
		Inner:      err.Error(),
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
	}
}

func (err Error) Error() string {
	s := err.Message
	if err.Inner != "" {
		return fmt.Sprintf("%s: %s", s, err.Inner)
	}
	return s
}

type ErrorWrapper struct {
	Message string
}

func NewErrorWrapper(messagef string, msgArgs ...interface{}) *ErrorWrapper {
	return &ErrorWrapper{
		Message: fmt.Sprintf(messagef, msgArgs...),
	}
}

func (w *ErrorWrapper) Error(err error) *Error {
	return WrapError(err, w.Message)
}
