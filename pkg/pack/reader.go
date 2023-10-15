package pack

import (
	"bufio"
	"io"
	"strconv"
)

type PackReader struct {
	r bufio.Reader
}

func (pr *PackReader) Read(b []byte) (n int, err error) {
	chunklen, err := pr.r.Peek(4)
	if err != nil {
		return
	}

	shouldRead, err := strconv.ParseInt(string(chunklen), 16, 16)
	if err != nil {
		return
	}

	if shouldRead == 0 {
		// "0000"
		shouldRead = 4
	}

	output, err := pr.r.Peek(int(shouldRead))
	if err != nil {
		return
	}

	copy(b, output)
	pr.r.Discard(len(output))
	n = len(output)
	return
}

func NewReader(r io.Reader) *PackReader {
	bufferedReader := bufio.NewReader(r)
	return &PackReader{*bufferedReader}
}
