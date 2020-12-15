package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aiocloud/shadowsocks/socks"
)

func tcpListen() {
	for {
		client, err := tcpListenter.Accept()
		if err != nil {
			fmt.Printf("[shadowsocks][tcpListenter.Accept] %v", err)
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
		fmt.Printf("[shadowsocks][net.Dial] %v", err)
		return
	}
	remote = cipher.StreamConn(remote)
	defer remote.Close()

	if _, err = remote.Write(addr); err != nil {
		fmt.Printf("[shadowsocks][remote.Write] %v", err)
		return
	}
	addr = nil

	fmt.Printf("[shadowsocks] New TCP connection from %s to %s", client.RemoteAddr(), tcpRemoteAddr)

	go func() {
		_, _ = io.CopyBuffer(remote, client, make([]byte, 1400))
		_ = client.SetDeadline(time.Now())
		_ = remote.SetDeadline(time.Now())
	}()

	_, _ = io.CopyBuffer(client, remote, make([]byte, 1400))
	_ = client.SetDeadline(time.Now())
	_ = remote.SetDeadline(time.Now())
}
