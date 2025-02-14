package main

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	prefix string
}

func (s Storage) fullPath(filename string) string {
	path := filepath.Clean(filename)
	full := filepath.Join(filepath.Clean(s.prefix), path)
	return full
}

func (s Storage) Open(filename string) (io.ReadCloser, error) {
	return os.Open(s.fullPath(filename))
}

func (s Storage) OpenWrite(filename string, append bool) (io.WriteCloser, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if append {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	return os.OpenFile(s.fullPath(filename), flags, 0644)
}

func (s Storage) Glob(directory string, pattern string) ([]string, error) {
	glob := filepath.Join(s.fullPath(directory), pattern)
	return filepath.Glob(glob)
}

func (s Storage) Describe() string {
	return filepath.Clean(s.prefix)
}
