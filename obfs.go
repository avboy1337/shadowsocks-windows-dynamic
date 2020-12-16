package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
)

type HTTP struct {
	net.Conn

	hasSentHeader bool
	hasRecvHeader bool

	buffer []byte
	offset int
}

func (o *HTTP) Read(data []byte) (int, error) {
	if o.buffer != nil {
		size := copy(data, o.buffer[o.offset:])
		o.offset += size

		if len(o.buffer) == o.offset {
			o.buffer = nil
		}

		return size, nil
	}

	if !o.hasRecvHeader {
		buffer := make([]byte, 1400)
		size, err := o.Conn.Read(buffer)
		if err != nil {
			return 0, err
		}
		buffer = buffer[:size]

		idx := bytes.Index(buffer, []byte{0x0d, 0x0a, 0x0d, 0x0a})
		if idx == -1 {
			return 0, nil
		}
		buffer = buffer[idx+4:]

		o.hasRecvHeader = true

		length := copy(data, buffer)
		if len(buffer) > length {
			o.buffer = buffer
			o.offset = length
		}

		return length, nil
	}

	return o.Conn.Read(data)
}

func (o *HTTP) Write(data []byte) (int, error) {
	if !o.hasSentHeader {
		random := make([]byte, 16)
		_, _ = rand.Read(random)

		if _, err := o.Conn.Write(append([]byte(
			fmt.Sprintf(""+
				"GET / HTTP/1.1\r\n"+
				"Connection: Upgrade\r\n"+
				"Upgrade: websocket\r\n"+
				"Host: %s\r\n"+
				"User-Agent: curl/7.64.0\r\n"+
				"Sec-WebSocket-Key: %s\r\n"+
				"Content-Length: %d\r\n"+
				"\r\n"+
				"",
				OBFSParam,
				base64.URLEncoding.EncodeToString(random),
				len(data),
			),
		), data...)); err != nil {
			return 0, err
		}

		return len(data), nil
	}

	return o.Conn.Write(data)
}

func newHTTP(client net.Conn) net.Conn {
	return &HTTP{Conn: client}
}
