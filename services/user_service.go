package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/models"
	"backend/repositories"

	"github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	LlmRequest(ctx context.Context, req *models.LlmRequest) *models.LlmMessage
}

type service struct {
	repo repositories.Repository
}

func NewService(repo repositories.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) LlmRequest(ctx context.Context, req *models.LlmRequest) *models.LlmMessage {
	// Define the LLM API endpoint - you can make this configurable
	llmURL := "https://openai-hub.neuraldeep.tech/v1/chat/completions"
	standartModel := "gpt-4o-mini"
	api_key := "sk-roG3OusRr0TLCHAADks6lw"

	// Prepare the request payload for the LLM API
	req.Model = standartModel

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return nil
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", llmURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key)) // You should make this configurable

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return nil
	}

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("LLM API returned error status %d: %s\n", resp.StatusCode, string(body))
		return nil
	}

	// Parse the response
	var llmResponse struct {
		Choices []struct {
			Message models.LlmMessage `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &llmResponse); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return nil
	}

	// Return the first choice message
	if len(llmResponse.Choices) > 0 {
		return &llmResponse.Choices[0].Message
	}

	return nil
}

func (s *service) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Check if user with email already exists
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Email != nil {
		// Check if new email is already taken by another user
		existingUser, err := s.repo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("user with email %s already exists", *req.Email)
		}
		user.Email = *req.Email
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
