package router

import (
	"project/handlers"
	"project/middleware"
	"project/storage"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewRouter(store storage.Storage, logger *logrus.Logger) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware(logger))

	h := handlers.NewHandler(store)

	r.HandleFunc("/items", h.GetItems).Methods("GET")
	r.HandleFunc("/items", h.CreateItem).Methods("POST")
	r.HandleFunc("/items/{id}", h.UpdateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", h.DeleteItem).Methods("DELETE")

	return r
}
