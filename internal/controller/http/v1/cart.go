package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "domashka-backend/docs"
	"domashka-backend/internal/custom_errors"
	entity "domashka-backend/internal/entity/carts"
)

type cartHandler struct {
	cartUsecase cartUsecase
}

func registerCartHandler(rg *gin.RouterGroup, log logger, cartUsecase cartUsecase) {
	c := cartHandler{cartUsecase: cartUsecase}
	rg = rg.Group("/cart")
	{
		rg.POST("/create", c.createCart)
		rg.POST("/clear", c.clearCart)
		rg.GET("/:user_id", c.getCart)
		rg.POST("/add_item", c.addCartItem)
		rg.POST("/update_item", c.updateCartItem)
		rg.POST("/remove_item", c.removeCartItem)
		rg.GET("/:user_id/items", c.getCartItems)
	}
}

type CreateCartRequest struct {
	UserID string `json:"user_id"`
}

// createCart godoc
// @Summary      Создание корзины
// @Description  Создание корзины
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request   body      CreateCartRequest  true  "body"
// @Router       /cart/create [post]
func (c *cartHandler) createCart(g *gin.Context) {
	ctx := g.Request.Context()

	var req CreateCartRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.cartUsecase.CreateCart(ctx, &entity.Cart{
		UserID: req.UserID,
	})
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

// getCart godoc
// @Summary      Получение корзины
// @Description  Получение корзины
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "user id"
// @Router       /cart/{id} [get]
func (c *cartHandler) getCart(g *gin.Context) {
	ctx := g.Request.Context()

	userID := g.Param("user_id")
	if userID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	cart, err := c.cartUsecase.GetCart(ctx, userID)
	if err != nil {
		if errors.Is(err, custom_errors.ErrCartNotFound) {
			g.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, cart)
}

type AddCartItemRequest struct {
	UserID                   string  `json:"user_id"`
	DishID                   string  `json:"dish_id"`
	ChefID                   string  `json:"chef_id"`
	AdditionalIngredientsIDs []int64 `json:"additional_ingredients_ids"`
	RemovedIngredientsIDs    []int64 `json:"removed_ingredients_ids"`
	CustomerNotes            string  `json:"customer_notes,omitempty"`
}

// addCartItem godoc
// @Summary      Добавление товара в корзину
// @Description  Добавление товара в корзину
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request   body      AddCartItemRequest  true  "body"
// @Router       /cart/add_item [post]
func (c *cartHandler) addCartItem(g *gin.Context) {
	ctx := g.Request.Context()

	var req AddCartItemRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemID, err := c.cartUsecase.AddCartItem(ctx, &entity.CartItem{
		UserID:                   req.UserID,
		DishID:                   req.DishID,
		ChefID:                   req.ChefID,
		AdditionalIngredientsIDs: req.AdditionalIngredientsIDs,
		RemovedIngredientsIDs:    req.RemovedIngredientsIDs,
		CustomerNotes:            req.CustomerNotes,
	})
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"item_id": itemID})

}

type UpdateCartItemRequest struct {
	ID                       string  `json:"id"`
	UserID                   string  `json:"user_id"`
	DishID                   string  `json:"dish_id"`
	ChefID                   string  `json:"chef_id"`
	AdditionalIngredientsIDs []int64 `json:"additional_ingredients_ids"`
	RemovedIngredientsIDs    []int64 `json:"removed_ingredients_ids"`
	CustomerNotes            string  `json:"customer_notes,omitempty"`
}

// updateCartItem godoc
// @Summary      Обновление товара в корзине
// @Description  Обновление товара в корзине
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request   body      UpdateCartItemRequest  true  "body"
// @Router       /cart/update_item [post]
func (c *cartHandler) updateCartItem(g *gin.Context) {
	ctx := g.Request.Context()

	var req UpdateCartItemRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := c.cartUsecase.UpdateCartItem(ctx, &entity.CartItem{
		ID:                       req.ID,
		UserID:                   req.UserID,
		DishID:                   req.DishID,
		ChefID:                   req.ChefID,
		AdditionalIngredientsIDs: req.AdditionalIngredientsIDs,
		RemovedIngredientsIDs:    req.RemovedIngredientsIDs,
		CustomerNotes:            req.CustomerNotes,
	})
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "success"})
}

type ClearCartRequest struct {
	UserID string `json:"user_id"`
}

// clearCart godoc
// @Summary      Очищение корзины
// @Description  Очищение корзины
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request   body      ClearCartRequest  true  "body"
// @Router       /cart/clear [post]
func (c *cartHandler) clearCart(g *gin.Context) {
	ctx := g.Request.Context()

	var req ClearCartRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.cartUsecase.ClearCart(ctx, req.UserID)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "success"})
}

type RemoveCartItemRequest struct {
	ID string `json:"id"`
}

// removeCartItem godoc
// @Summary      Удаление товара из корзины
// @Description  Удаление товара из корзины
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request   body      RemoveCartItemRequest  true  "body"
// @Router       /cart/remove_item [post]
func (c *cartHandler) removeCartItem(g *gin.Context) {
	ctx := g.Request.Context()

	var req RemoveCartItemRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.cartUsecase.DeleteCartItem(ctx, req.ID)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "success"})
}

// getCartItems godoc
// @Summary      Получение всех товаров в корзине
// @Description  Получение всех товаров в корзине
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "user id"
// @Router       /cart/{id}/items [get]
func (c *cartHandler) getCartItems(g *gin.Context) {
	ctx := g.Request.Context()

	userID := g.Param("user_id")
	if userID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	items, err := c.cartUsecase.GetCartItems(ctx, userID)
	if err != nil {
		if errors.Is(err, custom_errors.ErrCartNotFound) {
			g.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"items": items})
}
