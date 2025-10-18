package handlers

import (
	"net/http"
	"strconv"

	"backend/models"
	"backend/services"

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

func (h *Handler) LlmRequest(c echo.Context) error {
	var req models.LlmRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	response := h.service.LlmRequest(c.Request().Context(), &req)
	if response == nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to process LLM request",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    response,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with email and name
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	user, err := h.service.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	user, err := h.service.GetUser(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "User data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	user, err := h.service.UpdateUser(c.Request().Context(), id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	err = h.service.DeleteUser(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// ListUsers godoc
// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.service.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}
