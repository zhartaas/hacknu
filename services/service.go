package services

import (
	"backend/models"
	"backend/repositories"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type Service interface {
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error)
	LLMRequestAndSave(ctx context.Context, message *models.Message, fullChat *models.Chat) (*models.Message, error)
}

func NewService(repo repositories.Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo repositories.Repository
}

func (s *service) LLMRequestAndSave(ctx context.Context, requestMessage *models.Message, fullChat *models.Chat) (*models.Message, error) {
	if err := s.repo.SaveMessage(ctx, requestMessage); err != nil {
		return nil, err
	}

	llmURL := "https://openai-hub.neuraldeep.tech/v1/chat/completions"
	standartModel := "gpt-4o-mini"
	api_key := "sk-roG3OusRr0TLCHAADks6lw"

	req := models.LLMAPIRequest{
		Model: standartModel,
	}
	messages := []models.MessagesAPI{}
	for _, m := range fullChat.Messages {
		mapi := models.MessagesAPI{m.Role, m.Content}
		messages = append(messages, mapi)
	}

	messages = append(messages, models.MessagesAPI{Role: requestMessage.Role, Content: requestMessage.Content})
	req.Messages = messages

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(jsonData))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", llmURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return nil, err
	}

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("LLM API returned error status %d: %s\n", resp.StatusCode, string(body))
		return nil, err
	}
	// Parse the response
	var llmResponse struct {
		Choices []struct {
			Message models.Message `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &llmResponse); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return nil, err
	}

	if len(llmResponse.Choices) == 0 {
		return nil, errors.New("error llm response")
	}

	responseMessage := &llmResponse.Choices[0].Message
	responseMessage.ChatID = fullChat.ID
	if err := s.repo.SaveMessage(ctx, responseMessage); err != nil {
		return nil, err
	}
	fmt.Println("resp", responseMessage)
	return responseMessage, nil
}

func (s *service) GetChatByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error) {
	ch, err := s.repo.GetChatAndMessages(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
