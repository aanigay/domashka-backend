package order

import (
	"context"
	"domashka-backend/internal/entity/orders"
	"fmt"
)

type Usecase struct {
	geoUsecase  geoUsecase
	cartUsecase cartUsecase
	shiftsRepo  shiftsRepo
	ordersRepo  ordersRepo
}

func New(
	geoUsecase geoUsecase,
	cartUsecase cartUsecase,
	shiftsRepo shiftsRepo,
	ordersRepo ordersRepo,
) *Usecase {
	return &Usecase{
		geoUsecase:  geoUsecase,
		cartUsecase: cartUsecase,
		shiftsRepo:  shiftsRepo,
		ordersRepo:  ordersRepo,
	}
}

func (u *Usecase) CreateOrder(ctx context.Context, userID int64, leaveByTheDoor, callBeforehand bool) (int64, error) {
	// Get client address
	address, err := u.geoUsecase.GetLastUpdatedClientAddress(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Get Cart Items
	cartItems, err := u.cartUsecase.GetCartItems(ctx, userID)
	if len(cartItems) == 0 {
		return 0, fmt.Errorf("нет товаров в корзине")
	}
	// Calculate total profit
	totalProfit := float32(0)
	chefID := int64(0)
	for _, item := range cartItems {
		totalProfit += float32(item.Quantity) * item.Size.PriceValue
		chefID = item.Dish.ChefID
	}

	// Get ChefID

	// Get Active Shift
	shiftID, err := u.shiftsRepo.GetActiveShiftIDByChefID(ctx, chefID)
	if err != nil {
		return 0, err
	}

	// Create order
	orderID, err := u.ordersRepo.CreateOrder(
		ctx,
		chefID,
		shiftID,
		address.ID,
		totalProfit,
		leaveByTheDoor,
		callBeforehand,
	)

	// Add cart items to order
	for _, item := range cartItems {
		err = u.ordersRepo.AddCartItemToOrder(ctx, item.ID, orderID)
		if err != nil {
			return 0, err
		}
	}
	return orderID, nil
}

func (u *Usecase) GetOrdersByShiftID(ctx context.Context, shiftID int64) ([]orders.Order, error) {
	return u.ordersRepo.GetOrdersByShiftID(ctx, shiftID)
}
