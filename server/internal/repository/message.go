package repository

import (
	"context"
	"fmt"
	appmodels "server/internal/app_models"

	"github.com/sirupsen/logrus"
)

type MessageRepository struct {
	db DBTX
}

func NewMessageRepository(db DBTX) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) AddMessage(ctx context.Context, req appmodels.AddMessageReq) error {
	op := "repository.AddMessage"

	clQuery := "INSERT INTO messages (client_id, room_id, content) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, clQuery, req.UserID, req.RoomID, req.Content)
	if err != nil {
		logrus.Error(op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
