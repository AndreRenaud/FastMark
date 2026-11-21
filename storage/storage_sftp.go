package storage

import (
	"bytes"
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
		return nil, err
	}

	return &SFTPStorage{
		client: client,
		cmd:    cmd,
		server: server,
		prefix: prefix,
	}, nil
}

func (s *SFTPStorage) fullPath(filename string) string {
	return filepath.Join(s.prefix, filename)
}

// Open fetches the whole file into memory before returning. A plain
// sftp.File costs one network round trip per Read call, which is very slow
// under the small buffered reads the image decoders make; WriteTo instead
// pipelines the read requests for the whole file concurrently.
func (s *SFTPStorage) Open(filename string) (io.ReadCloser, error) {
	f, err := s.client.Open(s.fullPath(filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		return nil, err
	}
	return io.NopCloser(&buf), nil
}

func (s *SFTPStorage) Describe() string {
	return fmt.Sprintf("sftp://%s/%s", s.server, s.prefix)
}

func (s *SFTPStorage) OpenWrite(filename string, append bool) (io.WriteCloser, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if append {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}
	fullname := s.fullPath(filename)

	// Ensure the directory exists
	s.client.MkdirAll(filepath.Dir(fullname))

	return s.client.OpenFile(fullname, flags)
}

func (s *SFTPStorage) Glob(directory string, pattern string) ([]string, error) {
	fullname := s.fullPath(directory)
	return s.client.Glob(filepath.Join(fullname, pattern))
}
