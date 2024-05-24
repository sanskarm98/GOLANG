package models

import "errors"

var (
	ErrItemExists   = errors.New("item already exists")
	ErrItemNotFound = errors.New("item not found")
)

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
