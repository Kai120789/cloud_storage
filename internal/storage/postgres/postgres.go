package postgres

import (
	"cloud/internal/dto"
	"cloud/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PostgreStorage struct {
	Conn   *pgxpool.Pool
	Logger *zap.Logger
}

func NewPostgresStorage(dbConn *pgxpool.Pool, log *zap.Logger) *PostgreStorage {
	return &PostgreStorage{
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

// register new user
func (s *PostgreStorage) RegisterNewUser(body dto.User) (*models.User, error) {
	var userRet models.User
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, login, password, created_at`
	row := s.Conn.QueryRow(context.Background(), query, body.Login, body.Password)

	err := row.Scan(&userRet.ID,
		&userRet.Login,
		&userRet.Password,
		&userRet.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &userRet, nil
}

// login user
func (s *PostgreStorage) AuthorizateUser(body dto.User) (*uint, *string, error) {
	var id uint
	var passwordHash string

	query := `SELECT id, password FROM users WHERE login=$1`
	row := s.Conn.QueryRow(context.Background(), query, body.Login)

	err := row.Scan(&id, &passwordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	return &id, &passwordHash, nil
}

// get auth user
func (s *PostgreStorage) GetAuthUser(id uint) (*models.UserToken, error) {
	query := `SELECT * FROM user_token WHERE user_id=$1`
	row := s.Conn.QueryRow(context.Background(), query, id)

	var token models.UserToken
	err := row.Scan(&token.ID, &token.UserID, &token.RefreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}
