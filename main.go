package main

import "C"
import (
	"fmt"
	"net"

	"github.com/aiocloud/shadowsocks/core"
	"github.com/aiocloud/shadowsocks/socks"
)

var (
	ListenAddr string
	RemoteAddr string

	cipher        core.Cipher
	tcpListenter  net.Listener
	udpListenter  net.PacketConn
	tcpRemoteAddr net.Addr
	udpRemoteAddr net.Addr
)

//export ServerInfo
func ServerInfo(client, remote, passwd, method *C.char) bool {
	socks.UDPEnabled = true

	ListenAddr = C.GoString(client)
	RemoteAddr = C.GoString(remote)

	var err error
	if cipher, err = core.PickCipher(C.GoString(method), nil, C.GoString(passwd)); err != nil {
		fmt.Printf("[shadowsocks][core.PickCipher] %v", err)
		return false
	}

	if tcpRemoteAddr, err = net.ResolveTCPAddr("tcp", RemoteAddr); err != nil {
		fmt.Printf("[shadowsocks][net.ResolveTCPAddr] %v", err)
		return false
	}

	if udpRemoteAddr, err = net.ResolveUDPAddr("udp", RemoteAddr); err != nil {
		fmt.Printf("[shadowsocks][net.ResolveUDPAddr] %v", err)
		return false
	}

	return true
}

//export Create
func Create() bool {
	var err error

	if tcpListenter, err = net.Listen("tcp", ListenAddr); err != nil {
		fmt.Printf("[shadowsocks][net.Listen] %v", err)

		Delete()
		return false
	}

	if udpListenter, err = net.ListenPacket("udp", ListenAddr); err != nil {
		fmt.Printf("[shadowsocks][net.ListenPacket] %v", err)
		Delete()

		return false
	}

	return true
}

//export Delete
func Delete() {
	if tcpListenter != nil {
		tcpListenter.Close()
	}
	tcpListenter = nil

	if udpListenter != nil {
		udpListenter.Close()
	}
	udpListenter = nil
}

func main() {

}
