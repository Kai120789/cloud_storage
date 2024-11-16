package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	FileStorage FileStorage
	UserStorage UserStorage
}

func New(Conn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{
		FileStorage: *NewFileStorage(Conn, log),
		UserStorage: *NewUserStorage(Conn, log),
	}
}

func Connection(connectionStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}
