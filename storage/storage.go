package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Storage interface {
	Open(filename string) (io.ReadCloser, error)
	OpenWrite(filename string, append bool) (io.WriteCloser, error)
	Glob(directory string, pattern string) ([]string, error)
	Describe() string
	Disconnect()
}

type LocalStorage struct {
	prefix string
}

type DummyStorage struct{}

var _ Storage = &LocalStorage{}
var _ Storage = &DummyStorage{}

func NewStorage(directory string) Storage {
	if strings.HasPrefix(directory, "sftp://") {
		parts, err := url.Parse(directory)
		if err != nil {
			return &DummyStorage{}
		}
		path := parts.Path
		// Haul the directory out separately
		parts.Path = ""
		//parts.Scheme = ""
		server := parts.String()
		server = strings.TrimPrefix(server, "sftp://")
		if s, err := NewSFTPStorage(server, path); err == nil {
			return s
		}
		return &DummyStorage{}
	}

	return &LocalStorage{prefix: directory}
}

func (s LocalStorage) fullPath(filename string) string {
	path := filepath.Clean(filename)
	full := filepath.Join(filepath.Clean(s.prefix), path)
	return full
}

func (s LocalStorage) Open(filename string) (io.ReadCloser, error) {
	return os.Open(s.fullPath(filename))
}

func (s LocalStorage) OpenWrite(filename string, append bool) (io.WriteCloser, error) {
	fullname := s.fullPath(filename)

	// Ensure the directory exists
	dirname := filepath.Dir(fullname)
	os.MkdirAll(dirname, 0755)

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if append {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	return os.OpenFile(fullname, flags, 0644)
}

func (s LocalStorage) Glob(directory string, pattern string) ([]string, error) {
	glob := filepath.Join(s.fullPath(directory), pattern)
	return filepath.Glob(glob)
}

func (s LocalStorage) Describe() string {
	return filepath.Clean(s.prefix)
}

func (s LocalStorage) Disconnect() {

}

func (d DummyStorage) Open(filename string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("dummy storage")
}
func (d DummyStorage) OpenWrite(filename string, append bool) (io.WriteCloser, error) {
	return nil, fmt.Errorf("dummy storage")
}
func (d DummyStorage) Glob(directory string, pattern string) ([]string, error) {
	return nil, fmt.Errorf("dummy storage")
}
func (d DummyStorage) Describe() string {
	return "dummy storage"
}

func (d DummyStorage) Disconnect() {}
