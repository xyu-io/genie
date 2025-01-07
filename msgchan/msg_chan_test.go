package msgchan

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

// PubFunc | SubFunc : Service interface
type PubFunc func(string, ...[]byte) error
type SubFunc func(topic string, handler func([]byte) error)

func TestMsgChan(t *testing.T) {
	var msgChanIns = NewMsgChan(
		WithTopic("tp", true),
		WithTopic("tp1", true),
	)
	_, err := msgChanIns.CreatePublisher(nil)
	if err != nil {
		return
	}

	go func() {
		for i := 0; i < 10000; i++ {
			if i%2 == 0 {
				err := msgChanIns.PublishMsg("tp", []byte("hello world 1111111"))
				if err != nil {
					continue
				}
			} else {
				err := msgChanIns.PublishMsg("tp1", []byte("hello world 0000001"))
				if err != nil {
					continue
				}
			}
		}
	}()

	msgChanIns.Subscribe("tp", func(bytes []byte) error {
		log.Info("topic tp >> received bytes:", string(bytes))
		return nil
	})

	msgChanIns.Subscribe("tp1", func(bytes []byte) error {
		log.Info("topic tp1 >> received bytes:", string(bytes))
		return nil
	})

	time.Sleep(time.Second * time.Duration(5))
}

// ------------------ with obj ---------------- //

// Email obj
type SendEmail struct {
	sender PubFunc
}

func (se *SendEmail) Send(topic string, msg []byte) error {
	return se.sender(topic, msg)
}

func (se *SendEmail) Run() {
	go func() {
		for i := 0; i < 10; i++ {
			err := se.Send("email", []byte("this is a 【email】 msg"))
			if err != nil {
				continue
			}
		}
	}()
}

// Notice obj
type SendNotice struct {
	sender PubFunc
}

func (sn *SendNotice) Send(topic string, msg []byte) error {
	return sn.sender(topic, msg)
}

func (sn *SendNotice) Run() {
	go func() {
		for i := 0; i < 10; i++ {
			err := sn.Send("notice", []byte("this is a 【notice】 msg"))
			if err != nil {
				continue
			}
		}
	}()
}

// Test
func TestMsgChanOfObj(t *testing.T) {
	var msgChanIns = NewMsgChan(
		WithTopic("email", true),
		WithTopic("notice", true),
	)
	_, err := msgChanIns.CreatePublisher(nil)
	if err != nil {
		return
	}

	var emailIns SendEmail
	emailIns.sender = msgChanIns.PublishMsg

	var noticeIns SendNotice
	noticeIns.sender = msgChanIns.PublishMsg

	emailIns.Run()
	noticeIns.Run()

	msgChanIns.Subscribe("email", func(bytes []byte) error {
		log.Info("【topic email】 >> received bytes:", string(bytes))
		return nil
	})

	msgChanIns.Subscribe("notice", func(bytes []byte) error {
		log.Info("【topic notice】 >> received bytes:", string(bytes))
		return nil
	})

	time.Sleep(time.Second * time.Duration(5))
}
