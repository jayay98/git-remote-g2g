package pack

import (
	"bufio"
	"io"
	"strconv"
)

func packSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	chunklen, err := strconv.ParseInt(string(data[:4]), 16, 16)
	if err != nil {
		return 0, nil, err
	}

	if chunklen == 0 {
		// "0000"
		chunklen = 4
	}
	return int(chunklen), data[:chunklen], nil
}

func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(NewReader(r))
	scanner.Split(packSplit)
	return scanner
}
