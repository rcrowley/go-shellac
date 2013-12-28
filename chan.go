package shellac

import (
	"bytes"
	"io"
)

// ChanReader is a bridge between the receive side of a channel and io.Reader.
type ChanReader struct {
	buffer bytes.Buffer
	ch     <-chan string
}

// NewChanReader constructs a ChanReader that reads from ch.
func NewChanReader(ch <-chan string) *ChanReader {
	return &ChanReader{
		buffer: bytes.Buffer{},
		ch:     ch,
	}
}

// Read reads lines from a channel and returns them to callers.  It eventually
// returns io.EOF after the channel is closed and the buffer is emptied.
func (r *ChanReader) Read(p []byte) (int, error) {
	n, err := r.buffer.Read(p)
	if io.EOF == err {
		s, ok := <-r.ch
		if !ok {
			return 0, io.EOF
		}
		if n, err := r.buffer.Write([]byte(s)); nil != err || len(s) != n {
			return 0, err
		}
		if 0 == len(s) || '\n' != s[len(s)-1] {
			if _, err := r.buffer.Write([]byte{'\n'}); nil != err {
				return 0, err
			}
		}
		return r.Read(p)
	}
	return n, err
}

// ChanWriter is a bridge between the sending side of a channel and io.Writer.
type ChanWriter struct {
	buffer bytes.Buffer
	ch     chan<- string
}

// NewChanWriter constructs a ChanWriter that sends to ch.
func NewChanWriter(ch chan<- string) *ChanWriter {
	return &ChanWriter{
		buffer: bytes.Buffer{},
		ch:     ch,
	}
}

// Close drains the buffer and closes the channel.
func (w *ChanWriter) Close() error {
	for {
		s, err := w.buffer.ReadString('\n')
		if nil != err {
			if "" != s {
				w.ch <- s
			}
			break
		}
		w.ch <- s[:len(s)-1]
	}
	close(w.ch)
	return nil
}

// Write sends lines to a channel as they are passed by callers.
func (w *ChanWriter) Write(p []byte) (int, error) {
	n, err := w.buffer.Write(p)
	for i := bytes.Count(p, []byte{'\n'}); i > 0; i-- {
		s, err := w.buffer.ReadString('\n')
		if nil != err {
			break
		}
		w.ch <- s[:len(s)-1]
	}
	return n, err
}
