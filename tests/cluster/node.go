package cluster

import (
	"io"
)

type Node interface {
	Reboot(duration int64, timeout bool) error
	io.Closer
}
