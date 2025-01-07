package rudp

import (
	"log"
	"net"
	"sync"
)

// UDPServer 服务端
type UDPServer struct {
	socket *net.UDPConn
	mu     sync.Mutex
}

// NewUDPServer 创建实例,监听端口和IP
func NewUDPServer(ip string, port int) (*UDPServer, error) {
	listenAddr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	listen, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &UDPServer{socket: listen}, nil
}

// Send 发送数据，传入数据和地址
func (s *UDPServer) Send(data []byte, addr *net.UDPAddr) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.socket.WriteToUDP(data, addr)
	return err
}

// Receive 接收数据包
func (s *UDPServer) Receive() ([]byte, *net.UDPAddr, error) {
	var data [1024]byte
	s.mu.Lock()
	defer s.mu.Unlock()
	n, addr, err := s.socket.ReadFromUDP(data[:])
	if err != nil {
		return nil, nil, err
	}
	return data[:n], addr, nil
}

func (s *UDPServer) Close() error {
	err := s.socket.Close()
	if err != nil {
		log.Printf("关闭服务器失败，err: %v", err)
	}
	return err
}
