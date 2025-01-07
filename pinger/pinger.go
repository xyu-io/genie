package pinger

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)

func Ping(address string) (bool, error) {
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false, err
	}
	defer conn.Close()

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello, world!"),
		},
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return false, err
	}

	dstAddr, err := net.ResolveIPAddr("ip4", address)
	if err != nil {
		return false, err
	}

	conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
	if _, err := conn.WriteTo(msgBytes, dstAddr); err != nil {
		return false, err
	}

	var buffer [1500]byte
	conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
	n, _, err := conn.ReadFrom(buffer[0:])
	if err != nil {
		return false, err
	}

	msg, err = icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buffer[0:n])
	if err != nil {
		return false, err
	}
	if msg.Type != ipv4.ICMPTypeEchoReply {
		return false, fmt.Errorf("unexpected message type: %v", msg.Type)
	}

	return true, nil
}
