package gorqlite

import (
	"fmt"
)

type gorqliteError struct {
	Inner   string
	Message string
}

func newError(messagef string, msgArgs ...interface{}) *gorqliteError {
	return &gorqliteError{
		Inner:   "",
		Message: fmt.Sprintf(messagef, msgArgs...),
	}
}

func wrapError(err error, messagef string) *gorqliteError {
	return &gorqliteError{
		Inner:   err.Error(),
		Message: messagef,
	}
}

func (err gorqliteError) Error() string {
	s := err.Message
	if err.Inner != "" {
		return fmt.Sprintf("%s: %s", s, err.Inner)
	}
	return s
}
