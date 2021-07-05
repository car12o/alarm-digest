package messenger

import "time"

type Service interface {
	TopicAlarmStatusChanged() TopicAlarmStatusChanged
	TopicSendAlarmDigest() TopicSendAlarmDigest
	TopicAlarmDigest() (TopicAlarmDigest, error)
}

type TopicAlarmStatusChanged interface {
	WithQueue(queue string) TopicAlarmStatusChanged
	Subscribe(handler func(msg AlarmStatusChanged)) error
}

type TopicSendAlarmDigest interface {
	WithQueue(queue string) TopicSendAlarmDigest
	Subscribe(handler func(msg SendAlarmDigest)) error
}

type TopicAlarmDigest interface {
	Publish(msg AlarmDigest)
}

type AlarmStatusChanged struct {
	AlarmID   string
	UserID    string
	Status    string
	ChangedAt time.Time
}

type SendAlarmDigest struct {
	UserID string
}

type AlarmDigest struct {
	UserID       string
	ActiveAlarms []*ActiveAlarm
}

type ActiveAlarm struct {
	AlarmID         string
	Status          string
	LatestChangedAt time.Time
}
