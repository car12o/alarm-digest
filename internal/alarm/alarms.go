package alarm

import "time"

type Service interface {
	TopicsSubscribe() error
	TopicsInitPublishers() error
}

type Alarm struct {
	ID        string
	UserID    string
	Status    string
	ChangedAt time.Time
}

type Repository interface {
	Store(alarm Alarm) error
	GetByID(userID string, ID string) (*Alarm, error)
	GetByUserID(userID string) ([]*Alarm, error)
	DeleteByID(userID string, ID ...string) error
}
