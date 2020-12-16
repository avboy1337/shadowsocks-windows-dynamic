package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func udpListen() {
	nm := newNAT()
	buffer := make([]byte, 1500)

	for {
		size, from, err := udpListenter.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("[shadowsocks][udpListenter.ReadFrom] %v", err)
			return
		}

		remote := nm.Get(from.String())
		if remote == nil {
			remote, err = net.ListenPacket("udp", "")
			if err != nil {
				fmt.Printf("[shadowsocks][net.ListenPacket] %v", err)
				continue
			}
			remote = cipher.PacketConn(remote)

			fmt.Printf("[shadowsocks] New UDP connection from %s to %s", from, udpRemoteAddr)
			nm.Add(from, udpListenter, remote)
		}

		if _, err = remote.WriteTo(buffer[3:size], udpRemoteAddr); err != nil {
			continue
		}
	}
}

type NAT struct {
	sync.RWMutex

	m map[string]net.PacketConn
}

func (m *NAT) Get(key string) net.PacketConn {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *NAT) Set(key string, pc net.PacketConn) {
	m.Lock()
	defer m.Unlock()

	m.m[key] = pc
}

func (m *NAT) Del(key string) net.PacketConn {
	m.Lock()
	defer m.Unlock()

	pc, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return pc
	}
	return nil
}

func (m *NAT) Add(peer net.Addr, dst, src net.PacketConn) {
	m.Set(peer.String(), src)

	go func() {
		timedCopy(dst, peer, src)
		if pc := m.Del(peer.String()); pc != nil {
			pc.Close()
		}
	}()
}

func newNAT() *NAT {
	return &NAT{m: make(map[string]net.PacketConn)}
}

func timedCopy(dst net.PacketConn, target net.Addr, src net.PacketConn) {
	buffer := make([]byte, 1500)

	for {
		_ = src.SetReadDeadline(time.Now().Add(time.Second * 120))

		size, _, err := src.ReadFrom(buffer)
		if err != nil {
			return
		}

		if _, err = dst.WriteTo(append([]byte{0, 0, 0}, buffer[:size]...), target); err != nil {
			return
		}
	}
}
