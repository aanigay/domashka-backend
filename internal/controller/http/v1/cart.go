package v1

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	cartentity "domashka-backend/internal/entity/cart"
	"domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/utils/pointers"
)

type cartHandler struct {
	chefUsecase   chefUsecase
	cartUsecase   cartUsecase
	dishesUsecase dishesUsecase
	geoUsecase    geoUsecase
}

func RegisterCartHandlers(
	rg *gin.RouterGroup,
	chefUsecase chefUsecase,
	cartUsecase cartUsecase,
	dishesUsecase dishesUsecase,
	geoUsecase geoUsecase) {
	c := cartHandler{
		chefUsecase:   chefUsecase,
		cartUsecase:   cartUsecase,
		dishesUsecase: dishesUsecase,
		geoUsecase:    geoUsecase,
	}

	rg = rg.Group("/cart")
	rg.GET("/view", c.GetCartView)
	rg.POST("/add", c.AddItemToCart)
	rg.POST("/remove", c.RemoveItem)
	rg.POST("/clear", c.ClearCart)
	rg.POST("/increment", c.IncrementCartItem)
	rg.POST("/decrement", c.DecrementCartItem)
}

type GetCartViewSuccessResponseData struct {
	Address       *Address      `json:"address,omitempty"`
	DishGroups    []DishGroup   `json:"dish_groups"`
	RelatedDishes []dishSnippet `json:"related_dishes"`
	DeliveryCost  Price         `json:"delivery_cost"`
	TotalCost     Price         `json:"total"`
}

type Address struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type DishGroup struct {
	ChefID   int64  `json:"chef_id"`
	ChefName string `json:"chef_name"`
	ImageURL string `json:"image_url"`
	Dishes   []Dish `json:"dishes"`
}

type Dish struct {
	DishID     int64  `json:"dish_id"`
	CartItemID *int64 `json:"cart_item_id,omitempty"`
	Title      string `json:"title"`
	Details    string `json:"details"`
	Price      *Price `json:"price,omitempty"`
	Quantity   int32  `json:"quantity"`
	TotalPrice Price  `json:"total_price"`
	ImageURL   string `json:"image_url"`
}

func (h *cartHandler) GetCartView(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: "Не удалось получить адрес клиента.",
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
				Message: "Internal server error.",
				Details: "Не удалось получить товары в корзине.",
			},
		})
		return
	}
	chefDishes := make(map[int64][]cartentity.CartItem, 0)
	for _, cartItem := range cartItems {
		chefDishes[cartItem.Dish.ChefID] = append(chefDishes[cartItem.Dish.ChefID], cartItem)
	}

	dishGroups := make([]DishGroup, 0)
	allDishIDs := make([]int64, 0)
	for chefID, items := range chefDishes {
		chef, err := h.chefUsecase.GetChefByID(ctx, chefID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    4003,
					Message: "Internal server error.",
					Details: "Не удалось получить товары в корзине.",
				},
			})
			return
		}
		d := make([]Dish, 0)
		for _, item := range items {
			addIngredientsLabels := make([]string, 0, len(item.AddedIngredients))
			for _, addIngredient := range item.AddedIngredients {
				addIngredientsLabels = append(addIngredientsLabels, addIngredient.Name)
			}
			removeIngredientsLabels := make([]string, 0, len(item.AddedIngredients))
			for _, removeIngredient := range item.RemovedIngredients {
				removeIngredientsLabels = append(removeIngredientsLabels, removeIngredient.Name)
			}
			addIngredientsText := strings.Join(addIngredientsLabels, ", ")
			removeIngredientsText := strings.Join(removeIngredientsLabels, ", ")
			totalPrice := item.GetTotalPrice()
			details := fmt.Sprintf("%s", item.Size.Label)
			if addIngredientsText != "" {
				details = fmt.Sprintf("%s, Добавить: %s", details, addIngredientsText)
			}
			if removeIngredientsText != "" {
				details = fmt.Sprintf("%s, Убрать: %s", details, removeIngredientsText)
			}
			d = append(d, Dish{
				DishID:     item.Dish.ID,
				CartItemID: pointers.To(item.ID),
				Title:      item.Dish.Name,
				Details:    details,
				Price: &Price{
					Value:    item.Size.PriceValue,
					Currency: item.Size.PriceCurrency,
				},
				Quantity: item.Quantity,
				TotalPrice: Price{
					Value:    totalPrice.Value,
					Currency: totalPrice.Currency,
				},
				ImageURL: item.Dish.ImageURL,
			})
			allDishIDs = append(allDishIDs, item.Dish.ID)
		}
		dishGroups = append(dishGroups, DishGroup{
			ChefID:   chefID,
			ChefName: chef.Name,
			ImageURL: chef.SmallImageURL,
			Dishes:   d,
		})
	}
	// TODO: calculate delivery cost
	dishSnippets := make([]dishSnippet, 0)
	if len(dishGroups) > 0 {
		relatedDishes, err := h.dishesUsecase.GetDishesByChefID(ctx, dishGroups[0].ChefID, 6)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    4003,
					Message: "Internal server error.",
					Details: "Не удалось связанные товары.",
				},
			})
			return
		}
		for _, dish := range relatedDishes {
			if slices.Contains(allDishIDs, dish.ID) {
				continue
			}
			price, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, dish.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse{
					Status: "error",
					Err: errorMessage{
						Code:    4003,
						Message: "Internal server error.",
						Details: "Не удалось связанные товары.",
					},
				})
				return
			}
			dishSnippets = append(dishSnippets, dishSnippet{
				ID:       dish.ID,
				Name:     dish.Name,
				ImageURL: dish.ImageURL,
				Price:    fmt.Sprintf("от %.0fр", price.Value),
				Rating:   pointers.From(dish.Rating),
			})
		}
	}
	deliveryCost := Price{
		Value:    169,
		Currency: "RUB",
	}
	totalCartPrice := cartentity.GetTotalCartPrice(cartItems)
	var addressResp *Address
	if address != nil {
		addressResp = &Address{
			ID:    address.ID,
			Title: pointers.From(address.Address),
		}
	}
	resp := GetCartViewSuccessResponseData{
		Address:       addressResp,
		DishGroups:    dishGroups,
		RelatedDishes: dishSnippets,
		DeliveryCost:  deliveryCost,
		TotalCost: Price{
			Value:    totalCartPrice.Value,
			Currency: totalCartPrice.Currency,
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

type AddItemToCartRequest struct {
	UserID                   int64   `json:"user_id"`
	DishID                   int64   `json:"dish_id"`
	ChefID                   int64   `json:"chef_id"`
	SizeID                   int64   `json:"size_id"`
	AdditionalIngredientsIDs []int64 `json:"additional_ingredients_ids"`
	RemovedIngredientsIDs    []int64 `json:"removed_ingredients_ids"`
	Notes                    string  `json:"notes"`
}

func (h *cartHandler) AddItemToCart(c *gin.Context) {
	ctx := c.Request.Context()

	var req AddItemToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid request body.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	cartItemID, err := h.cartUsecase.AddItem(
		ctx,
		req.UserID,
		dishes.Dish{
			ID:     req.DishID,
			ChefID: req.ChefID,
		},
		req.SizeID,
		req.AdditionalIngredientsIDs,
		req.RemovedIngredientsIDs,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"cart_item_id": cartItemID,
		},
	})
}

type UpdateCartItemRequest struct {
	CartItemID               int64   `json:"cart_item_id"`
	UserID                   int64   `json:"user_id"`
	DishID                   int64   `json:"dish_id"`
	ChefID                   int64   `json:"chef_id"`
	SizeID                   int64   `json:"size_id"`
	AdditionalIngredientsIDs []int64 `json:"additional_ingredients_ids"`
	RemovedIngredientsIDs    []int64 `json:"removed_ingredients_ids"`
	Notes                    string  `json:"notes"`
}

func (h *cartHandler) UpdateCartItem(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid request body.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}

	// TODO:
	// Сделать здесь изменение айтема
	cartItemID, err := h.cartUsecase.AddItem(
		ctx,
		req.UserID,
		dishes.Dish{
			ID:     req.DishID,
			ChefID: req.ChefID,
		},
		req.SizeID,
		req.AdditionalIngredientsIDs,
		req.RemovedIngredientsIDs,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"cart_item_id": cartItemID,
		},
	})
}

type IncrementCartItemRequest struct {
	UserID     int64 `json:"user_id"`
	CartItemID int64 `json:"cart_item_id"`
}

func (h *cartHandler) IncrementCartItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req IncrementCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid request body.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	quantity, err := h.cartUsecase.IncrementCartItem(ctx, req.CartItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"new_quantity": quantity,
		},
	})
}

type DecrementCartItemRequest struct {
	UserID     int64 `json:"user_id"`
	CartItemID int64 `json:"cart_item_id"`
}

func (h *cartHandler) DecrementCartItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req DecrementCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid request body.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	quantity, err := h.cartUsecase.DecrementCartItem(ctx, req.CartItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"new_quantity": quantity,
		},
	})
}

type ClearCartRequest struct {
	UserID int64 `json:"user_id"`
}

func (h *cartHandler) ClearCart(c *gin.Context) {
	ctx := c.Request.Context()
	var req ClearCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid request body.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	err := h.cartUsecase.ClearCart(ctx, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *cartHandler) RemoveItem(c *gin.Context) {
	ctx := c.Request.Context()
	cartItemID, err := strconv.ParseInt(c.Query("cart_item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{})
		return
	}
	err = h.cartUsecase.RemoveItem(ctx, cartItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error.",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
