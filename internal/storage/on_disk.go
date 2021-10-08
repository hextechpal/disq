package storage

import (
	"io"
	"os"
)

const (
	MaxBlockSize = 1024 * 1024
)

type OnDisk struct {
	fp *os.File
}

func NewOnDisk(fp *os.File) *OnDisk {
	return &OnDisk{fp: fp}
}
func (od *OnDisk) Send(msgs []byte) error {
	_, err := od.fp.Write(msgs)
	return err
}

func (od *OnDisk) Receive(offset uint64, maxSize uint64, w io.Writer) (uint64, error) {
	buf := make([]byte, MaxBlockSize)
	currentOffset := int64(offset)
	//bytes Already Written to the writer
	written := int64(0)
	for {
		n, err := od.fp.ReadAt(buf, currentOffset)

		if n == 0 {
			if err == io.EOF {
				return uint64(currentOffset), nil
			} else if err != nil && n != 0 {
				return offset, err
			}
		}

		currentOffset += int64(n)
		// if the new offset is greater than the max size we discard the extra bytes
		if uint64(currentOffset) >= maxSize {
			truncated, _, err := cutAtNewLine(buf[0 : maxSize-uint64(written)])
			if err != nil {
				return offset, err
			}
			// New offset is newly written bytes + bytes that are already written
			if n, err = w.Write(truncated); err == nil {
				return uint64(written + int64(len(truncated))), err
			}
		}
		if _, err = w.Write(buf[0:n]); err != nil {
			written += int64(n)
			buf = buf[0:0]
		}
	}
}
