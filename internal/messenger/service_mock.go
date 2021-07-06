package messenger

import "github.com/stretchr/testify/mock"

type MockService struct {
	mock.Mock
}

func (m *MockService) TopicAlarmStatusChanged() TopicAlarmStatusChanged {
	args := m.Called()
	return args.Get(0).(TopicAlarmStatusChanged)
}

func (m *MockService) TopicSendAlarmDigest() TopicSendAlarmDigest {
	args := m.Called()
	return args.Get(0).(TopicSendAlarmDigest)
}

func (m *MockService) TopicAlarmDigest() (TopicAlarmDigest, error) {
	args := m.Called()
	return args.Get(0).(TopicAlarmDigest), args.Error(1)
}

type MockTopicAlarmStatusChanged struct {
	mock.Mock
}

func (m *MockTopicAlarmStatusChanged) WithQueue(queue string) TopicAlarmStatusChanged {
	args := m.Called(queue)
	return args.Get(0).(TopicAlarmStatusChanged)
}

func (m *MockTopicAlarmStatusChanged) Subscribe(handler func(msg AlarmStatusChanged)) error {
	args := m.Called(handler)
	return args.Error(0)
}

type MockTopicSendAlarmDigest struct {
	mock.Mock
}

func (m *MockTopicSendAlarmDigest) WithQueue(queue string) TopicSendAlarmDigest {
	args := m.Called(queue)
	return args.Get(0).(TopicSendAlarmDigest)
}

func (m *MockTopicSendAlarmDigest) Subscribe(handler func(msg SendAlarmDigest)) error {
	args := m.Called(handler)
	return args.Error(0)
}

type MockTopicAlarmDigest struct {
	mock.Mock
}

func (m *MockTopicAlarmDigest) Publish(msg AlarmDigest) {
	m.Called(msg)
}
