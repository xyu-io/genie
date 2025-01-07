package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestProducer(t *testing.T) {
	d := ProducerDriver{
		KafkaVersion:          "2.8.0",
		KafkaHashing:          false,
		KafkaTLS:              false,
		KafkaSASL:             "none",
		KafkaSrv:              "",
		KafkaBrk:              "192.168.200.222:9092",
		kafkaMaxMsgBytes:      1000000,
		kafkaFlushBytes:       int(sarama.MaxRequestSize),
		kafkaFlushFrequency:   time.Second * 5,
		kafkaCompressionCodec: "",
		producer:              nil,
		q:                     nil,
		errors:                make(chan error),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	d.Init()
	message := "hello world-111"
	topic := "topic-ids"

	for i := 0; i < 3; i++ {
		err := d.Exec(ctx, unsafe.Slice(unsafe.StringData(topic), len(topic)), unsafe.Slice(unsafe.StringData(message), len(message)))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func TestProducer01(t *testing.T) {
	// 配置生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 5                    // 重试次数
	config.Producer.Return.Successes = true          // 成功交付回调

	// 构建Kafka生产者
	brokers := []string{"192.168.200.222:9092"} // 替换为你的Kafka broker地址
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Panicf("Failed to start Sarama producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close producer: %v", err)
		}
	}()

	// 设置消息发送成功的回调
	go func() {
		for msg := range producer.Successes() {
			fmt.Printf("Message sent to topic/partition/offset %s/%d/%d\n", msg.Topic, msg.Partition, msg.Offset)
		}
	}()

	// 设置错误处理回调
	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send message: %v", err)
		}
	}()

	// 发送消息
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		msg := &sarama.ProducerMessage{
			Topic: "topic-ids", // 替换为你的topic名称
			Value: sarama.StringEncoder("Hello, Kafka!"),
		}
		producer.Input() <- msg
	}()

	// 设置信号处理器，用于优雅关闭
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Println("Gracefully shutting down...")

	// 等待所有消息发送完毕
	wg.Wait()
	log.Println("Shutdown complete.")
}
