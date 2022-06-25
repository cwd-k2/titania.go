package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// compare two files' contents
func Equal(a, b string) (bool, error) {
	fpA, err := os.Open(a)
	if err != nil {
		return false, err
	}
	defer fpA.Close()

	fpB, err := os.Open(b)
	if err != nil {
		return false, err
	}
	defer fpB.Close()

	bytesA, err := io.ReadAll(bufio.NewReader(fpA))
	if err != nil {
		return false, nil
	}

	bytesB, err := io.ReadAll(bufio.NewReader(fpB))
	if err != nil {
		return false, nil
	}

	return bytes.Equal(bytesA, bytesB), nil
}
