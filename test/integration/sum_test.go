package integration

import (
	"bytes"
	"fmt"
	"github.com/ppal31/disq/cli/client"
	"io"
	"log"
	"strconv"
	"strings"
	"testing"
)

const (
	MaxNumbers = 1000000
	ChunkSize  = 1024 * 1024
)

func TestDisq_Sum(t *testing.T) {
	c := client.NewClient([]string{"http://localhost:8081"})
	want, _ := send(c)
	got, _ := receive(c)

	if want != got {
		t.Errorf("Different sums calculated want %d, got %d", want, got)
		t.Fail()
	} else {
		t.Logf("Test Passes want %d, got %d", want, got)
	}
}

func send(c *client.Client) (sum int64, err error) {
	var b bytes.Buffer
	for i := 1; i <= MaxNumbers; i++ {
		sum += int64(i)
		_, _ = fmt.Fprintf(&b, "%d\n", i)
		if b.Len() >= ChunkSize {
			err := c.Send(b.Bytes())
			if err != nil {
				return 0, err
			}
			b.Reset()
		}
	}

	if b.Len() != 0 {
		err := c.Send(b.Bytes())
		if err != nil {
			return 0, err
		}
		b.Reset()
	}
	return sum, nil
}

func receive(c *client.Client) (sum int64, err error) {
	buf := make([]byte, ChunkSize)
	for {
		b, err := c.Receive(buf)
		if err == io.EOF {
			return sum, nil
		} else if err != nil {
			log.Printf(err.Error())
			return 0, err
		} else {
			ints := strings.Split(string(b), "\n")
			for _, val := range ints {
				num, _ := strconv.Atoi(val)
				//log.Printf("Number Received: %d\n", num)
				sum += int64(num)
			}
			buf = make([]byte, ChunkSize)
		}
	}
}
