package alarm

import (
	"sort"

	"github.com/car12o/netdata-digest/internal/messenger"
	"github.com/car12o/netdata-digest/pkg/logger"
	"github.com/pkg/errors"
)

const (
	StatusCleared  = "CLEARED"
	StatusWarning  = "WARNING"
	StatusCritical = "CRITICAL"
	queue          = "netdata-digest"
)

type service struct {
	log        logger.Service
	repository Repository
	messenger  messenger.Service
	publishers struct {
		alarmDigest messenger.TopicAlarmDigest
	}
}

func NewService(repository Repository, messenger messenger.Service, log logger.Service) Service {
	return &service{
		repository: repository,
		messenger:  messenger,
		log:        log,
	}
}

func (s *service) TopicsSubscribe() error {
	if err := s.messenger.TopicAlarmStatusChanged().WithQueue(queue).Subscribe(
		s.topicAlarmStatusChangedHandler,
	); err != nil {
		return err
	}

	if err := s.messenger.TopicSendAlarmDigest().WithQueue(queue).Subscribe(
		s.topicSendAlarmDigestHandler,
	); err != nil {
		return err
	}

	return nil
}

func (s *service) TopicsInitPublishers() error {
	var err error
	s.publishers.alarmDigest, err = s.messenger.TopicAlarmDigest()
	return err
}

func (s *service) topicAlarmStatusChangedHandler(msg messenger.AlarmStatusChanged) {
	if msg.Status == StatusCleared {
		return
	}

	alarm, err := s.repository.GetByID(msg.UserID, msg.AlarmID)
	if err != nil {
		s.log.Error(errors.Wrap(err, "error fetching alarm"))
		return
	}

	if alarm != nil && alarm.ChangedAt.After(msg.ChangedAt) {
		return
	}

	if err := s.repository.Store(Alarm{
		ID:        msg.AlarmID,
		UserID:    msg.UserID,
		Status:    msg.Status,
		ChangedAt: msg.ChangedAt,
	}); err != nil {
		s.log.Error(errors.Wrap(err, "error storing alarm"))
	}
}

func (s *service) topicSendAlarmDigestHandler(msg messenger.SendAlarmDigest) {
	alarms, err := s.repository.GetByUserID(msg.UserID)
	if err != nil {
		s.log.Error(errors.Wrap(err, "error fetching all alarms"))
		return
	}

	sort.Slice(alarms, func(i, j int) bool {
		return alarms[i].ChangedAt.Before(alarms[j].ChangedAt)
	})

	activeAlarms := make([]*messenger.ActiveAlarm, len(alarms))
	alarmsToRemove := make([]string, len(alarms))
	for i, alarm := range alarms {
		activeAlarms[i] = &messenger.ActiveAlarm{
			AlarmID:         alarm.ID,
			Status:          alarm.Status,
			LatestChangedAt: alarm.ChangedAt,
		}
		alarmsToRemove[i] = alarm.ID
	}

	s.publishers.alarmDigest.Publish(
		messenger.AlarmDigest{
			UserID:       msg.UserID,
			ActiveAlarms: activeAlarms,
		},
	)

	s.repository.DeleteByID(msg.UserID, alarmsToRemove...)
}
