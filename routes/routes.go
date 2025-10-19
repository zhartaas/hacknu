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

	v1.GET("/get-chat/:id", handler.GetChatByID)
	v1.POST("/llm-prompt/:id", handler.LLMChat)
	v1.POST("/start", handler.StartNewChat)

}
