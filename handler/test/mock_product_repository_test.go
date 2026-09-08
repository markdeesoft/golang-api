package handler_test

import (
	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/mock"
)

type MockProductRepositoryTest struct {
	mock.Mock
}

func (m *MockProductRepositoryTest) Store(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepositoryTest) GetByID(id uint) (*model.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepositoryTest) GetAll(limit, offset int) ([]model.Product, error) {

	args := m.Called(limit, offset)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepositoryTest) Update(id uint, product *model.Product) (*model.Product, error) {
	args := m.Called(id, product)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepositoryTest) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepositoryTest) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
