package storage

import (
	"io"
)

type InMemory struct {
	b []byte
}

func NewInMemory() *InMemory {
	return &InMemory{}
}
func (im *InMemory) Send(msgs []byte) error {
	im.b = append(im.b, msgs...)
	return nil
}

func (im *InMemory) Receive(offset uint64, maxSize uint64, w io.Writer) (uint64, error) {

	if offset >= uint64(len(im.b)) {
		return offset, nil
	}

	if offset+maxSize >= uint64(len(im.b)) {
		write, err := w.Write(im.b[offset:])
		if err != nil {
			return offset, err
		}
		return offset + uint64(write), nil
	}

	ret, _, err := cutAtNewLine(im.b[offset : offset+maxSize])
	if err != nil {
		return 0, err
	}
	write, err := w.Write(ret)
	if err != nil {
		return 0, err
	}
	return offset + uint64(write), nil

}
