package handler_test

import (
	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/mock"
)

type MockProductCategoryRepositoryTest struct {
	mock.Mock
}

func (m *MockProductCategoryRepositoryTest) Store(product *model.ProductCategory) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductCategoryRepositoryTest) GetByID(id uint) (*model.ProductCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepositoryTest) GetAll(limit, offset int) ([]model.ProductCategory, error) {

	args := m.Called(limit, offset)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepositoryTest) Update(id uint, product *model.ProductCategory) (*model.ProductCategory, error) {
	args := m.Called(id, product)
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepositoryTest) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductCategoryRepositoryTest) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
