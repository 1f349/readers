package readers

import (
	"bytes"
	"errors"
	"io"
)

var EOL error = &eolError{}

type eolError struct{}

func (e *eolError) Error() string { return "EOL" }
func (e *eolError) Unwrap() error { return io.EOF }

type LineReader struct {
	err  error
	r    io.Reader
	over []byte
}

var _ io.Reader = &LineReader{}

func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{nil, r, nil}
}

func (l *LineReader) Err() error {
	if errors.Is(l.err, io.EOF) {
		return nil
	}
	return l.err
}

func (l *LineReader) Next() bool {
	if l.err == nil || errors.Is(l.err, EOL) {
		l.err = nil
		return true
	}
	return false
}

func (l *LineReader) Read(p []byte) (n int, err error) {
	if errors.Is(l.err, EOL) {
		return 0, io.EOF
	}

	rOver := bytes.NewBuffer(l.over)
	r := io.MultiReader(rOver, l.r)
	n, l.err = r.Read(p)
	if l.err != nil {
		err = eolToEof(l.err)
		return
	}
	nFull := n

	nStart := n
	n2 := bytes.IndexByte(p, '\n')
	if n2 != -1 {
		l.err = EOL
		n = n2
		nStart = n2 + 1
	}

	if n > 0 && p[n-1] == '\r' {
		n--
	}

	l.over = p[nStart:nFull]
	if rOver.Len() > 0 {
		l.over = append(p[nStart:nFull], rOver.Bytes()...)
	}
	return
}

func eolToEof(err error) error {
	if errors.Is(err, EOL) {
		return io.EOF
	}
	return err
}
