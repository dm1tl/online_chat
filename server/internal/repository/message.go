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

func (r *MessageRepository) GetAllMessages(ctx context.Context) (appmodels.BackupMessages, error) {
	op := "repository.GetAllMessage"

	query := "SELECT m.client_id, m.room_id, m.content, c.username FROM messages AS m JOIN clients AS c ON c.id = m.client_id ORDER BY m.created_at ASC"

	output := make(appmodels.BackupMessages)
	res, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer res.Close()
	if !res.Next() {
		return output, nil
	}
	for {
		var msg appmodels.BackupMessage
		if err := res.Scan(&msg.UserID, &msg.RoomID, &msg.Content, &msg.Username); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		output[msg.RoomID] = append(output[msg.RoomID], msg)
		if !res.Next() {
			break
		}
	}
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return output, nil

}
