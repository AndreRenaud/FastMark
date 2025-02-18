package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/sftp"
)

type SFTPStorage struct {
	client *sftp.Client
	cmd    *exec.Cmd
	server string
	prefix string
}

var _ Storage = &SFTPStorage{}

func NewSFTPStorage(server string, prefix string) (*SFTPStorage, error) {
	// Connect to a remote host and request the sftp subsystem via the 'ssh'
	// command.  This assumes that passwordless login is correctly configured.
	log.Printf("SSH connecting to %s directory %s", server, prefix)
	cmd := exec.Command("ssh", server, "-s", "sftp")

	// send errors from ssh to stderr
	// TODO: Expose these somehow?
	cmd.Stderr = os.Stderr

	// get stdin and stdout
	wr, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	rd, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// start the process
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// open the SFTP session
	client, err := sftp.NewClientPipe(rd, wr)
	if err != nil {
		cmd.Process.Kill()
		return nil, nil
	}

	return &SFTPStorage{
		client: client,
		cmd:    cmd,
		server: server,
		prefix: prefix,
	}, nil
}

func (s *SFTPStorage) fullPath(filename string) string {
	path := filepath.Clean(filename)
	full := filepath.Join(filepath.Clean(s.prefix), path)
	return full
}

func (s *SFTPStorage) Open(filename string) (io.ReadCloser, error) {
	fullname := s.fullPath(filename)
	return s.client.Open(fullname)
}

func (s *SFTPStorage) Describe() string {
	return fmt.Sprintf("sftp://%s/%s", s.server, s.prefix)
}

func (s *SFTPStorage) Disconnect() {
	s.client.Close()
	s.cmd.Process.Kill()
}

func (s *SFTPStorage) OpenWrite(filename string, append bool) (io.WriteCloser, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if append {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}
	fullname := s.fullPath(filename)
	return s.client.OpenFile(fullname, flags)
}

func (s *SFTPStorage) Glob(directory string, pattern string) ([]string, error) {
	fullname := s.fullPath(directory)
	return s.client.Glob(filepath.Join(fullname, pattern))
}
