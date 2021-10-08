package storage

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
)

const (
	MaxNumbers = 1000000
	ChunkSize  = 1024 * 1024
)

func TestInMemory_SendReceive(t *testing.T) {
	im := NewInMemory()
	want, _ := send(im)
	got, _ := receive(im)

	if want != got {
		t.Errorf("Different sums calculated want %d, got %d", want, got)
	}
	t.Logf("Test Passes want %d, got %d", want, got)
}

func send(im *InMemory) (sum int64, err error) {
	var b bytes.Buffer
	for i := 0; i < MaxNumbers; i++ {
		sum += int64(i)
		_, _ = fmt.Fprintf(&b, "%d\n", i)
		if b.Len() >= ChunkSize {
			err := im.Send(b.Bytes())
			if err != nil {
				return 0, err
			}
			b.Reset()
		}
	}

	if b.Len() != 0 {
		err := im.Send(b.Bytes())
		if err != nil {
			return 0, err
		}
		b.Reset()
	}
	return sum, nil
}

func TestInMemory_CutToNewLine(t *testing.T) {
	res := []byte("100\n101\n10")
	wantOut, wantRes := []byte("100\n101\n"), []byte("10")
	gotOut, gotRes, err := cutAtNewLine(res)

	if !bytes.Equal(wantOut, gotOut) || !bytes.Equal(wantRes, gotRes) || err != nil {
		t.Errorf("bute arrays not equest gotOut %v, gotRes %v", gotOut, gotRes)
	}

}

func receive(im *InMemory) (sum int64, err error) {
	off := uint64(0)
	buf := bytes.NewBuffer(nil)
	for {
		off, err = im.Receive(off, ChunkSize, buf)
		if err == io.EOF {
			return sum, nil
		} else if err != nil {
			return 0, err
		} else {
			ints := strings.Split(string(buf.Bytes()), "\n")
			for _, val := range ints {
				num, _ := strconv.Atoi(val)
				sum += int64(num)
			}
			buf.Reset()
		}
	}
}
