package main

import "C"
import (
	"fmt"
	"net"
	"strings"

	"github.com/aiocloud/shadowsocks/core"
	"github.com/aiocloud/shadowsocks/socks"
)

var (
	ListenAddr string
	RemoteAddr string
	Passwd     string
	Method     string
	OBFS       string
	OBFSParam  string

	cipher        core.Cipher
	tcpListenter  net.Listener
	udpListenter  net.PacketConn
	tcpRemoteAddr net.Addr
	udpRemoteAddr net.Addr
)

const (
	TYPE_LISN int = iota
	TYPE_HOST
	TYPE_PASS
	TYPE_METH
	TYPE_OBFS
	TYPE_OBPA
)

//export ServerInfo
func ServerInfo(name int, value *C.char) bool {
	switch name {
	case TYPE_LISN:
		ListenAddr = C.GoString(value)
	case TYPE_HOST:
		RemoteAddr = C.GoString(value)
	case TYPE_PASS:
		Passwd = C.GoString(value)
	case TYPE_METH:
		Method = C.GoString(value)
	case TYPE_OBFS:
		OBFS = strings.ToUpper(C.GoString(value))
	case TYPE_OBPA:
		OBFSParam = C.GoString(value)
	}

	return true
}

//export Create
func Create() bool {
	var err error

	if cipher, err = core.PickCipher(Method, nil, Passwd); err != nil {
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

	socks.UDPEnabled = true

	go tcpListen()
	go udpListen()
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
