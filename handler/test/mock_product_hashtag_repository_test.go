package handler_test

import (
	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/mock"
)

type MockProductHashtagRepositoryTest struct {
	mock.Mock
}

func (m *MockProductHashtagRepositoryTest) Store(product *model.ProductHashtag) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductHashtagRepositoryTest) GetByID(id uint) (*model.ProductHashtag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductHashtag), args.Error(1)
}

func (m *MockProductHashtagRepositoryTest) GetAll(limit, offset int) ([]model.ProductHashtag, error) {

	args := m.Called(limit, offset)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.ProductHashtag), args.Error(1)
}

func (m *MockProductHashtagRepositoryTest) Update(id uint, product *model.ProductHashtag) (*model.ProductHashtag, error) {
	args := m.Called(id, product)
	return args.Get(0).(*model.ProductHashtag), args.Error(1)
}

func (m *MockProductHashtagRepositoryTest) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
