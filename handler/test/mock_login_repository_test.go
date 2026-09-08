package handler_test

import (
	"github.com/markdeesoft/golang-api/model"
)

func (m *MockUserRepositoryTest) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
