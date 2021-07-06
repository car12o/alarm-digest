package logger

import "github.com/stretchr/testify/mock"

type MockService struct {
	mock.Mock
}

func (m *MockService) Debug(msg string, fields ...map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockService) Info(msg string, fields ...map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockService) Warn(msg string, fields ...map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockService) Error(err error, fields ...map[string]interface{}) {
	m.Called(err, fields)
}

func (m *MockService) Fatal(msg string, fields ...map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockService) Panic(msg string, fields ...map[string]interface{}) {
	m.Called(msg, fields)
}
