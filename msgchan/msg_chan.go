package msgchan

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	log "github.com/sirupsen/logrus"
	"sync"
)

type (
	Msg     <-chan *message.Message
	MsgChan struct {
		lock      sync.RWMutex
		producer  message.Publisher
		topicList map[string]Msg
	}
)

func (m *MsgChan) CreatePublisher(opt *gochannel.Config) (message.Publisher, error) {
	config := gochannel.Config{}
	if opt != nil {
		config = *opt
	}

	pub := gochannel.NewGoChannel(
		config,
		watermill.NewStdLogger(false, false),
	)

	m.producer = pub

	for topic, _ := range m.topicList {
		m.topicList[topic], _ = pub.Subscribe(context.Background(), topic)
	}

	return pub, nil
}

func (m *MsgChan) PublishMsg(topic string, msg ...[]byte) error {
	if msg == nil || len(msg) == 0 {
		return errors.New("msg list is nil or empty")
	}
	var messageList []*message.Message
	for _, bytes := range msg {
		if bytes == nil || len(bytes) == 0 {
			log.Warn("msg is nil or empty")
			continue
		}
		messageList = append(messageList, message.NewMessage(watermill.NewUUID(), bytes))
	}

	if len(messageList) == 0 {
		return errors.New("messageList is empty")
	}

	return m.producer.Publish(topic, messageList...)
}

func (m *MsgChan) Subscribe(topic string, handler func([]byte) error) {
	if _, ok := m.topicList[topic]; !ok {
		log.Warn("topic not found")
		return
	}

	go func() {
		for {
			select {
			case msg := <-m.topicList[topic]:
				err := handler(msg.Payload)
				if err != nil {
					log.Errorf("handle message, uuid: %s error: %s", msg.UUID, err.Error())
					break
				}
				msg.Ack()
			default:

			}
		}
	}()
}
