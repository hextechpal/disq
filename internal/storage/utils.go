package storage

import (
	"bytes"
	"errors"
)

func cutAtNewLine(in []byte) ([]byte, []byte, error) {
	size := len(in)
	bIdx := bytes.LastIndexByte(in, '\n')

	if bIdx < 0 {
		return nil, nil, errors.New("buffer too small")
	}

	if bIdx == size-1 {
		return in, nil, nil
	}

	return in[:bIdx+1], in[bIdx+1:], nil
}
