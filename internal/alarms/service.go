package alarms

import (
	"fmt"

	"github.com/car12o/netdata-digest/internal/messagebroker"
)

const queue = "netdata-digest"

type service struct {
	messageBroker messagebroker.Service
}

func NewService(mb messagebroker.Service) Service {
	return &service{mb}
}

func (s *service) SubscribeTopics() error {
	if err := s.messageBroker.TopicAlarmStatusChanged().WithQueue(queue).Subscribe(
		func(msg messagebroker.AlarmStatusChanged) {
			fmt.Println("## AlarmStatusChanged ##", "AlarmID ", msg.AlarmID, " UserID ", msg.UserID, " Status ", msg.Status, " ChangedAt ", msg.ChangedAt)
		},
	); err != nil {
		return err
	}

	if err := s.messageBroker.TopicSendAlarmDigest().WithQueue(queue).Subscribe(
		func(msg messagebroker.SendAlarmDigest) {
			fmt.Println("## SendAlarmDigest ##", "UserID ", msg.UserID)
		},
	); err != nil {
		return err
	}

	return nil
}
