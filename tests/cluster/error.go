package cluster

import (
	"fmt"
)

type clusterError struct {
	Inner   string
	Message string
}

func newError(messagef string, msgArgs ...interface{}) *clusterError {
	return &clusterError{
		Inner:   "",
		Message: fmt.Sprintf(messagef, msgArgs...),
	}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) *clusterError {
	return &clusterError{
		Inner:   err.Error(),
		Message: fmt.Sprintf(messagef, msgArgs...),
	}
}

func (err clusterError) Error() string {
	s := err.Message
	if err.Inner != "" {
		return fmt.Sprintf("%s: %s", s, err.Inner)
	}
	return s
}

type errorWrapper struct {
	Message string
}

func newErrorWrapper(messagef string) *errorWrapper {
	return &errorWrapper{
		Message: messagef,
	}
}

func (w *errorWrapper) Error(err error) *clusterError {
	return wrapError(err, w.Message)
}
