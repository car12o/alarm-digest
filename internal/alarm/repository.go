package alarm

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type repository struct {
	rc *redis.Client
}

func NewRepository(rc *redis.Client) Repository {
	return &repository{rc}
}

func (r *repository) Store(alarm Alarm) error {
	b, err := json.Marshal(alarm)
	if err != nil {
		return err
	}

	if err := r.rc.HSet(context.Background(), alarm.UserID, alarm.ID, b).Err(); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(userID string, ID string) (*Alarm, error) {
	str, err := r.rc.HGet(context.Background(), userID, ID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var alarm Alarm
	if err := json.Unmarshal([]byte(str), &alarm); err != nil {
		return nil, err
	}

	return &alarm, nil
}

func (r *repository) GetByUserID(userID string) ([]*Alarm, error) {
	m, err := r.rc.HGetAll(context.Background(), userID).Result()
	if err != nil {
		return nil, err
	}

	alarms := make([]*Alarm, len(m))
	i := 0
	for _, str := range m {
		var alarm Alarm
		if err := json.Unmarshal([]byte(str), &alarm); err != nil {
			return nil, err
		}
		alarms[i] = &alarm
		i += 1
	}

	return alarms, nil
}

func (r *repository) DeleteByID(userID string, ID ...string) error {
	return r.rc.HDel(context.Background(), userID, ID...).Err()
}
