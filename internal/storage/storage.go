package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	Conn   *pgxpool.Pool
	Logger *zap.Logger
}

func New(dbConn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{
		Conn:   dbConn,
		Logger: log,
	}
}

func Connection(connectionStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}
