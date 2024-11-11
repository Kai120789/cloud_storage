package tokens

import (
	"cloud/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint, expiresAt time.Time) (string, error) {
	// check SecretKey is string
	if config.AppConfig.SecretKey == "" {
		zap.S().Error("Secret key is empty")
		return "", errors.New("secret key is empty")
	}

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(), // token expire time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.AppConfig.SecretKey)) // key to type []byte
	if err != nil {
		zap.S().Errorf("Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
}
