package goversion

import (
	"fmt"
	"io"
	"os"
)

// NewLumberjack logger file.
func NewLumberjack(filename string, opts ...FileOption) (io.Writer, error) {
	fo := &fileOptions{
		filename: filename,
		version:  "null",
	}

	for _, opt := range opts {
		opt.apply(fo)
	}
	return fo, nil
}

func (fo *fileOptions) Write(p []byte) (n int, err error) {
	fmt.Println(fo)
	if fo.file == nil {
		fo.openExistingOrNew()
	}
	if fo.f != nil {
		fo.f()
	}
	return fo.file.Write(p)
}

type fileOptions struct {
	filename string
	version  string
	size     int64
	backup   bool
	file     *os.File
	f        func()
}

func (fo *fileOptions) openExistingOrNew() error {
	info, err := os.Stat(fo.filename)
	if os.IsNotExist(err) {
		return fo.openNew()
	}

	file, err := os.OpenFile(fo.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fo.openNew()
	}

	fo.file = file
	fo.size = info.Size()
	return nil
}

func (fo *fileOptions) openNew() error {
	name := fo.filename
	mode := os.FileMode(0600)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	fo.file = f
	fo.size = 0
	return nil
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

// WithBufferSize set buffer size.
func WithBufferSize(size int64) FileOption {
	return newFuncFileOption(func(o *fileOptions) {
		o.size = size
	})
}

// WithBackup set backup.
func WithBackup() FileOption {
	return newFuncFileOption(func(o *fileOptions) {
		o.backup = true
	})
}

// WithCallBack ...
func WithCallBack(f func()) FileOption {
	return newFuncFileOption(func(o *fileOptions) {
		o.f = f
	})
}
