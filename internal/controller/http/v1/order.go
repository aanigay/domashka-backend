package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	cartentity "domashka-backend/internal/entity/cart"
	"domashka-backend/internal/utils/pointers"
)

type orderHandler struct {
	geoUsecase    geoUsecase
	cartUsecase   cartUsecase
	orderUsecase  orderUsecase
	shiftsUsecase shiftsUsecase
	chefUsecase   chefUsecase
}

func RegisterOrderHandlers(
	rg *gin.RouterGroup,
	geoUsecase geoUsecase,
	cartUsecase cartUsecase,
	orderUsecase orderUsecase,
	shiftUsecase shiftsUsecase,
	chefUsecase chefUsecase,
) {
	c := orderHandler{
		geoUsecase:    geoUsecase,
		cartUsecase:   cartUsecase,
		orderUsecase:  orderUsecase,
		shiftsUsecase: shiftUsecase,
		chefUsecase:   chefUsecase,
	}

	rg.GET("/chef/home", c.chefMain)
	rg = rg.Group("/order")
	rg.GET("/details_form", c.GetDetailsForm)
	rg.GET("/final_form", c.GetFinalForm)
	rg.POST("/create", c.createOrder)
	rg.GET("/")
	rg.POST("/accept", c.accept)
	rg.POST("/reject", c.reject)
	rg.POST("/call_delivery", c.callDelivery)
	rg.POST("/pickup", c.pickUp)
	rg.POST("/deliver", c.deliver)
	rg.GET("/status", c.getStatus)
}

type GetOrderDetailsFormData struct {
	Address        *Address        `json:"address,omitempty"`
	TotalPrice     Price           `json:"total"`
	PaymentOptions []PaymentOption `json:"payment_options"`
	DeliveryCost   Price           `json:"delivery_cost"`
}

type PaymentOption struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

func (h *orderHandler) GetDetailsForm(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid user ID",
				Details: "Передан некорректный ID пользователя.",
			},
		})
		return
	}
	address, err := h.geoUsecase.GetLastUpdatedClientAddress(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: "Ошибка на сервере. Попробуйте позже.",
			},
		})
		return
	}
	cartItems, err := h.cartUsecase.GetCartItems(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: "Ошибка на сервере. Попробуйте позже.",
			},
		})
		return
	}
	totalCartPrice := cartentity.GetTotalCartPrice(cartItems)
	var addressResp *Address
	if address != nil {
		addressResp = &Address{
			ID:    address.ID,
			Title: pointers.From(address.Address),
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": GetOrderDetailsFormData{
			Address: addressResp,
			TotalPrice: Price{
				Value:    totalCartPrice.Value,
				Currency: totalCartPrice.Currency,
			},
			DeliveryCost: Price{Value: 169, Currency: totalCartPrice.Currency},
			PaymentOptions: []PaymentOption{
				{
					ID:    1234,
					Type:  "MASTERCARD",
					Title: "1234",
				},
			},
		},
	})
}

type GetOrderFinalFormData struct {
	Dishes          []Dish           `json:"dishes"`
	DeliveryOptions []DeliveryOption `json:"delivery_options"`
	Total           Price            `json:"total"`
	DeliveryCost    Price            `json:"delivery_cost"`
}

type DeliveryOption struct {
	Title       string       `json:"title"`
	Date        string       `json:"date"`
	TimeOptions []TimeOption `json:"time_options"`
}

type TimeOption struct {
	Title         string `json:"title"`
	IntervalStart string `json:"interval_start"`
	IntervalEnd   string `json:"interval_end"`
}

func (h *orderHandler) GetFinalForm(c *gin.Context) {
	ctx := c.Request.Context()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userId == 400 {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid dish ID",
				Details: "Передан некорректный ID пользователя.",
			},
		})
		return
	}
	if userId == 404 {
		c.JSON(http.StatusNotFound, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4001,
				Message: "User not found",
				Details: "Пользователя с указанным ID не найдено.",
			},
		})
		return
	}
	if userId == 500 {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: "Ошибка на сервере. Попробуйте позже.",
			},
		})
		return
	}
	// TODO: временные интервалы брать из отдельной таблицы

	cartItems, err := h.cartUsecase.GetCartItems(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: "Ошибка на сервере. Попробуйте позже. " + err.Error(),
			},
		})
		return
	}
	dishes := make([]Dish, 0)
	for _, item := range cartItems {
		dish := Dish{
			DishID:     item.Dish.ID,
			CartItemID: pointers.To(item.ID),
			Title:      item.Dish.Name,
			Details:    item.GetDetailsString(),
			Price: &Price{
				Value:    item.Size.PriceValue,
				Currency: item.Size.PriceCurrency,
			},
			Quantity: item.Quantity,
			TotalPrice: Price{
				Value:    item.GetTotalPrice().Value,
				Currency: item.GetTotalPrice().Currency,
			},
			ImageURL: item.Dish.ImageURL,
		}
		dishes = append(dishes, dish)
	}

	totalPrice := cartentity.GetTotalCartPrice(cartItems)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": GetOrderFinalFormData{
			Dishes: dishes,
			DeliveryOptions: []DeliveryOption{
				{
					Title: "Сегодня",
					Date:  time.Now().Format("2006-01-02"),
					TimeOptions: []TimeOption{
						{
							Title:         "30-40 минут",
							IntervalStart: time.Now().Add(time.Minute * 30).Format("2006-01-02T15:04:05"),
							IntervalEnd:   time.Now().Add(time.Minute * 40).Format("2006-01-02T15:04:05"),
						},
						{
							Title:         "16:30-17:30",
							IntervalStart: time.Now().Add(time.Minute * 30).Format("2006-01-02T15:04:05"),
							IntervalEnd:   time.Now().Add(time.Minute * 40).Format("2006-01-02T15:04:05"),
						},
						{
							Title:         "17:30-18:30",
							IntervalStart: time.Now().Add(time.Minute * 30).Format("2006-01-02T15:04:05"),
							IntervalEnd:   time.Now().Add(time.Minute * 40).Format("2006-01-02T15:04:05"),
						},
					},
				},
				{
					Title: "Завтра",
					Date:  time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					TimeOptions: []TimeOption{
						{
							Title:         "16:30-17:30",
							IntervalStart: time.Now().Add(time.Minute * 30).Format("2006-01-02T15:04:05"),
							IntervalEnd:   time.Now().Add(time.Minute * 40).Format("2006-01-02T15:04:05"),
						},
						{
							Title:         "17:30-18:30",
							IntervalStart: time.Now().Add(time.Minute * 30).Format("2006-01-02T15:04:05"),
							IntervalEnd:   time.Now().Add(time.Minute * 40).Format("2006-01-02T15:04:05"),
						},
					},
				},
			},
			Total: Price{
				Value:    totalPrice.Value,
				Currency: totalPrice.Currency,
			},
			DeliveryCost: Price{
				Value:    169.0,
				Currency: "RUB",
			},
		},
	})
}

type CreateOrderRequest struct {
	UserID         int64 `json:"user_id"`
	LeaveByTheDoor bool  `json:"leave_by_the_door"`
	CallBeforehand bool  `json:"call_beforehand"`
}

func (h *orderHandler) createOrder(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	orderID, err := h.orderUsecase.CreateOrder(ctx, req.UserID, req.LeaveByTheDoor, req.CallBeforehand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"order_id": orderID,
		},
	})
}

func (h *orderHandler) chefMain(c *gin.Context) {
	ctx := c.Request.Context()
	chefID, err := strconv.ParseInt(c.Query("chef_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid user ID",
				Details: "Передан некорректный ID пользователя.",
			},
		})
		return
	}
	response := map[string]interface{}{}
	activeShift, err := h.shiftsUsecase.GetActiveShiftByChefID(ctx, chefID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	if activeShift == nil {
		response["is_active"] = false
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   response,
		})
		return
	}
	response["shift_id"] = activeShift.ID
	response["is_active"] = true
	response["elapsed"] = fmt.Sprintf("%.0f мин", time.Since(activeShift.CreatedAt).Minutes())
	response["total_profit"] = fmt.Sprintf("%.0f", activeShift.TotalProfit)
	// get orders by shift id
	o, err := h.orderUsecase.GetOrdersByShiftID(ctx, activeShift.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}

	// for each order get cartItems
	type kek struct {
		orderID   int64
		items     []cartentity.CartItem
		totalCost float32
		status    int32
	}
	orderCartItemsMap := map[int64]kek{}
	for _, order := range o {
		cartItems, err := h.orderUsecase.GetCartItemsByOrderID(ctx, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
					Details: err.Error(),
				},
			})
			return
		}
		orderCartItemsMap[order.ID] = kek{
			orderID:   order.ID,
			items:     cartItems,
			totalCost: cartentity.GetTotalCartPrice(cartItems).Value,
			status:    order.Status,
		}
	}
	ordersByStatusMap := map[int32][]kek{}
	for _, k := range orderCartItemsMap {
		ordersByStatusMap[k.status] = append(ordersByStatusMap[k.status], k)
	}
	var orderGroups []map[string]interface{}
	for status, os := range ordersByStatusMap {
		group := map[string]interface{}{}
		group["order_status"] = status
		ordersResp := []map[string]interface{}{}
		for _, o := range os {
			order := map[string]interface{}{}
			order["order_id"] = o.orderID
			order["total_cost"] = o.totalCost
			dishes := []map[string]interface{}{}
			for _, item := range o.items {
				dish := map[string]interface{}{}
				dish["dish_id"] = item.Dish.ID
				dish["name"] = item.Dish.Name
				dish["details"] = item.GetDetailsString()
				dish["image_url"] = item.Dish.ImageURL
				dish["price"] = item.Size.PriceValue
				dish["quantity"] = item.Quantity
				dishes = append(dishes, dish)
			}
			order["dishes"] = dishes
			ordersResp = append(ordersResp, order)
		}
		group["orders"] = ordersResp
		orderGroups = append(orderGroups, group)
	}
	response["order_groups"] = orderGroups
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

func (h *orderHandler) accept(c *gin.Context) {
	ctx := c.Request.Context()
	orderID, err := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid order ID",
				Details: "Передан некорректный ID заказа.",
			},
		})
		return
	}
	err = h.orderUsecase.Accept(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *orderHandler) callDelivery(c *gin.Context) {
	ctx := c.Request.Context()
	orderID, err := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid order ID",
				Details: "Передан некорректный ID заказа.",
			},
		})
		return
	}
	err = h.orderUsecase.CallDelivery(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *orderHandler) pickUp(c *gin.Context) {
	ctx := c.Request.Context()
	orderID, err := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid order ID",
				Details: "Передан некорректный ID заказа.",
			},
		})
		return
	}
	err = h.orderUsecase.PickUp(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *orderHandler) deliver(c *gin.Context) {
	ctx := c.Request.Context()
	orderID, err := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid order ID",
				Details: "Передан некорректный ID заказа.",
			},
		})
		return
	}
	err = h.orderUsecase.Deliver(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *orderHandler) reject(c *gin.Context) {
	ctx := c.Request.Context()
	orderID, err := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid order ID",
				Details: "Передан некорректный ID заказа.",
			},
		})
		return
	}
	err = h.orderUsecase.Reject(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *orderHandler) getStatus(c *gin.Context) {
	ctx := c.Request.Context()
	orderID, err := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid order ID",
				Details: "Передан некорректный ID заказа.",
			},
		})
		return
	}
	order, err := h.orderUsecase.GetOrderByID(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	// получить адрес клиента
	clientAddress, err := h.geoUsecase.GetAddressByID(ctx, order.ClientAddressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	// получить адрес повара
	chefAddress, err := h.geoUsecase.GetChefAddress(ctx, int(order.ChefID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}

	// получить айтемы заказа
	items, err := h.orderUsecase.GetCartItemsByOrderID(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	dishes := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		chef, err := h.chefUsecase.GetChefByDishID(ctx, item.Dish.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
					Details: err.Error(),
				},
			})
			return
		}
		dish := map[string]interface{}{
			"dish_id":    item.Dish.ID,
			"name":       item.Dish.Name,
			"quantity":   item.Quantity,
			"image_url":  item.Dish.ImageURL,
			"avatar_url": chef.SmallImageURL,
			"price":      item.Size.PriceValue,
			"chef_id":    item.Dish.ChefID,
		}
		dishes = append(dishes, dish)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"dishes": dishes,
			"status": order.Status,
			"client_address": gin.H{
				"address":   clientAddress.Address,
				"longitude": clientAddress.Longitude,
				"latitude":  clientAddress.Latitude,
			},
			"chef_address": gin.H{
				"longitude": chefAddress.Longitude,
				"latitude":  chefAddress.Latitude,
			},
			"leave_by_the_door": order.LeaveByTheDoor,
			"comment":           clientAddress.Comment,
		},
	})
}
