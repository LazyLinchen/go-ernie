package go_ernie

import (
	"errors"
	"io"
)

type peekingReader struct {
	r       io.Reader
	buf     []byte
	readIdx int
}

func newPeekingReader(r io.Reader, peekSize int) (*peekingReader, error) {
	buf := make([]byte, peekSize)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}
	return &peekingReader{
		r:       r,
		buf:     buf[:n],
		readIdx: 0,
		
	}, nil
}

func (pr *peekingReader) Read(p []byte) (n int, err error) {
	if len(pr.buf[pr.readIdx:]) > 0 {
		n = copy(p, pr.buf[pr.readIdx:])
		pr.readIdx += n
	}
	if n == len(p) || pr.readIdx == len(pr.buf) {
		return n, nil
	}
	moreN, moreErr := pr.r.Read(p[n:])
	return n + moreN, moreErr
}
