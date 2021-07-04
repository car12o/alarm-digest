package alarms

import "time"

type Service interface {
	SubscribeTopics() error
}

type Alarm struct {
	ID        string
	UserID    string
	Status    string
	ChangedAt time.Time
}
