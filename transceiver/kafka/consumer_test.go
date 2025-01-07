package kafka

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestConsumer(t *testing.T) {
	config := ConsumerDriver{
		Brokers:          "192.168.200.222:9093",
		Version:          "2.1.0",
		Group:            "hermes-group",
		CompressionCodec: "",
		TlsEnable:        false,
		SaslEnable:       false,
		SaslConf: MqKafkaSaslConf{
			MaxWaitTime: 250,
		},
		Offsets: MqKafkaOffsets{},
		Status:  0,
	}

	ctx, cancel := context.WithCancel(context.Background())
	driver, err := NewConsumerDriver().Prepare(&config).Init()
	if err != nil {
		t.Fatal(err)
	}
	var topics = []string{"test-kafka", "topic-ids"}
	tpBytes, _ := json.Marshal(topics)

	go driver.Exec(ctx, []byte{}, tpBytes)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigterm:
			ctx.Done()
			cancel()
		}
	}

}
