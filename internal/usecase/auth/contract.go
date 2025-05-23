package auth

import (
	"context"
	"time"

	usersentity "domashka-backend/internal/entity/users"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type usersRepo interface {
	CreateWithPhone(ctx context.Context, phone string) (*usersentity.User, error)
	GetByPhone(ctx context.Context, phone string) (*usersentity.User, error)
	Update(ctx context.Context, id int64, user usersentity.User) error
	CheckIfUserIsChef(ctx context.Context, userID int64) (*int64, bool, error)
}

type redisClient interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
	IsExpired(key string) (bool, error)
}

type jwtUsecase interface {
	ValidateJWT(tokenString string) (map[string]interface{}, error)
	GenerateJWT(userID int64, chefID *int64, role string) (string, error)
}

type SMSClient interface {
	Send(phone, message string) error
}
