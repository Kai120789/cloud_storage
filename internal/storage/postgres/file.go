package postgres

import (
	"cloud/internal/dto"
	"cloud/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type FileStorage struct {
	Conn   *pgxpool.Pool
	Logger *zap.Logger
}

func NewFileStorage(dbConn *pgxpool.Pool, log *zap.Logger) *FileStorage {
	return &FileStorage{
		Conn:   dbConn,
		Logger: log,
	}
}

// func upload file and create new record in files table
func (s *FileStorage) UploadFile(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

// func create folder and create new record in files table
func (s *FileStorage) CreateFolder(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

// func delete file or folder
func (s *FileStorage) DeleteItem(path string) error {
	return nil
}

// func update record(file or folder) in db and rename it
func (s *FileStorage) RenameItem(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

// func search file by name
func (s *FileStorage) SearchFiles() error {
	return nil
}

// func return items in directory
func (s *FileStorage) ListDirectory() error {
	return nil
}
