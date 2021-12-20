package cluster

import (
	"os"
	"path/filepath"
)

const (
	logDir = "log"
)

type logger struct {
	Path string
	f    *os.File
}

func newLogger(file string) (*logger, error) {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	path := filepath.Join(logDir, file)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &logger{
		Path: path,
		f:    f,
	}, nil
}

func (l *logger) Write(b []byte) (int, error) {
	return l.f.Write(b)
}
