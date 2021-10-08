package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const BinaryContentType = "application/octet-stream"

type Client struct {
	addresses []string
	cl        http.Client
	offset    uint64
}

func NewClient(addresses []string) *Client {
	return &Client{
		addresses: addresses,
		cl:        http.Client{},
		offset:    0,
	}
}

func (c *Client) Receive(scratch []byte) ([]byte, error) {
	resp, err := c.cl.Get(fmt.Sprintf(c.addresses[0]+"/receive?offset=%d&maxSize=%d", c.offset, len(scratch)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)
		return nil, fmt.Errorf("http code %d, %s", resp.StatusCode, b.String())
	}

	off, err := strconv.Atoi(resp.Header.Get("offset"))
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(scratch[0:0])
	written, err := io.Copy(b, resp.Body)

	if err != nil {
		return nil, err
	}

	if written == 0 {
		return nil, io.EOF
	}

	c.offset += uint64(off)
	return b.Bytes(), nil
}

func (c *Client) Send(msgs []byte) error {
	resp, err := c.cl.Post(c.addresses[0]+"/send", BinaryContentType, bytes.NewReader(msgs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)
		return fmt.Errorf("http code %d, %s", resp.StatusCode, b.String())
	}

	io.Copy(ioutil.Discard, resp.Body)
	return nil
}
