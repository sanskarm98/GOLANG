package storage

import (
	"project/models"
)

type Storage interface {
	GetItems() ([]models.Item, error)
	CreateItem(models.Item) error
	UpdateItem(string, models.Item) error
	DeleteItem(string) error
}

type InMemoryStorage struct {
	items map[string]models.Item
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		items: make(map[string]models.Item),
	}
}

func (s *InMemoryStorage) GetItems() ([]models.Item, error) {
	var items []models.Item
	for _, item := range s.items {
		items = append(items, item)
	}
	return items, nil
}

func (s *InMemoryStorage) CreateItem(item models.Item) error {
	if _, ok := s.items[item.ID]; ok {
		return models.ErrItemExists
	}
	s.items[item.ID] = item
	return nil
}

func (s *InMemoryStorage) UpdateItem(id string, item models.Item) error {
	if _, ok := s.items[id]; !ok {
		return models.ErrItemNotFound
	}
	s.items[id] = item
	return nil
}

func (s *InMemoryStorage) DeleteItem(id string) error {
	if _, ok := s.items[id]; !ok {
		return models.ErrItemNotFound
	}
	delete(s.items, id)
	return nil
}
