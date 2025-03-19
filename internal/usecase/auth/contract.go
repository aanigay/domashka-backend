package auth

import (
	"context"
	"time"

	usersentity "domashka-backend/internal/entity/users"
)

type usersRepo interface {
	CreateWithPhone(ctx context.Context, phone string) error
	GetByPhone(ctx context.Context, phone string) (*usersentity.User, error)
}

type redisClient interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
}

type jwtUsecase interface {
	ValidateJWT(tokenString string) (string, error)
	GenerateJWT() (string, error)
}
