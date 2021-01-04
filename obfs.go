package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
)

type HTTP struct {
	net.Conn

	hasRecvHeader bool
	hasSentHeader bool

	buffer []byte
	offset int
}

func (c *HTTP) Read(data []byte) (int, error) {
	if c.buffer != nil {
		size := copy(data, c.buffer[c.offset:])
		c.offset += size

		if len(c.buffer) == c.offset {
			c.buffer = nil
		}

		return size, nil
	}

	if !c.hasRecvHeader {
		buffer := make([]byte, 1400)
		length, err := c.Conn.Read(buffer)
		if err != nil {
			return 0, err
		}

		idx := bytes.Index(buffer[:length], []byte{0x0d, 0x0a, 0x0d, 0x0a})
		if idx == -1 {
			return 0, io.EOF
		}
		buffer = buffer[idx+4 : length]

		c.hasRecvHeader = true

		size := copy(data, buffer)
		if len(buffer) != size {
			c.buffer = buffer
			c.offset = size
		}

		return size, nil
	}

	return c.Conn.Read(data)
}

func (c *HTTP) Write(data []byte) (int, error) {
	if !c.hasSentHeader {
		random := make([]byte, 16)
		_, _ = rand.Read(random)

		header := []byte(fmt.Sprintf(""+
			"GET / HTTP/1.1\r\n"+
			"Content-Length: %d\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Host: %s\r\n"+
			"User-Agent: curl/\r\n"+
			"Sec-WebSocket-Key: %s\r\n"+
			"\r\n",
			len(data),
			OBFSParam,
			base64.URLEncoding.EncodeToString(random),
		))

		if _, err := c.Conn.Write(append(header, data...)); err != nil {
			return 0, err
		}

		return len(data), nil
	}

	return c.Conn.Write(data)
}
