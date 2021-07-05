package alarm

import (
	"fmt"

	"github.com/car12o/netdata-digest/internal/messenger"
)

const queue = "netdata-digest"

type service struct {
	// repository Repository
	messenger messenger.Service
}

// func NewService(repository Repository, messenger messenger.Service) Service {
// 	return &service{repository, messenger}
// }
func NewService(messenger messenger.Service) Service {
	return &service{messenger}
}

func (s *service) SubscribeTopics() error {
	if err := s.messenger.TopicAlarmStatusChanged().WithQueue(queue).Subscribe(
		func(msg messenger.AlarmStatusChanged) {
			fmt.Println("## AlarmStatusChanged ##", "AlarmID ", msg.AlarmID, " UserID ", msg.UserID, " Status ", msg.Status, " ChangedAt ", msg.ChangedAt)
		},
	); err != nil {
		return err
	}

	if err := s.messenger.TopicSendAlarmDigest().WithQueue(queue).Subscribe(
		func(msg messenger.SendAlarmDigest) {
			fmt.Println("## SendAlarmDigest ##", "UserID ", msg.UserID)
		},
	); err != nil {
		return err
	}

	return nil
}
