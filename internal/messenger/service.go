package messenger

import (
	"github.com/nats-io/nats.go"
)

type service struct {
	ec *nats.EncodedConn
}

func NewService(ec *nats.EncodedConn) Service {
	return &service{ec}
}

type topicAlarmStatusChanged struct {
	ec           *nats.EncodedConn
	subscription *nats.Subscription
	topic        string
	queue        string
	ch           chan AlarmStatusChanged
}

func (s *service) TopicAlarmStatusChanged() TopicAlarmStatusChanged {
	return &topicAlarmStatusChanged{
		ec:    s.ec,
		topic: "AlarmStatusChanged",
	}
}

func (t *topicAlarmStatusChanged) WithQueue(queue string) TopicAlarmStatusChanged {
	t.queue = queue
	return t
}

func (t *topicAlarmStatusChanged) Subscribe(handler func(msg AlarmStatusChanged)) error {
	t.ch = make(chan AlarmStatusChanged)

	var err error
	if t.queue != "" {
		t.subscription, err = t.ec.BindRecvQueueChan(t.topic, t.queue, t.ch)
	} else {
		t.subscription, err = t.ec.BindRecvChan(t.topic, t.ch)
	}
	if err != nil {
		return err
	}

	go func(ch chan AlarmStatusChanged) {
		for msg := range ch {
			handler(msg)
		}
	}(t.ch)

	return nil
}

type topicSendAlarmDigest struct {
	ec           *nats.EncodedConn
	subscription *nats.Subscription
	topic        string
	queue        string
	ch           chan SendAlarmDigest
}

func (s *service) TopicSendAlarmDigest() TopicSendAlarmDigest {
	return &topicSendAlarmDigest{
		ec:    s.ec,
		topic: "SendAlarmDigest",
	}
}

func (t *topicSendAlarmDigest) WithQueue(queue string) TopicSendAlarmDigest {
	t.queue = queue
	return t
}

func (t *topicSendAlarmDigest) Subscribe(handler func(msg SendAlarmDigest)) error {
	t.ch = make(chan SendAlarmDigest)

	var err error
	if t.queue != "" {
		t.subscription, err = t.ec.BindRecvQueueChan(t.topic, t.queue, t.ch)
	} else {
		t.subscription, err = t.ec.BindRecvChan(t.topic, t.ch)
	}
	if err != nil {
		return err
	}

	go func(ch chan SendAlarmDigest) {
		for msg := range ch {
			handler(msg)
		}
	}(t.ch)

	return nil
}

type topicAlarmDigest struct {
	ec    *nats.EncodedConn
	topic string
	ch    chan AlarmDigest
}

func (s *service) TopicAlarmDigest() (TopicAlarmDigest, error) {
	t := &topicAlarmDigest{
		ec:    s.ec,
		topic: "AlarmDigest",
		ch:    make(chan AlarmDigest),
	}

	if err := t.ec.BindSendChan(t.topic, t.ch); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *topicAlarmDigest) Publish(msg AlarmDigest) {
	t.ch <- msg
}
