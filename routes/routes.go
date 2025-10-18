package routes

import (
	"backend/handlers"
	"backend/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo, service services.Service) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Initialize handlers
	handler := handlers.NewHandler(service)

	v1.POST("/request", handler.LlmRequest)
	// User routes
	users := v1.Group("/users")
	users.POST("", handler.CreateUser)
	users.GET("", handler.ListUsers)
	users.GET("/:id", handler.GetUser)
	users.PUT("/:id", handler.UpdateUser)
	users.DELETE("/:id", handler.DeleteUser)

}
