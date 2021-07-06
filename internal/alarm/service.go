package alarm

import (
	"fmt"
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
	s.log.Debug(fmt.Sprintf("topicAlarmStatusChangedHandler: %v", msg))

	// Since CLEARED status is a non active status I'm simply discarding it.
	// It was not clear for me if this is the intended behavior or if it was to clear any
	// previous record for the specified alarmID.
	// However, change this to the clear logic would be just a matter of calling:
	// s.repository.DeleteByID(msg.UserID, msg.AlarmID)
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
	s.log.Debug(fmt.Sprintf("topicSendAlarmDigestHandler: %v", msg))

	alarms, err := s.repository.GetByUserID(msg.UserID)
	if err != nil {
		s.log.Error(errors.Wrap(err, "error fetching all alarms"))
		return
	}
	if len(alarms) == 0 {
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

	if err := s.repository.DeleteByID(msg.UserID, alarmsToRemove...); err != nil {
		s.log.Error(errors.Wrap(err, "error deleting alarms"))
	}
}
