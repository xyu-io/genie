package rudp

import (
	"log"
	"net"
	"time"
)

// UDPClient 客户端
type UDPClient struct {
	socket          *net.UDPConn
	timeoutDuration time.Duration
}

func NewUDPClient(ip string, port int) (*UDPClient, error) {
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	//localAddr := &net.UDPAddr{
	//	IP:   net.ParseIP(clientIp),
	//	Port: clientPort,
	//}
	socket, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}
	client := &UDPClient{socket: socket}
	// 设置超时时间为 5 秒
	client.timeoutDuration = 5 * time.Second
	err = client.socket.SetDeadline(time.Now().Add(client.timeoutDuration))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *UDPClient) Send(data []byte) error {
	_, err := c.socket.Write(data)
	return err
}

func (c *UDPClient) Receive() ([]byte, error) {
	data := make([]byte, 4096)
	n, _, err := c.socket.ReadFromUDP(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (c *UDPClient) Close() error {
	err := c.socket.Close()
	if err != nil {
		log.Printf("关闭客户端失败，err: %v", err)
	}
	return err
}
