package middleware

import (
	"cloud/internal/config"
	"cloud/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// Claims для JWT
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// middleware for Access token check
func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logger init
		zapLog, err := logger.New(config.AppConfig.LogLevel)
		if err != nil {
			zap.S().Fatalf("init logger error", zap.Error(err))
		}

		log := zapLog.ZapLogger

		// extract token from header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token is missing", http.StatusUnauthorized)
			return
		}

		// check is token start with 'Bearer '
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		// trim prefix 'Bearer '
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// get token
		log.Info("Token:", zap.String("token", tokenString))

		// where save info from token
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// check method is true
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("method is not correct: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.SecretKey), nil
		})

		// check token is valid
		if err != nil {
			log.Error("token parsing error: %v", zap.Error(err))
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Error("invalid token")
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		zap.S().Infof("claims: %+v", claims)

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
