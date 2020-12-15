package main

import (
	"io"
	"net"
	"time"

	"github.com/aiocloud/shadowsocks/socks"
)

func tcpListen() {
	for {
		client, err := tcpListenter.Accept()
		if err != nil {
			return
		}

		go tcpHandle(client)
	}
}

func tcpHandle(client net.Conn) {
	defer client.Close()

	addr, err := socks.Handshake(client)
	if err != nil {
		if err == socks.InfoUDPAssociate {
			buffer := make([]byte, 1)

			for {
				if _, err = client.Read(buffer); err != nil {
					return
				}
			}
		}

		return
	}

	remote, err := net.Dial("tcp", tcpRemoteAddr.String())
	if err != nil {
		return
	}
	remote = cipher.StreamConn(remote)
	defer remote.Close()

	if _, err = remote.Write(addr); err != nil {
		return
	}
	addr = nil

	go func() {
		io.CopyBuffer(remote, client, make([]byte, 1400))
		client.SetDeadline(time.Now())
		remote.SetDeadline(time.Now())
	}()

	io.CopyBuffer(client, remote, make([]byte, 1400))
	client.SetDeadline(time.Now())
	remote.SetDeadline(time.Now())
}
