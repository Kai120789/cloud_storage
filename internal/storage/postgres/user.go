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

type UserStorage struct {
	Conn   *pgxpool.Pool
	Logger *zap.Logger
}

func NewUserStorage(dbConn *pgxpool.Pool, log *zap.Logger) *UserStorage {
	return &UserStorage{
		Conn:   dbConn,
		Logger: log,
	}
}

// register new user
func (s *UserStorage) RegisterNewUser(body dto.User) (*models.UserToken, error) {
	var id uint
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := s.Conn.QueryRow(context.Background(), query, body.Login, body.Password).Scan(&id)
	if err != nil {
		return nil, err
	}

	userRet, err := s.GetAuthUser(uint(id))
	if err != nil {
		return nil, err
	}

	return userRet, nil
}

// login user
func (s *UserStorage) AuthorizateUser(body dto.User) (*uint, *string, error) {
	var id uint
	var passwordHash string

	query := `SELECT id, password FROM users WHERE username=$1`
	err := s.Conn.QueryRow(context.Background(), query, body.Login).Scan(&id, &passwordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	return &id, &passwordHash, nil
}

// get auth user
func (s *UserStorage) GetAuthUser(id uint) (*models.UserToken, error) {
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

// logout user
func (s *UserStorage) UserLogout(id uint) error {
	query := `DELETE FROM user_token WHERE user_id=$1`
	_, err := s.Conn.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// add refresh token to db
func (d *UserStorage) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	query := `INSERT INTO user_token (user_id, refresh_token) VALUES ($1, $2)`
	_, err := d.Conn.Exec(context.Background(), query, userId, refreshTokenValue)
	if err != nil {
		return err
	}

	return nil
}
