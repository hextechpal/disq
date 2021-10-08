package storage

import "io"

type Storage interface {
	Send(msgs []byte) error
	Receive(offset uint64, maxSize uint64, w io.Writer) (newOffset uint64, err error)
}
