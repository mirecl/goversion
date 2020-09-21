package goversion

import (
	"fmt"
	"io"
)

// NewLumberjack logger file.
func NewLumberjack(path string, opts ...FileOption) (io.Writer, error) {
	fo := &fileOptions{
		path:    path,
		version: "null",
	}

	for _, opt := range opts {
		opt.apply(fo)
	}
	return fo, nil
}

func (fo *fileOptions) Write(p []byte) (n int, err error) {
	fmt.Println(fo.version)
	return 0, nil
}

type fileOptions struct {
	path    string
	version string
}

// FileOption configures file.
type FileOption interface {
	apply(*fileOptions)
}

type funcFileOption struct {
	f func(*fileOptions)
}

func (fdo *funcFileOption) apply(do *fileOptions) {
	fdo.f(do)
}

func newFuncFileOption(f func(*fileOptions)) *funcFileOption {
	return &funcFileOption{
		f: f,
	}
}

// WithVersion set version.
func WithVersion(version string) FileOption {
	return newFuncFileOption(func(o *fileOptions) {
		o.version = version
	})
}
