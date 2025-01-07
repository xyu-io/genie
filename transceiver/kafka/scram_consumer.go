package kafka

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	log2 "github.com/sirupsen/logrus"
)

// Filter 增加过滤器
type Filter func(message []byte) bool

// ReceiverFn 流量接收
type ReceiverFn func(topic string, host []byte, message []byte, filter ...Filter) error

type ConsumerDriver struct {
	Brokers          string          `json:"brokers"`
	Version          string          `json:"version"`
	Group            string          `json:"group"`
	CompressionCodec string          `json:"compression_codec"`
	TlsEnable        bool            `json:"tls_enable"`
	SaslEnable       bool            `json:"sasl_enable"`
	SaslConf         MqKafkaSaslConf `json:"sasl_conf"`
	Offsets          MqKafkaOffsets  `json:"offsets"`
	ReceiverFn

	consumer      sarama.Consumer
	consumerGroup sarama.ConsumerGroup

	q       chan bool
	errors  chan error
	AppName string
	Status  int
}

func (c *ConsumerDriver) Send(topic string, key, data []byte) error {
	return nil
}

func (c *ConsumerDriver) Close() error {
	c.consumerGroup.Close()

	return nil
}

type MqKafkaSaslConf struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Mechanism   string `json:"mechanism"`
	MaxWaitTime int    `json:"max_waitTime"`
}

type MqKafkaOffsets struct {
	AutocommitEnable bool `json:"autocommit_enable"`
	Interval         int  `json:"interval"`
	Retry            int  `json:"retry"`
}

func NewConsumerDriver() *ConsumerDriver {
	return &ConsumerDriver{
		q:             make(chan bool),
		errors:        make(chan error),
		SaslConf:      MqKafkaSaslConf{},
		Offsets:       MqKafkaOffsets{},
		consumerGroup: nil,
		consumer:      nil,
	}
}

func (c *ConsumerDriver) Prepare(conf interface{}) *ConsumerDriver {
	if conf == nil {
		return nil
	}
	if cfg, ok := conf.(ConsumerDriver); ok {
		if cfg.Version == "" {
			cfg.Version = "2.8.0"
		}
		c.Brokers = cfg.Brokers
		c.Version = cfg.Version
		c.Group = cfg.Group
		c.CompressionCodec = cfg.CompressionCodec
		c.TlsEnable = cfg.TlsEnable
		c.SaslEnable = cfg.SaslEnable
		c.SaslConf = cfg.SaslConf
		c.AppName = cfg.AppName
		c.Status = cfg.Status

		return c
	}

	return nil
}

func (c *ConsumerDriver) Init() (*ConsumerDriver, error) {

	//	endpoint, _ := os.Hostname() //"cn-beijing.log.aliyuncs.com"
	version := c.Version // "2.1.0"
	group := c.Group     //"test-groupId"

	log2.Info("Starting a new Sarama consumer")
	newVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log2.Errorf("Error parsing Kafka version: %v", err)
	}

	/**
	 	 * 构建一个新的Sarama配置。
		 * 在初始化消费者/生产者之前，必须定义Kafka集群版本。
	*/
	brokers := []string{c.Brokers} // []string{fmt.Sprintf("%s.%s:%s", project, endpoint, port)}

	conf := sarama.NewConfig()
	conf.Version = newVersion

	conf.Net.TLS.Enable = c.TlsEnable
	conf.Net.SASL.Enable = c.SaslEnable
	conf.Net.SASL.User = c.SaslConf.Username
	conf.Net.SASL.Password = c.SaslConf.Password
	conf.Net.SASL.Mechanism = sarama.SASLMechanism(c.SaslConf.Mechanism) //"PLAIN"

	conf.Consumer.MaxWaitTime = time.Duration(c.SaslConf.MaxWaitTime) * time.Millisecond //250 * time.Millisecond
	conf.Consumer.MaxProcessingTime = 100 * time.Millisecond
	conf.Consumer.Fetch.Min = 1
	conf.Consumer.Fetch.Default = 1024 * 1024
	conf.Consumer.Retry.Backoff = 2 * time.Second
	conf.Consumer.Return.Errors = false
	conf.Consumer.Offsets.AutoCommit.Enable = true
	conf.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Consumer.Offsets.Retry.Max = 3

	client, err := sarama.NewConsumerGroup(brokers, group, conf)
	if err != nil {
		log2.Errorf("Error creating consumer group client: %v", err)
		return nil, fmt.Errorf("error creating consumer group client: %s", err.Error())
	}

	c.consumerGroup = client

	return c, nil
}

func (c *ConsumerDriver) Exec(ctx context.Context, key, data []byte) error {
	if c.consumerGroup == nil {
		log2.Errorf("consumer group not initialized")
		return errors.New("consumer group not initialized")
	}

	var topic []string
	json.Unmarshal(data, &topic)

	var fn = func(tp string, key, value []byte) error {
		if c.ReceiverFn != nil {
			err := c.ReceiverFn(tp, key, value)
			if err != nil {
				log2.Errorf("Error on receiver function: %v", err)
				return err
			}
		}
		return nil
	}

	keepRunning := true
	/**
	 * 设置一个新的Sarama消费者组
	 */
	consumer := Consumer{
		ready: make(chan bool),
		Send:  fn,
	}

	// ctx, cancel := context.WithCancel(context.Background())
	// 是否暂停信号
	// consumptionIsPaused := false

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume`应该在一个无限循环内调用，当服务器端重新平衡时，消费者会话将需要重新创建以获取新的声明
			if err := c.consumerGroup.Consume(ctx, topic, &consumer); err != nil {
				log2.Errorf("Error from consumer: %v", err)
			}
			// 检查上下文是否被取消，表示消费者应该停止
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // 等待消费者设置完成
	log2.Info("Sarama consumer up and running!...")

	// 是否暂停消费/恢复消费
	// sigusr1 := make(chan os.Signal, 1)
	// 退出信号，ctl+c
	//sigterm := make(chan os.Signal, 1)
	//signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
			//case <-sigterm:
			//	log.Println("terminating: via signal")
			//	keepRunning = false
			//	//case <-sigusr1:
			//	//	toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	// cancel()
	wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		log2.Errorf("Error closing client: %v", err)
	}
	log2.Info("Sarama consumer stop and close!")

	return nil
}

// 暂停与恢复消费
func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

type Send func(string, []byte, []byte) error

// Consumer 表示Sarama消费者组消费者
type Consumer struct {
	ready chan bool
	Send
}

// Setup在新会话开始时运行，在ConsumeClaim之前
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// 将消费者标记为已准备好
	close(consumer.ready)
	return nil
}

// Cleanup在会话结束时运行，一旦所有ConsumeClaim goroutine退出
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim必须启动ConsumerGroupClaim的Messages()的消费者循环
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 注意：
	// 不要把下面的代码移到goroutine里。
	// ConsumeClaim本身在goroutine中调用，参见：https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				continue
			}
			realUnixTimeSeconds := message.Timestamp.Unix()
			if realUnixTimeSeconds < 2000000 {
				realUnixTimeSeconds = message.Timestamp.UnixMicro() / 1000
			}
			err := consumer.Send(message.Topic, message.Key, message.Value)
			if err != nil {
				log2.Errorf("Error sending message: %v", err)
			}
			// log2.Info("Message claimed: key = %s, value = %s, timestamp = %d, topic = %s", string(message.Key), string(message.Value), realUnixTimeSeconds, message.Topic)
			session.MarkMessage(message, "")

		// 当session.Context()完成时应返回。
		// 如果不这样做，当kafka重新平衡时，会引发ErrRebalanceInProgress或read tcp <ip>:<port>: i/o timeout错误。参见：https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *ConsumerDriver) GetStatus() int {

	return c.Status
}

func (c *ConsumerDriver) UpdateStatus(status int) {
	c.Status = status

}
