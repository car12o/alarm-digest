package alarm

import (
	"testing"
	"time"

	"github.com/car12o/netdata-digest/internal/messenger"
	"github.com/car12o/netdata-digest/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopicsSubscribe(t *testing.T) {
	messengerService := new(messenger.MockService)
	topicAlarmStatusChanged := new(messenger.MockTopicAlarmStatusChanged)
	topicSendAlarmDigest := new(messenger.MockTopicSendAlarmDigest)
	messengerService.On("TopicAlarmStatusChanged").Return(topicAlarmStatusChanged)
	messengerService.On("TopicSendAlarmDigest").Return(topicSendAlarmDigest)
	topicAlarmStatusChanged.On("WithQueue", queue).Return(topicAlarmStatusChanged)
	topicAlarmStatusChanged.On("Subscribe", mock.Anything).Return(nil)
	topicSendAlarmDigest.On("WithQueue", queue).Return(topicSendAlarmDigest)
	topicSendAlarmDigest.On("Subscribe", mock.Anything).Return(nil)

	s := NewService(&MockRepository{}, messengerService, &logger.MockService{})
	err := s.TopicsSubscribe()
	assert.Nil(t, err)

	messengerService.AssertNumberOfCalls(t, "TopicAlarmStatusChanged", 1)
	messengerService.AssertNumberOfCalls(t, "TopicSendAlarmDigest", 1)
	topicAlarmStatusChanged.AssertNumberOfCalls(t, "WithQueue", 1)
	topicAlarmStatusChanged.AssertNumberOfCalls(t, "Subscribe", 1)
	topicSendAlarmDigest.AssertNumberOfCalls(t, "WithQueue", 1)
	topicSendAlarmDigest.AssertNumberOfCalls(t, "Subscribe", 1)
}

func TestTopicsInitPublishers(t *testing.T) {
	messengerService := new(messenger.MockService)
	topicAlarmDigest := new(messenger.MockTopicAlarmDigest)
	messengerService.On("TopicAlarmDigest").Return(topicAlarmDigest, nil)

	s := NewService(&MockRepository{}, messengerService, &logger.MockService{})
	err := s.TopicsInitPublishers()
	assert.Nil(t, err)

	messengerService.AssertNumberOfCalls(t, "TopicAlarmDigest", 1)
}

func TestTopicAlarmStatusChangedHandler(t *testing.T) {
	t.Run("new unknown alarm", func(t *testing.T) {
		msg := messenger.AlarmStatusChanged{
			AlarmID:   "alarmID-test",
			UserID:    "userID-test",
			Status:    StatusCritical,
			ChangedAt: time.Now(),
		}

		repository := new(MockRepository)
		repository.On("GetByID", msg.UserID, msg.AlarmID).Return(nil, nil)
		repository.On("Store", mock.Anything).Return(nil)

		s := service{repository: repository}
		s.topicAlarmStatusChangedHandler(msg)

		repository.AssertNumberOfCalls(t, "GetByID", 1)
		repository.AssertNumberOfCalls(t, "Store", 1)
		repository.AssertCalled(t, "Store", Alarm{
			ID:        msg.AlarmID,
			UserID:    msg.UserID,
			Status:    msg.Status,
			ChangedAt: msg.ChangedAt,
		})
	})

	t.Run("order delivery most recent", func(t *testing.T) {
		msg := messenger.AlarmStatusChanged{
			AlarmID:   "alarmID-test",
			UserID:    "userID-test",
			Status:    StatusWarning,
			ChangedAt: time.Now(),
		}

		repository := new(MockRepository)
		repository.On("GetByID", msg.UserID, msg.AlarmID).Return(
			&Alarm{
				ID:        msg.AlarmID,
				UserID:    msg.UserID,
				Status:    msg.Status,
				ChangedAt: msg.ChangedAt.Add(-time.Second * 10),
			},
			nil,
		)
		repository.On("Store", mock.Anything).Return(nil)

		s := service{repository: repository}
		s.topicAlarmStatusChangedHandler(msg)

		repository.AssertNumberOfCalls(t, "GetByID", 1)
		repository.AssertNumberOfCalls(t, "Store", 1)
		repository.AssertCalled(t, "Store", Alarm{
			ID:        msg.AlarmID,
			UserID:    msg.UserID,
			Status:    msg.Status,
			ChangedAt: msg.ChangedAt,
		})
	})

	t.Run("order delivery outdated", func(t *testing.T) {
		msg := messenger.AlarmStatusChanged{
			AlarmID:   "alarmID-test",
			UserID:    "userID-test",
			Status:    StatusCritical,
			ChangedAt: time.Now(),
		}

		repository := new(MockRepository)
		repository.On("GetByID", msg.UserID, msg.AlarmID).Return(
			&Alarm{
				ID:        msg.AlarmID,
				UserID:    msg.UserID,
				Status:    msg.Status,
				ChangedAt: msg.ChangedAt.Add(time.Second * 10),
			},
			nil,
		)
		repository.On("Store", mock.Anything).Return(nil)

		s := service{repository: repository}
		s.topicAlarmStatusChangedHandler(msg)

		repository.AssertNumberOfCalls(t, "GetByID", 1)
		repository.AssertNotCalled(t, "Store")
	})

	t.Run("non active alarm", func(t *testing.T) {
		msg := messenger.AlarmStatusChanged{
			AlarmID:   "alarmID-test",
			UserID:    "userID-test",
			Status:    StatusCleared,
			ChangedAt: time.Now(),
		}

		repository := new(MockRepository)
		repository.On("GetByID", msg.UserID, msg.AlarmID).Return(nil, nil)
		repository.On("Store", mock.Anything).Return(nil)

		s := service{repository: repository}
		s.topicAlarmStatusChangedHandler(msg)

		repository.AssertNotCalled(t, "GetByID")
		repository.AssertNotCalled(t, "Store")
	})
}

func TestTopicSendAlarmDigestHandler(t *testing.T) {
	t.Run("active alarms order", func(t *testing.T) {
		msg := messenger.SendAlarmDigest{
			UserID: "userID-test",
		}

		now := time.Now()
		alarms := []*Alarm{
			{
				ID:        "test-2",
				UserID:    "userID-test",
				Status:    StatusCritical,
				ChangedAt: now.Add(-time.Second * 10),
			},
			{
				ID:        "test-3",
				UserID:    "userID-test",
				Status:    StatusWarning,
				ChangedAt: now,
			},
			{
				ID:        "test-1",
				UserID:    "userID-test",
				Status:    StatusCritical,
				ChangedAt: now.Add(-time.Second * 20),
			},
		}

		repository := new(MockRepository)
		repository.On("GetByUserID", msg.UserID).Return(alarms, nil)
		repository.On("DeleteByID", msg.UserID, mock.Anything).Return(nil)

		topicAlarmDigest := new(messenger.MockTopicAlarmDigest)
		topicAlarmDigest.On("Publish", mock.Anything)

		s := service{
			repository: repository,
			publishers: struct{ alarmDigest messenger.TopicAlarmDigest }{
				alarmDigest: topicAlarmDigest,
			},
		}
		s.topicSendAlarmDigestHandler(msg)

		repository.AssertNumberOfCalls(t, "GetByUserID", 1)
		repository.AssertNumberOfCalls(t, "DeleteByID", 1)
		repository.AssertCalled(t, "DeleteByID", msg.UserID, []string{"test-1", "test-2", "test-3"})
		topicAlarmDigest.AssertNumberOfCalls(t, "Publish", 1)
		topicAlarmDigest.AssertCalled(t, "Publish", messenger.AlarmDigest{
			UserID: msg.UserID,
			ActiveAlarms: []*messenger.ActiveAlarm{
				{
					AlarmID:         "test-1",
					Status:          StatusCritical,
					LatestChangedAt: now.Add(-time.Second * 20),
				},
				{
					AlarmID:         "test-2",
					Status:          StatusCritical,
					LatestChangedAt: now.Add(-time.Second * 10),
				},
				{
					AlarmID:         "test-3",
					Status:          StatusWarning,
					LatestChangedAt: now,
				},
			},
		})
	})

	t.Run("no alarms to digest", func(t *testing.T) {
		msg := messenger.SendAlarmDigest{
			UserID: "userID-test",
		}

		repository := new(MockRepository)
		repository.On("GetByUserID", msg.UserID).Return([]*Alarm{}, nil)
		repository.On("DeleteByID", msg.UserID, mock.Anything).Return(nil)

		topicAlarmDigest := new(messenger.MockTopicAlarmDigest)
		topicAlarmDigest.On("Publish", mock.Anything)

		s := service{
			repository: repository,
			publishers: struct{ alarmDigest messenger.TopicAlarmDigest }{
				alarmDigest: topicAlarmDigest,
			},
		}
		s.topicSendAlarmDigestHandler(msg)

		repository.AssertNumberOfCalls(t, "GetByUserID", 1)
		repository.AssertNotCalled(t, "DeleteByID")
		topicAlarmDigest.AssertNotCalled(t, "Publish")
	})
}
