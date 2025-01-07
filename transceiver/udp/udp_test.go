package rudp

import (
	"sync"
	"testing"
	"time"
)

func TestUdp(t *testing.T) {
	localIP := "127.0.0.1"
	port := 8100
	//clientIp := "127.0.0.1"
	//clientPort := 8101
	t.Run("TestServerReceive", func(t *testing.T) {
		udpServer, err := NewUDPServer(localIP, port)
		if err != nil {
			t.Error(err)
			return
		}
		defer udpServer.Close()

		receivedCount := 0
		maxReceivedCount := 3 // 设置接收次数上限

		for receivedCount < maxReceivedCount {
			data, addr, err := udpServer.Receive()
			if err != nil {
				t.Error(err)
				continue
			}
			t.Logf("接收到的数据：%s ------服务器地址：%s", data, addr)
			receivedCount++
		}
	})
	t.Run("TestClientSend", func(t *testing.T) {
		udpClient, err := NewUDPClient(localIP, port)
		if err != nil {
			t.Error(err)
			return
		}
		defer udpClient.Close()

		messages := []string{
			"这是第一次发送的数据",
			"这是第二次发送的数据",
			"这是第三次发送的数据",
		}

		for _, message := range messages {
			err = udpClient.Send([]byte(message))
			if err != nil {
				t.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	})

}
func TestUdp2(t *testing.T) {
	localIP := "127.0.0.1"
	port := 8100
	//clientIp := "127.0.0.1"
	//clientPort := 8101

	// 启动 UDP 服务器
	udpServer, err := NewUDPServer(localIP, port)
	if err != nil {
		t.Fatal(err)
	}
	defer udpServer.Close()

	// 启动 UDP 客户端
	udpClient, err := NewUDPClient(localIP, port)
	if err != nil {
		t.Fatal(err)
	}
	defer udpClient.Close()

	// 定义要发送的消息
	messages := []string{
		"这是第一次发送的数据",
		"这是第二次发送的数据",
		"这是第三次发送的数据",
	}

	// 启动 goroutine 用于接收数据
	var wg sync.WaitGroup
	receivedCount := 0
	maxReceivedCount := len(messages) // 设置接收次数上限

	wg.Add(1)
	go func() {
		defer wg.Done()
		for receivedCount < maxReceivedCount {
			data, addr, err := udpServer.Receive()
			if err != nil {
				t.Error(err)
				continue
			}
			t.Logf("接收到的数据：%s ------服务器地址：%s", data, addr)
			receivedCount++
		}
	}()

	// 发送消息
	for _, message := range messages {
		err = udpClient.Send([]byte(message))
		if err != nil {
			t.Error(err)
		}
		time.Sleep(1 * time.Second) // 等待一段时间，确保数据发送
	}

	// 等待接收 goroutine 完成
	wg.Wait()

	// 确保所有消息都已接收
	if receivedCount != maxReceivedCount {
		t.Errorf("Expected %d messages to be received, but got %d", maxReceivedCount, receivedCount)
	}
}
