package postgres

import (
	"cloud/internal/dto"
	"cloud/internal/models"
	"context"
	"time"

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
func (s *FileStorage) CreateNewFileOrFold(dto dto.Object) (*models.Object, error) {
	var id uint
	var createdAt time.Time
	query := `INSERT INTO files (name, path, user_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	row := s.Conn.QueryRow(context.Background(), query, dto.Name, dto.Path, dto.UserID)

	err := row.Scan(&id, &createdAt)
	if err != nil {
		s.Logger.Error("create file or folder error", zap.Error(err))
		return nil, err
	}

	var obj models.Object = models.Object{
		ID:        id,
		Name:      dto.Name,
		Path:      dto.Path,
		UserID:    dto.UserID,
		CreatedAt: createdAt,
	}

	return &obj, nil
}

// func delete file or folder
func (s *FileStorage) DeleteItem(path string) error {
	query := `DELETE FROM files WHERE path=$1`
	_, err := s.Conn.Exec(context.Background(), query, path)
	if err != nil {
		s.Logger.Error("delete file or folder error", zap.Error(err))
		return err
	}
	return nil
}

// func update record(file or folder) in db and rename it
func (s *FileStorage) RenameItem(dto dto.Object) (*models.Object, error) {
	query := `UPDATE files SET name=$1 WHERE path=$2`
	_, err := s.Conn.Exec(context.Background(), query, dto.Name, dto.Path)
	if err != nil {
		s.Logger.Error("rename file or folder error", zap.Error(err))
		return nil, err
	}

	query = `SELECT * FROM files WHERE path=$1`
	row := s.Conn.QueryRow(context.Background(), query, dto.Path)

	var obj models.Object
	err = row.Scan(&obj.ID, &obj.Name, &obj.Path, &obj.UserID, &obj.CreatedAt)
	if err != nil {
		s.Logger.Error("set file or folder error", zap.Error(err))
		return nil, err
	}
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
