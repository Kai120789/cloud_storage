package storage

import (
	"cloud/internal/dto"
	"cloud/internal/models"
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

// func upload file and create new record in files table
func (s *Storage) UploadFile(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

// func create folder and create new record in files table
func (s *Storage) CreateFolder(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

// func delete file or folder
func (s *Storage) DeleteItem(path string) error {
	return nil
}

// func update record(file or folder) in db and rename it
func (s *Storage) RenameItem(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

// func search file by name
func (s *Storage) SearchFiles() error {
	return nil
}

// func return items in directory
func (s *Storage) ListDirectory() error {
	return nil
}
