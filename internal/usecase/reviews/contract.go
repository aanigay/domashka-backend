package reviews

import (
	"context"

	"github.com/segmentio/kafka-go"

	cartentity "domashka-backend/internal/entity/cart"
	reviewEntity "domashka-backend/internal/entity/reviews"
	userEntity "domashka-backend/internal/entity/users"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type reviewRepo interface {
	GetByID(ctx context.Context, id int64) (*reviewEntity.Review, error)
	ListByChef(ctx context.Context, chefID int64, includeOnly bool, limit *int) ([]reviewEntity.Review, error)
	Create(ctx context.Context, rv *reviewEntity.Review) error
	Update(ctx context.Context, rv *reviewEntity.Review) error
	SoftDelete(ctx context.Context, id int64) error
	GetReviewByOrderAndUserID(ctx context.Context, chefID, userID int64) (*reviewEntity.Review, error)
}

type userRepo interface {
	GetByID(ctx context.Context, id int64) (*userEntity.User, error)
}

type KafkaWriter interface {
	WriteMessages(context.Context, ...kafka.Message) error
}

type orderRepo interface {
	GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error)
}
