package readers

import (
	"bytes"
	"io"
)

type CommentReader struct {
	r    io.Reader
	over []byte
	mark []string
}

var _ io.Reader = &CommentReader{}

func NewCommentReader(r io.Reader, mark []string) *CommentReader {
	return &CommentReader{r, nil, mark}
}

func (s *CommentReader) Read(p []byte) (n int, err error) {
	if s.over != nil {
		return 0, io.EOF
	}
	n, err = s.r.Read(p)
	if err != nil {
		return
	}

	n2 := s.matchesMark(p[:n])
	if n2 != -1 {
		s.over = p[n2:n]
		n = n2
		err = io.EOF
	}
	return
}

func (s *CommentReader) matchesMark(p []byte) int {
	for _, i := range s.mark {
		n := bytes.Index(p, []byte(i))
		if n != -1 {
			return n
		}
	}
	return -1
}

func (s *CommentReader) Comment() io.Reader {
	return io.MultiReader(bytes.NewReader(s.over), s.r)
}
