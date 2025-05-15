package order

import (
	"context"
	cartentity "domashka-backend/internal/entity/cart"
	chefEntity "domashka-backend/internal/entity/chefs"
	dishEntity "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/entity/orders"
	"fmt"
)

// ReviewDetail описывает отзыв пользователя на заказ
type ReviewDetail struct {
	CanWriteReview bool    // можно ли писать отзыв
	Rating         *int    // оценка пользователя
	Comment        *string // комментарий пользователя
}

type Usecase struct {
	geoUsecase    geoUsecase
	cartUsecase   cartUsecase
	dishesUsecase dishesUsecase
	chefsUsecase  chefsUsecase
	reviewUsecase reviewUsecase
	shiftsRepo    shiftsRepo
	ordersRepo    ordersRepo
}

func New(
	geoUsecase geoUsecase,
	cartUsecase cartUsecase,
	shiftsRepo shiftsRepo,
	ordersRepo ordersRepo,
	dishesUsecase dishesUsecase,
	chefsUsecase chefsUsecase,
	reviewUsecase reviewUsecase,
) *Usecase {
	return &Usecase{
		geoUsecase:    geoUsecase,
		cartUsecase:   cartUsecase,
		shiftsRepo:    shiftsRepo,
		ordersRepo:    ordersRepo,
		dishesUsecase: dishesUsecase,
		chefsUsecase:  chefsUsecase,
		reviewUsecase: reviewUsecase,
	}
}

func (u *Usecase) CreateOrder(ctx context.Context, userID int64, leaveByTheDoor, callBeforehand bool) (int64, error) {
	// Get client address
	address, err := u.geoUsecase.GetLastUpdatedClientAddress(ctx, userID)
	if err != nil {
		return 0, err
	}
	if address == nil {
		return 0, fmt.Errorf("address is nil")
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

	if shiftID == 0 {
		return 0, fmt.Errorf("повар сейчас не работает")
	}

	// Create order
	orderID, err := u.ordersRepo.CreateOrder(
		ctx,
		chefID,
		shiftID,
		userID,
		address.ID,
		totalProfit,
		leaveByTheDoor,
		callBeforehand,
	)

	// Add cart items to order
	for _, item := range cartItems {
		err = u.ordersRepo.AddCartItemToOrder(ctx, &item, orderID, userID)
		if err != nil {
			return 0, err
		}
	}
	return orderID, nil
}

func (u *Usecase) GetOrdersByShiftID(ctx context.Context, shiftID int64) ([]orders.Order, error) {
	return u.ordersRepo.GetOrdersByShiftID(ctx, shiftID)
}

func (u *Usecase) GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error) {
	return u.ordersRepo.GetCartItemsByOrderID(ctx, orderID)
}

func (u *Usecase) SetStatus(ctx context.Context, orderID int64, status int32) error {
	return u.ordersRepo.SetStatus(ctx, orderID, status)
}

func (u *Usecase) Accept(ctx context.Context, orderID int64) error {
	return u.ordersRepo.SetStatus(ctx, orderID, orders.StatusAccepted)
}

func (u *Usecase) CallDelivery(ctx context.Context, orderID int64) error {
	return u.ordersRepo.SetStatus(ctx, orderID, orders.StatusCooked)
}

func (u *Usecase) PickUp(ctx context.Context, orderID int64) error {
	return u.ordersRepo.SetStatus(ctx, orderID, orders.StatusInDelivery)
}

func (u *Usecase) Deliver(ctx context.Context, orderID int64) error {
	shiftID, err := u.ordersRepo.GetShiftIDByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	order, err := u.ordersRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	err = u.shiftsRepo.AddToTotalProfit(ctx, shiftID, order.TotalCost)
	if err != nil {
		return err
	}
	return u.ordersRepo.SetStatus(ctx, orderID, orders.StatusDelivered)
}

func (u *Usecase) Reject(ctx context.Context, orderID int64) error {
	return u.ordersRepo.SetStatus(ctx, orderID, orders.StatusRejected)
}

func (u *Usecase) GetOrderedDishesAndChefsByUserID(ctx context.Context, userID int64) ([]dishEntity.Dish, []chefEntity.Chef, error) {
	dishIDs, chefsIDs, err := u.ordersRepo.GetOrderedDishesIDsAndChefsIDs(ctx, userID, 100, 100)
	if err != nil {
		return nil, nil, err
	}
	dishes := make([]dishEntity.Dish, 0, len(dishIDs))
	for _, dishID := range dishIDs {
		dish, err := u.dishesUsecase.GetDishByID(ctx, dishID)
		if err != nil {
			return nil, nil, err
		}
		dishes = append(dishes, *dish)
	}
	chefs := make([]chefEntity.Chef, 0, len(chefsIDs))
	for _, chefID := range chefsIDs {
		chef, err := u.chefsUsecase.GetChefByID(ctx, chefID)
		if err != nil {
			return nil, nil, err
		}
		chefs = append(chefs, *chef)
	}
	return dishes, chefs, nil
}

func (u *Usecase) GetStatus(ctx context.Context, orderID int64) (int32, error) {
	return u.ordersRepo.GetStatus(ctx, orderID)
}

func (u *Usecase) GetOrderByID(ctx context.Context, orderID int64) (*orders.Order, error) {
	return u.ordersRepo.GetOrderByID(ctx, orderID)
}
func (u *Usecase) GetOrdersByUserID(ctx context.Context, userID int64) ([]orders.OrderProfile, error) {
	// 1) Получаем все заказы пользователя
	ordersList, err := u.ordersRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for user %d: %w", userID, err)
	}

	// 2) Подготавливаем срез под профили
	profiles := make([]orders.OrderProfile, 0, len(ordersList))

	// 3) Для каждого заказа по порядку собираем профиль
	for _, ord := range ordersList {
		// 3a) Позиции заказа
		items, err := u.ordersRepo.GetCartItemsByOrderID(ctx, ord.ID)
		if err != nil {
			return nil, fmt.Errorf("GetCartItemsByOrderID(order=%d): %w", ord.ID, err)
		}
		for i, item := range items {
			rating, err := u.dishesUsecase.GetDishRatingByID(ctx, item.Dish.ID)
			if err != nil {
				return nil, fmt.Errorf("GetDishRatingByID(dish=%d): %w", item.Dish.ID, err)
			}
			items[i].Dish.Rating = rating.Rating
			items[i].Dish.ReviewsCount = item.Dish.ReviewsCount
		}
		// 3b) Данные шефа
		chef, err := u.chefsUsecase.GetChefByID(ctx, ord.ChefID)
		if err != nil {
			return nil, fmt.Errorf("GetChefByID(%d): %w", ord.ChefID, err)
		}

		// 3c) Отзыв пользователя по этому заказу
		rv, err := u.reviewUsecase.GetReviewByOrderAndUserID(ctx, ord.ID, userID)
		var reviewDetail *orders.ReviewDetail
		if err != nil {
			// если отзыва нет — можно писать
			reviewDetail = &orders.ReviewDetail{CanWriteReview: true}
		} else {
			// отзыв есть — нельзя писать повторно
			stars := int(rv.Stars)
			comment := rv.Comment
			reviewDetail = &orders.ReviewDetail{
				CanWriteReview: false,
				Rating:         &stars,
				Comment:        &comment,
			}
		}

		profiles = append(profiles, orders.OrderProfile{
			Order:  &ord,
			Items:  items,
			Chef:   chef,
			Review: reviewDetail,
		})
	}

	// 4) Возвращаем готовые профили
	return profiles, nil
}

func (u *Usecase) GetActiveOrdersByUserID(ctx context.Context, userID int64) ([]orders.Order, error) {
	allOrders, err := u.ordersRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for user %d: %w", userID, err)
	}
	activeOrders := make([]orders.Order, 0, len(allOrders))
	for _, order := range allOrders {
		if order.Status == orders.StatusCooked ||
			order.Status == orders.StatusInDelivery ||
			order.Status == orders.StatusAccepted ||
			order.Status == orders.StatusCreated {
			activeOrders = append(activeOrders, order)
		}
	}
	return activeOrders, nil
}

func (u *Usecase) CountDishInOrders(ctx context.Context, dishID int64) (int, error) {
	return u.ordersRepo.CountDishInOrders(ctx, dishID)
}
