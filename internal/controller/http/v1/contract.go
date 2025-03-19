package v1

import (
	"context"
	geoEntity "domashka-backend/internal/entity/geo"
	notifEntity "domashka-backend/internal/entity/notifications"
	entities "domashka-backend/internal/entity/carts"
	usersEntity "domashka-backend/internal/entity/users"
)

type logger interface {
	Error(ctx context.Context, args ...interface{})
}

type authUsecase interface {
	Register(ctx context.Context, phone string) error
	Verify(ctx context.Context, phone, otp string) (string, error)
	Login(ctx context.Context, phone string) error
}

type jwtUsecase interface {
	ValidateJWT(token string) (string, error)
}

type usersUsecase interface {
	Create(ctx context.Context, user *usersEntity.User) (*string, error)
	GetByID(ctx context.Context, id string) (*usersEntity.User, error)
	Update(ctx context.Context, id string, user usersEntity.User) error
	Delete(ctx context.Context, id string) error
}

type geoUsecase interface {
	AddClientAddress(ctx context.Context, clientID int, address geoEntity.Address) error
	AddChefAddress(ctx context.Context, chefID int, address geoEntity.Address) error
	GetClientAddresses(ctx context.Context, clientID int) ([]geoEntity.Address, error)
	GetChefAddress(ctx context.Context, chefID int) (geoEntity.Address, error)
	UpdateClientAddress(ctx context.Context, clientID int, addressID int, address geoEntity.Address) error
	UpdateChefAddress(ctx context.Context, chefID int, address geoEntity.Address) error
	FindChefsNearAddress(ctx context.Context, clientAddressID int, radius float64) ([]geoEntity.Address, error)
	FindClientsNearAddress(ctx context.Context, chefID int, radius float64) ([]geoEntity.Address, error)
}

type notificationUsecase interface {
	CreateNotification(ctx context.Context, n notifEntity.Notification) (int, error)
	SendEmailNotification(ctx context.Context, n notifEntity.Notification) error
	ResendNotification(ctx context.Context, id int) error
	GetNotificationByID(ctx context.Context, id int) (*notifEntity.Notification, error)
	GetNotifications(ctx context.Context, filters map[string]string, page, limit int) ([]notifEntity.Notification, int, error)
}

type cartUsecase interface {
	CreateCart(ctx context.Context, cart *entities.Cart) error
	GetCart(ctx context.Context, userId string) (*entities.Cart, error)
	AddCartItem(ctx context.Context, cartItem *entities.CartItem) (itemID *string, err error)
	UpdateCartItem(ctx context.Context, cartItem *entities.CartItem) (itemID *string, err error)
	DeleteCartItem(ctx context.Context, itemID string) error
	GetCartItems(ctx context.Context, userID string) ([]entities.CartItem, error)
	ClearCart(ctx context.Context, userID string) error
}
