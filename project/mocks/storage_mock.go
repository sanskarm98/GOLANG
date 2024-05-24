package mocks

import (
	"project/models"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetItems() ([]models.Item, error) {
	args := m.Called()
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *MockStorage) CreateItem(item models.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockStorage) UpdateItem(id string, item models.Item) error {
	args := m.Called(id, item)
	return args.Error(0)
}

func (m *MockStorage) DeleteItem(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
