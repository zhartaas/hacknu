package handlers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *Handler) StartNewChat(c echo.Context) error {
	chatID := uuid.New()

	var req *models.LLMChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid JSON body",
		})
	}

	userMessage := &models.Message{
		ChatID:  chatID,
		Role:    "user",
		Content: req.Content,
	}
	chat, err := h.service.CreateNewChat(c.Request().Context(), chatID, userMessage)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false, "error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    chat,
	})
}

func (h *Handler) LLMChat(c echo.Context) error {
	chatIDStr := c.Param("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false, "error": "invalid chat id",
		})
	}

	var req *models.LLMChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid JSON body",
		})
	}

	fullChat, err := h.service.GetChatByID(c.Request().Context(), chatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false, "error": err.Error(),
		})
	}

	userMessage := &models.Message{
		ChatID:  chatID,
		Role:    "user",
		Content: req.Content,
	}

	responseMessage, err := h.service.LLMRequestAndSave(c.Request().Context(), userMessage, fullChat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false, "error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    responseMessage,
	})
}

func (h *Handler) GetChatByID(c echo.Context) error {
	chatIDStr := c.Param("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false, "error": "invalid chat id",
		})
	}

	chat, err := h.service.GetChatByID(c.Request().Context(), chatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false, "error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    chat,
	})
}
