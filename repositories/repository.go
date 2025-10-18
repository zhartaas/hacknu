package repositories

import (
	"backend/database"
	"backend/models"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetChatAndMessages(ctx context.Context, id uuid.UUID) (*models.Chat, error)
	SaveMessage(ctx context.Context, message *models.Message) error
}

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveMessage(ctx context.Context, message *models.Message) error {
	query := `
INSERT INTO messages(chat_id, role, content) VALUES($1,$2,$3) 
`
	_, err := r.db.Pool.Exec(ctx, query, message.ChatID, message.Role, message.Content)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetChatAndMessages(ctx context.Context, id uuid.UUID) (*models.Chat, error) {
	chat := &models.Chat{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, title, model, created_at
		FROM chats
		WHERE id = $1
	`, id).Scan(&chat.ID, &chat.Title, &chat.Model, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, chat_id, role, content, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC, id ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chat.Messages = make([]models.Message, 0, 32)
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		chat.Messages = append(chat.Messages, m)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return chat, nil

}
