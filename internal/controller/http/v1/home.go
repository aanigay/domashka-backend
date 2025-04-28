package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type homeHandler struct {
	geoUsecase
	dishesUsecase
	chefUsecase
}

func NewHomeHandler(rg *gin.RouterGroup, jwt jwtUsecase, geo geoUsecase, dishes dishesUsecase, chef chefUsecase) {
	hh := homeHandler{
		geoUsecase:    geo,
		dishesUsecase: dishes,
		chefUsecase:   chef,
	}

	rg.GET("/home", AuthMiddleware(jwt), hh.getHomePage)
}

type ChefSnippet struct {
}

func (h *homeHandler) getHomePage(c *gin.Context) {
	ctx := c.Request.Context()
	// получить последний адрес
	userIDstr := c.GetString("user_id")
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if userIDstr == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Неверный юзер айди",
			"message": "Пожалуйста, авторизуйтесь заново.",
		})
		return
	}

	address, err := h.GetLastUpdatedClientAddress(ctx, userID)
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
	// получить блюда с лучшим рейтингом
	topDishes, err := h.dishesUsecase.GetTopDishes(ctx, 100)
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

	response := map[string]interface{}{}

	// получить для каждого блюда chef_avatar_url
	dishes := make([]map[string]interface{}, 0)
	for _, dish := range topDishes {
		avatarURL, err := h.chefUsecase.GetChefAvatarURLByChefID(ctx, dish.ChefID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    4003,
					Message: "Internal server error",
					Details: err.Error(),
				},
			})
			return
		}
		// TODO: добавить в саджесты категорию
		suggests := []string{"15-20 мин", "Веган"}
		minPrice, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, dish.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    4003,
					Message: "Internal server error",
					Details: err.Error(),
				},
			})
			return
		}
		dishes = append(dishes, map[string]interface{}{
			"id":              dish.ID,
			"dish_image_url":  dish.ImageURL,
			"chef_avatar_url": avatarURL,
			"chef_id":         dish.ChefID,
			"suggests":        suggests,
			"name":            dish.Name,
			"rating":          dish.Rating,
			"price":           fmt.Sprintf("от %.0fр", minPrice.Value),
		})
	}
	response["top_dishes"] = dishes
	topChefs, err := h.GetTopChefs(ctx, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: err.Error(),
			},
		})
		return
	}
	chefs := make([]map[string]interface{}, 0)
	for _, topChef := range topChefs {
		chefInfo := make(map[string]interface{})
		chefInfo["id"] = topChef.ID
		chefInfo["name"] = topChef.Name
		if topChef.Rating != nil && topChef.ReviewsCount != nil {
			chefInfo["rating_string"] = fmt.Sprintf("%.1f (%d)", *topChef.Rating, *topChef.ReviewsCount)
		}
		chefInfo["suggests"] = []string{"15-20 мин", "Мед. книжка"}
		chefInfo["avatar_url"] = topChef.ImageURL
		review := make(map[string]interface{})
		review["name"] = "Ульяна"
		review["text"] = "Очень понравилось! Обязательно закажем еще"

		chefDishesMap := make([]map[string]interface{}, 0)
		chefDishes, err := h.dishesUsecase.GetDishesByChefID(ctx, topChef.ID)
		if err != nil {
			chefDishesMap = make([]map[string]interface{}, 0)
			fmt.Println(err)
			continue
		}
		for _, dish := range chefDishes {
			minPrice, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, dish.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			chefDishesMap = append(chefDishesMap, map[string]interface{}{
				"dish_id":        dish.ID,
				"dish_image_url": dish.ImageURL,
				"name":           dish.Name,
				"price":          fmt.Sprintf("от %.0fр", minPrice.Value),
				"rating":         dish.Rating,
			})
		}
		chefs = append(chefs, map[string]interface{}{
			"chef_info": chefInfo,
			"dishes":    chefDishesMap,
			"review":    review,
		})
	}
	response["chefs"] = chefs
	if address != nil {
		response["address"] = address.Address
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   response,
		"status": "success",
	})
}
