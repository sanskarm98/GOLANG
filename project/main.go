package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/router"
	"project/storage"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := logrus.New()
	logger.Out = logFile

	store := storage.NewInMemoryStorage()
	r := router.NewRouter(store, logger)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel to listen for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine to allow signal handling
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	logger.Println("Server started on :8080")

	// Block until we receive a signal
	<-quit
	logger.Println("Shutting down server...")

	// Create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}
