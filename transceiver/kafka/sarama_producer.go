package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"genie/iper"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	sarama "github.com/IBM/sarama"
)

type ProducerDriver struct {
	KafkaTLS     bool   `json:"kafka_tls"`
	KafkaSASL    string `json:"kafka_sasl"`
	KafkaSrv     string `json:"kafka_srv"`
	KafkaBrk     string `json:"kafka_brk"`
	KafkaHashing bool   `json:"kafka_hashing"`
	KafkaVersion string `json:"kafka_version"`
	Status       int    `json:"status"`
	AppName      string `json:"app_name"`

	kafkaCompressionCodec string
	kafkaMaxMsgBytes      int
	kafkaFlushBytes       int
	kafkaFlushFrequency   time.Duration

	producer sarama.AsyncProducer

	q chan bool

	errors chan error
}

type KafkaSASLAlgorithm string

const (
	KAFKA_SASL_NONE         KafkaSASLAlgorithm = "none"
	KAFKA_SASL_PLAIN        KafkaSASLAlgorithm = "plain"
	KAFKA_SASL_SCRAM_SHA256 KafkaSASLAlgorithm = "scram-sha256"
	KAFKA_SASL_SCRAM_SHA512 KafkaSASLAlgorithm = "scram-sha512"
)

var (
	compressionCodecs = map[string]sarama.CompressionCodec{
		strings.ToLower(sarama.CompressionNone.String()):   sarama.CompressionNone,
		strings.ToLower(sarama.CompressionGZIP.String()):   sarama.CompressionGZIP,
		strings.ToLower(sarama.CompressionSnappy.String()): sarama.CompressionSnappy,
		strings.ToLower(sarama.CompressionLZ4.String()):    sarama.CompressionLZ4,
		strings.ToLower(sarama.CompressionZSTD.String()):   sarama.CompressionZSTD,
	}

	saslAlgorithms = map[KafkaSASLAlgorithm]bool{
		KAFKA_SASL_PLAIN:        true,
		KAFKA_SASL_SCRAM_SHA256: true,
		KAFKA_SASL_SCRAM_SHA512: true,
	}
	saslAlgorithmsList = []string{
		string(KAFKA_SASL_NONE),
		string(KAFKA_SASL_PLAIN),
		string(KAFKA_SASL_SCRAM_SHA256),
		string(KAFKA_SASL_SCRAM_SHA512),
	}
)

func NewProducerDriver() *ProducerDriver {
	return &ProducerDriver{
		q:      make(chan bool),
		errors: make(chan error),

		kafkaMaxMsgBytes:    1000000,
		kafkaFlushBytes:     int(sarama.MaxRequestSize),
		kafkaFlushFrequency: time.Second * 5,
	}
}

func (d *ProducerDriver) Prepare(conf interface{}) *ProducerDriver {

	if conf == nil {
		return nil
	}

	if cfg, ok := conf.(ProducerDriver); ok {
		if cfg.KafkaVersion == "" {
			cfg.KafkaVersion = "2.8.0"
		}

		d.KafkaTLS = cfg.KafkaTLS
		d.KafkaSASL = cfg.KafkaSASL
		d.KafkaSrv = cfg.KafkaSrv
		d.KafkaBrk = cfg.KafkaBrk
		d.kafkaMaxMsgBytes = 1000000
		d.kafkaFlushBytes = int(sarama.MaxRequestSize)
		d.kafkaFlushFrequency = time.Second * 5
		d.KafkaHashing = cfg.KafkaHashing
		d.kafkaCompressionCodec = ""
		d.KafkaVersion = cfg.KafkaVersion
		d.Status = cfg.Status
		d.AppName = cfg.AppName
		d.q = make(chan bool)
		d.errors = make(chan error)

		return d
	}

	return nil
}

func (d *ProducerDriver) Errors() <-chan error {
	return d.errors
}

func (d *ProducerDriver) Init() (*ProducerDriver, error) {
	kafkaConfigVersion, err := sarama.ParseKafkaVersion(d.KafkaVersion)
	if err != nil {
		return nil, err
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Version = kafkaConfigVersion
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.ClientID = "server-producer"
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.MaxMessageBytes = d.kafkaMaxMsgBytes
	kafkaConfig.Producer.Flush.Bytes = d.kafkaFlushBytes
	kafkaConfig.Producer.Flush.Frequency = d.kafkaFlushFrequency
	kafkaConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	if d.kafkaCompressionCodec != "" {
		/*
			// when upgrading sarama, replace with:
			// note: if the library adds more codecs, they will be supported natively
			var cc *sarama.CompressionCodec

			if err := cc.UnmarshalText([]byte(d.kafkaCompressionCodec)); err != nil {
				return err
			}
			kafkaConfig.Producer.Compression = *cc
		*/

		if cc, ok := compressionCodecs[strings.ToLower(d.kafkaCompressionCodec)]; !ok {
			return nil, fmt.Errorf("compression codec does not exist")
		} else {
			kafkaConfig.Producer.Compression = cc
		}
	}

	if d.KafkaTLS {
		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("error initializing TLS: %v", err)
		}
		kafkaConfig.Net.TLS.Enable = true
		kafkaConfig.Net.TLS.Config = &tls.Config{
			RootCAs:    rootCAs,
			MinVersion: tls.VersionTLS12,
		}
	}

	if d.KafkaHashing {
		kafkaConfig.Producer.Partitioner = sarama.NewHashPartitioner
	}

	kafkaSASL := KafkaSASLAlgorithm(d.KafkaSASL)
	if d.KafkaSASL != "" && kafkaSASL != KAFKA_SASL_NONE {
		_, ok := saslAlgorithms[KafkaSASLAlgorithm(strings.ToLower(d.KafkaSASL))]
		if !ok {
			return nil, errors.New("SASL algorithm does not exist")
		}

		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = os.Getenv("KAFKA_SASL_USER")
		kafkaConfig.Net.SASL.Password = os.Getenv("KAFKA_SASL_PASS")
		if kafkaConfig.Net.SASL.User == "" && kafkaConfig.Net.SASL.Password == "" {
			return nil, fmt.Errorf("Kafka SASL config from environment was unsuccessful. KAFKA_SASL_USER and KAFKA_SASL_PASS need to be set.")
		}

		if kafkaSASL == KAFKA_SASL_SCRAM_SHA256 || kafkaSASL == KAFKA_SASL_SCRAM_SHA512 {
			kafkaConfig.Net.SASL.Handshake = true

			if kafkaSASL == KAFKA_SASL_SCRAM_SHA512 {
				kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
					return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
				}
				kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
			} else if kafkaSASL == KAFKA_SASL_SCRAM_SHA256 {
				kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
					return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
				}
				kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
			}
		}
	}

	var addrs []string
	if d.KafkaSrv != "" {
		addrs, _ = GetServiceAddresses(d.KafkaSrv)
	} else {
		addrs = strings.Split(d.KafkaBrk, ",")
	}

	// 抽离
	kafkaProducer, err := sarama.NewAsyncProducer(addrs, kafkaConfig)
	if err != nil {
		return nil, err
	}

	d.producer = kafkaProducer

	d.q = make(chan bool)

	go func() {
		for {
			select {
			case msg := <-kafkaProducer.Errors():
				var err error
				if msg != nil {
					err = errors.New(fmt.Sprintf("kafka transport %s", msg.Error()))
				}
				select {
				case d.errors <- err:
				default:
				}

				if msg == nil {
					return
				}
			case <-d.q:
				return
			}
		}
	}()

	return d, err
}

// todo: 这里需要区分事件和流量的topic
func (d *ProducerDriver) Exec(ctx context.Context, key, data []byte) error {
	if d.producer == nil {
		return fmt.Errorf("producer not initialized")
	}

	host := iper.LocalIP()
	d.producer.Input() <- &sarama.ProducerMessage{
		Topic: string(key),
		Key:   sarama.ByteEncoder(host),
		Value: sarama.ByteEncoder(data),
	}

	// 设置消息发送成功的回调
	go func() {
		for msg := range d.producer.Successes() {
			fmt.Printf("Message sent to topic/partition/offset %s/%d/%d\n", msg.Topic, msg.Partition, msg.Offset)
		}
	}()
	return nil
}

func (d *ProducerDriver) Close() error {
	if d.producer != nil {
		d.producer.Close()
		close(d.q)
	}
	return nil
}

func GetServiceAddresses(srv string) (addrs []string, err error) {
	_, srvs, err := net.LookupSRV("", "", srv)
	if err != nil {
		return nil, fmt.Errorf("service discovery: %v\n", err)
	}
	for _, srv := range srvs {
		addrs = append(addrs, net.JoinHostPort(srv.Target, strconv.Itoa(int(srv.Port))))
	}
	return addrs, nil
}

func (d *ProducerDriver) GetStatus() int {

	return d.Status
}

func (d *ProducerDriver) UpdateStatus(status int) {
	d.Status = status
}
