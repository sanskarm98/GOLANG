package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/handlers"
	"project/mocks"
	"project/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	mockItems := []models.Item{{ID: "1", Name: "Test Item 1", Price: 100}}
	mockStore.On("GetItems").Return(mockItems, nil)

	req, err := http.NewRequest("GET", "/items", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items", handler.GetItems).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var items []models.Item
	err = json.NewDecoder(rr.Body).Decode(&items)
	assert.NoError(t, err)
	assert.Equal(t, mockItems, items)
}

func TestGetItemsError(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	mockStore.On("GetItems").Return([]models.Item{}, assert.AnError)

	req, err := http.NewRequest("GET", "/items", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items", handler.GetItems).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCreateItem(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	newItem := models.Item{ID: "1", Name: "Test Item 1", Price: 100}
	mockStore.On("CreateItem", newItem).Return(nil)

	body, err := json.Marshal(newItem)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items", handler.CreateItem).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var returnedItem models.Item
	err = json.NewDecoder(rr.Body).Decode(&returnedItem)
	assert.NoError(t, err)
	assert.Equal(t, newItem, returnedItem)
}

func TestCreateItemConflict(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	newItem := models.Item{ID: "1", Name: "Test Item 1", Price: 100}
	mockStore.On("CreateItem", newItem).Return(models.ErrItemExists)

	body, err := json.Marshal(newItem)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items", handler.CreateItem).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCreateItemBadRequest(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items", handler.CreateItem).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateItem(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	updatedItem := models.Item{ID: "1", Name: "Updated Item", Price: 200}
	mockStore.On("UpdateItem", "1", updatedItem).Return(nil)

	body, err := json.Marshal(updatedItem)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var returnedItem models.Item
	err = json.NewDecoder(rr.Body).Decode(&returnedItem)
	assert.NoError(t, err)
	assert.Equal(t, updatedItem, returnedItem)
}

func TestUpdateItemNotFound(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	updatedItem := models.Item{ID: "1", Name: "Updated Item", Price: 200}
	mockStore.On("UpdateItem", "1", updatedItem).Return(models.ErrItemNotFound)

	body, err := json.Marshal(updatedItem)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateItemBadRequest(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	req, err := http.NewRequest("PUT", "/items/1", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteItem(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	mockStore.On("DeleteItem", "1").Return(nil)

	req, err := http.NewRequest("DELETE", "/items/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteItemNotFound(t *testing.T) {
	mockStore := new(mocks.MockStorage)
	handler := handlers.NewHandler(mockStore)

	mockStore.On("DeleteItem", "1").Return(models.ErrItemNotFound)

	req, err := http.NewRequest("DELETE", "/items/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
