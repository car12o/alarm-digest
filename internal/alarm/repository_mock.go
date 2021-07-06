package alarm

import "github.com/stretchr/testify/mock"

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Store(alarm Alarm) error {
	args := m.Called(alarm)
	return args.Error(0)
}

func (m *MockRepository) GetByID(userID string, ID string) (*Alarm, error) {
	args := m.Called(userID, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Alarm), args.Error(1)
}

func (m *MockRepository) GetByUserID(userID string) ([]*Alarm, error) {
	args := m.Called(userID)
	return args.Get(0).([]*Alarm), args.Error(1)
}

func (m *MockRepository) DeleteByID(userID string, ID ...string) error {
	args := m.Called(userID, ID)
	return args.Error(0)
}
