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
	CreateNewChat(ctx context.Context, chatID uuid.UUID, req *models.Message) (*models.Chat, error)
	LLMRequest(ctx context.Context, requestMessage *models.Message, fullChat *models.Chat) (*models.Message, error)
}

func NewService(repo repositories.Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo repositories.Repository
}

func (s *service) LLMRequest(ctx context.Context, requestMessage *models.Message, fullChat *models.Chat) (*models.Message, error) {
	llmURL := "https://openai-hub.neuraldeep.tech/v1/chat/completions"
	api_key := "sk-roG3OusRr0TLCHAADks6lw"

	req := models.LLMAPIRequest{
		Model: models.LLMModel,
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
	return responseMessage, nil
}

func (s *service) LLMRequestAndSave(ctx context.Context, requestMessage *models.Message, fullChat *models.Chat) (*models.Message, error) {
	if err := s.repo.SaveMessage(ctx, requestMessage); err != nil {
		return nil, err
	}
	responseMessage, err := s.LLMRequest(ctx, requestMessage, fullChat)
	if err != nil {
		return nil, err
	}
	responseMessage.ChatID = fullChat.ID
	if err := s.repo.SaveMessage(ctx, responseMessage); err != nil {
		return nil, err
	}
	fmt.Println("resp", responseMessage)
	return responseMessage, nil
}

func (s *service) CreateNewChat(ctx context.Context, chatID uuid.UUID, req *models.Message) (*models.Chat, error) {
	systemMessage := &models.Message{ChatID: chatID, Role: "system", Content: models.BasePrompt}
	chat := &models.Chat{ID: chatID, Messages: []models.Message{*systemMessage}}
	response, err := s.LLMRequest(ctx, req, chat)
	if err != nil {
		return nil, err
	}

	allMessages := chat.Messages
	allMessages = append(allMessages, *req)
	allMessages = append(allMessages, *response)

	chat.Messages = allMessages

	createNamePrompt := "Generate a short and concise title for a chat based on the user prompt. The title should be no more than 5 words"
	chatNameMessage, err := s.LLMRequest(ctx, &models.Message{ChatID: chatID, Role: "user", Content: createNamePrompt}, chat)
	if err != nil {
		return nil, errors.New("llmreq for chat name failed: " + err.Error())
	}

	err = s.repo.CreateNewChat(ctx, chatID, chatNameMessage.Content)
	if err != nil {
		return nil, errors.New("create new chat in repo failed: " + err.Error())
	}

	if err := s.repo.SaveMessage(ctx, systemMessage); err != nil {
		return nil, errors.New("save system message failed: " + err.Error())
	}
	if err := s.repo.SaveMessage(ctx, req); err != nil {
		return nil, errors.New("save req message failed: " + err.Error())
	}
	response.ChatID = chatID
	chat.Title = chatNameMessage.Content
	if err := s.repo.SaveMessage(ctx, response); err != nil {
		return nil, errors.New("save response message failed: " + err.Error())
	}
	return chat, nil
}

func (s *service) GetChatByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error) {
	ch, err := s.repo.GetChatAndMessages(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
