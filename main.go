package main

import (
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/router"
	"log"
	"os"
	"os/signal"
)

func main() {
	config.Init()
	logger.Init()

	db := config.DatabaseConnection()

	app := router.InitRouter(db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Run the app in a goroutine
	go func() {
		if err := app.Listen(":8888"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Gracefully shutting down server...")

	// Gracefully shutdown Fiber
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server successfully stopped")
}
