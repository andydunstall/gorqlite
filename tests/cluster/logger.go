package cluster

import (
	"os"
	"path/filepath"

	"github.com/dunstall/gorqlite"
)

const (
	logDir = "log"
)

type logger struct {
	Path string
	f    *os.File
}

func newLogger(file string) (*logger, error) {
	path := filepath.Join(logDir, file)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, gorqlite.WrapError(err, "failed to open logger: %s", file)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, gorqlite.WrapError(err, "failed to open logger: %s", path)
	}
	return &logger{
		Path: path,
		f:    f,
	}, nil
}

func (l *logger) Write(b []byte) (int, error) {
	n, err := l.f.Write(b)
	if err != nil {
		return n, gorqlite.WrapError(err, "failed to write to logger: %s", l.Path)
	}
	return n, nil
}
