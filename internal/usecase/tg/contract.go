package tg

import (
	"context"
	usersentity "domashka-backend/internal/entity/users"
	"time"
)

type redisClient interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	IsExpired(key string) (bool, error)
}

type jwtUsecase interface {
	GenerateJWT(userID int64, role string) (string, error)
}

type usersRepo interface {
	Create(ctx context.Context, user *usersentity.User) error
	GetByPhone(ctx context.Context, phone string) (*usersentity.User, error)
	Update(ctx context.Context, id int64, user usersentity.User) error
}
