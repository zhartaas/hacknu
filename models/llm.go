package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title,omitempty"`
	Model     string    `json:"model"` // lives on chat
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	ID        uuid.UUID `json:"id,omitempty"`
	ChatID    uuid.UUID `json:"chat_id,omitempty"`
	Role      string    `json:"role,omitempty"` // "user" | "assistant"
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type LLMChatRequest struct {
	Content string `json:"content"`
}

type LLMAPIRequest struct {
	Model    string        `json:"model"`
	Messages []MessagesAPI `json:"messages"`
}

type MessagesAPI struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
