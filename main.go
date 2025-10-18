package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/database"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"backend/validation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repo := repositories.NewRepository(db)

	// Initialize services
	service := services.NewService(repo)

	// Initialize Echo server
	e := echo.New()
	e.Validator = validation.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(e, service)

	// Start server in a goroutine
	go func() {
		if err := e.Start(cfg.ServerAddress()); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Printf("Server started on %s", cfg.ServerAddress())
	log.Printf("Environment: %s", cfg.Env)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
