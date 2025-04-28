package v1

import (
	"context"
	cartentity "domashka-backend/internal/entity/cart"
	"domashka-backend/internal/entity/orders"
	"domashka-backend/internal/entity/shifts"
	"mime/multipart"

	authEntity "domashka-backend/internal/entity/auth"
	chefEntity "domashka-backend/internal/entity/chefs"
	dishesEntity "domashka-backend/internal/entity/dishes"
	geoEntity "domashka-backend/internal/entity/geo"
	notifEntity "domashka-backend/internal/entity/notifications"
	usersEntity "domashka-backend/internal/entity/users"
)

type logger interface {
	Error(ctx context.Context, args ...interface{})
}

type authUsecase interface {
	Auth(ctx context.Context, req authEntity.Request) error
	Verify(ctx context.Context, phone, otp, role string) (int64, string, error)
	AuthViaTg(ctx context.Context, phoneNumber string) error
	AuthViaTgStatus(ctx context.Context, phoneNumber string) (string, error)
}

type jwtUsecase interface {
	ValidateJWT(token string) (map[string]interface{}, error)
}

type usersUsecase interface {
	Create(ctx context.Context, user *usersEntity.User) error
	GetByID(ctx context.Context, id int64) (*usersEntity.User, error)
	Update(ctx context.Context, id int64, user usersEntity.User) error
	Delete(ctx context.Context, id int64) error
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
	GetLastUpdatedClientAddress(ctx context.Context, clientID int64) (*geoEntity.Address, error)
	GetAddressByID(ctx context.Context, id int64) (*geoEntity.Address, error)
	PushClientAddress(ctx context.Context, addressID int64) error
}

type notificationUsecase interface {
	CreateNotification(ctx context.Context, n notifEntity.Notification) (int, error)
	SendEmailNotification(ctx context.Context, n notifEntity.Notification) error
	ResendNotification(ctx context.Context, id int) error
	GetNotificationByID(ctx context.Context, id int) (*notifEntity.Notification, error)
	GetNotifications(ctx context.Context, filters map[string]string, page, limit int) ([]notifEntity.Notification, int, error)
}

type dishesUsecase interface {
	GetDishByID(ctx context.Context, dishID int64) (*dishesEntity.Dish, error)
	GetDishesByChefID(ctx context.Context, chefID int64) ([]dishesEntity.Dish, error)
	GetNutritionByDishID(ctx context.Context, dishID int64) (*dishesEntity.Nutrition, error)
	GetDishSizesByDishID(ctx context.Context, dishID int64) ([]dishesEntity.Size, error)
	GetIngredientsByDishID(ctx context.Context, dishID int64) ([]dishesEntity.Ingredient, error)
	GetMinimalPriceByDishID(ctx context.Context, dishID int64) (*dishesEntity.Price, error)
	GetTopDishes(ctx context.Context, limit int) ([]dishesEntity.Dish, error)
}

type chefUsecase interface {
	GetChefByDishID(ctx context.Context, dishID int64) (*chefEntity.Chef, error)
	GetChefByID(ctx context.Context, chefID int64) (*chefEntity.Chef, error)
	UploadAvatar(ctx context.Context, chefID int64, fileHeader *multipart.FileHeader) (string, error)
	GetTopChefs(ctx context.Context, limit int) ([]chefEntity.Chef, error)
	GetChefAvatarURLByChefID(ctx context.Context, dishID int64) (string, error)
	GetChefExperienceYears(ctx context.Context, chefID int64) (int, error)
	GetChefCertifications(ctx context.Context, chefID int64) ([]chefEntity.Certification, error)
}

type cartUsecase interface {
	AddItem(
		ctx context.Context,
		userID int64,
		dish dishesEntity.Dish,
		sizeID int64,
		addedIngredients []int64,
		removedIngredients []int64,
		notes string,
	) (int64, error)
	RemoveItem(ctx context.Context, cartItemID int64) error
	GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error)

	ClearCart(ctx context.Context, userID int64) error
	IncrementCartItem(ctx context.Context, cartItemID int64) (newQuantity int32, err error)
	DecrementCartItem(ctx context.Context, cartItemID int64) (newQuantity int32, err error)
	GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error)
}

type orderUsecase interface {
	CreateOrder(ctx context.Context, userID int64, leaveByTheDoor, callBeforehand bool) (int64, error)
	GetOrdersByShiftID(ctx context.Context, shiftID int64) ([]orders.Order, error)
}

type shiftsUsecase interface {
	GetActiveShiftByChefID(ctx context.Context, chefID int64) (*shifts.Shift, error)
}
