package repository

import (
	"context"
	"database/sql"
	appmodels "server/internal/app_models"
)

type Repository struct {
	RoomManager
	ClientManager
	MessageManager
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func NewRepository(db DBTX) *Repository {
	return &Repository{
		RoomManager:    NewRoomRepository(db),
		ClientManager:  NewClientRepository(db),
		MessageManager: NewMessageRepository(db),
	}
}

type RoomManager interface {
	CreateRoom(ctx context.Context, req appmodels.CreateRoomReq) (int64, error)
	GetRoom(ctx context.Context, req appmodels.AddClientReq) (*appmodels.GetRoomResp, error)
	GetAllRooms(ctx context.Context) ([]appmodels.BackupRoom, error)
}

type ClientManager interface {
	AddClient(ctx context.Context, req appmodels.AddClientReq) error
}

type MessageManager interface {
	AddMessage(ctx context.Context, req appmodels.AddMessageReq) error
}
