package transport

import (
	"errors"
	"io"
	"sync"
)

// Stream joins owned read/write handles, for example child stdin/stdout pipes.
// Close must unblock Read and Write; do not pass a shared application stream.
type Stream struct {
	Reader io.ReadCloser
	Writer io.WriteCloser
	once   sync.Once
	err    error
}

func (s *Stream) Read(b []byte) (int, error)  { return s.Reader.Read(b) }
func (s *Stream) Write(b []byte) (int, error) { return s.Writer.Write(b) }
func (s *Stream) Close() error {
	s.once.Do(func() { s.err = errors.Join(s.Reader.Close(), s.Writer.Close()) })
	return s.err
}
