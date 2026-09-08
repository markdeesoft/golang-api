package handler_test

import (
	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepositoryTest struct {
	mock.Mock
}

func (m *MockUserRepositoryTest) Store(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepositoryTest) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepositoryTest) GetAll(limit, offset int) ([]model.User, error) {

	args := m.Called(limit, offset)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepositoryTest) Update(id uint, user *model.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepositoryTest) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepositoryTest) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepositoryTest) ResetPassword(id uint, hashedPassword string) error {
	args := m.Called(id)
	return args.Error(0)
}

// func (m *MockUserRepositoryTest) ResetPassword() (int64, error) {
// 	args := m.Called()
// 	return args.Get(0).(int64), args.Error(1)
// }
