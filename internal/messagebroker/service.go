package messagebroker

import (
	"time"

	"github.com/car12o/netdata-digest/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type service struct {
	ec *nats.EncodedConn
}

func NewService(url string, log logger.Service) (Service, error) {
	nc, err := nats.Connect(
		url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(15),
		nats.ReconnectWait(time.Second*2),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Error(errors.Wrap(err, "Disconnected from NATS server"))
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info("Successfully reconnected to NATS server")
		}),
	)
	if err != nil {
		return nil, err
	}

	if err := nc.Flush(); err != nil {
		return nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return &service{ec}, nil
}

type topicAlarmStatusChanged struct {
	ec    *nats.EncodedConn
	topic string
	queue string
	ch    chan AlarmStatusChanged
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
		_, err = t.ec.BindRecvQueueChan(t.topic, t.queue, t.ch)
	} else {
		_, err = t.ec.BindRecvChan(t.topic, t.ch)
	}
	if err != nil {
		return err
	}

	go func(ch chan AlarmStatusChanged) {
		for {
			handler(<-ch)
		}
	}(t.ch)

	return nil
}

type topicSendAlarmDigest struct {
	ec    *nats.EncodedConn
	topic string
	queue string
	ch    chan SendAlarmDigest
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
		_, err = t.ec.BindRecvQueueChan(t.topic, t.queue, t.ch)
	} else {
		_, err = t.ec.BindRecvChan(t.topic, t.ch)
	}
	if err != nil {
		return err
	}

	go func(ch chan SendAlarmDigest) {
		for {
			handler(<-ch)
		}
	}(t.ch)

	return nil
}

// type topicAlarmDigest struct {
// 	topic string
// 	ch    chan AlarmDigest
// }

// func (s *service) TopicAlarmDigest() TopicAlarmDigest {

// }

func (s *service) Close() {
	s.ec.Close()
}
