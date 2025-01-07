package msgchan

import (
	"sync"
)

type MsgChanOption = func(*MsgChan)

func NewMsgChan(opts ...MsgChanOption) *MsgChan {
	instance := &MsgChan{
		lock:      sync.RWMutex{},
		topicList: make(map[string]Msg),
	}
	for _, opt := range opts {
		opt(instance)
	}

	return instance
}

func WithTopic(topic string, isUsed bool) MsgChanOption {
	return func(cfg *MsgChan) {
		if !isUsed {
			if _, ok := cfg.topicList[topic]; ok {
				delete(cfg.topicList, topic)
			}
		} else {
			if _, ok := cfg.topicList[topic]; !ok {
				cfg.topicList[topic] = make(Msg)
			}
		}
	}
}
