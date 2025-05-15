package v1

import (
	chefEntity "domashka-backend/internal/entity/chefs"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	chefsDistanceRange = 15
	defaultLong        = 37.623150
	defaultLat         = 55.752507
)

type homeHandler struct {
	geoUsecase
	dishesUsecase
	chefUsecase
	orderUsecase
	reviewsUsecase
}

func NewHomeHandler(rg *gin.RouterGroup, jwt jwtUsecase, geo geoUsecase, dishes dishesUsecase, chef chefUsecase, order orderUsecase, reviews reviewsUsecase) {
	hh := homeHandler{
		geoUsecase:     geo,
		dishesUsecase:  dishes,
		chefUsecase:    chef,
		orderUsecase:   order,
		reviewsUsecase: reviews,
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
		chef, err := h.chefUsecase.GetChefByID(ctx, dish.ChefID)
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
		category, err := h.dishesUsecase.GetCategoryTitleByDishID(ctx, dish.ID)
		if err != nil {
			fmt.Println(err)
		}
		suggests := []string{category}
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
			"chef_avatar_url": chef.SmallImageURL,
			"chef_id":         dish.ChefID,
			"suggests":        suggests,
			"name":            dish.Name,
			"rating":          dish.Rating,
			"price":           fmt.Sprintf("от %.0fр", minPrice.Value),
		})
	}
	response["top_dishes"] = dishes
	var nearestChefs []chefEntity.Chef
	if address == nil {
		nearestChefs, err = h.chefUsecase.GetNearestChefs(ctx, defaultLat, defaultLong, chefsDistanceRange, 6)
	} else {
		nearestChefs, err = h.chefUsecase.GetNearestChefs(ctx, address.Latitude, address.Longitude, chefsDistanceRange, 6)
	}
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
	for _, topChef := range nearestChefs {
		chefInfo := make(map[string]interface{})
		chefInfo["id"] = topChef.ID
		chefInfo["name"] = topChef.Name
		if topChef.Rating != nil && topChef.ReviewsCount != nil {
			chefInfo["rating_string"] = fmt.Sprintf("%.1f (%d)", *topChef.Rating, *topChef.ReviewsCount)
		}

		var certs []chefEntity.Certification

		certs, err := h.chefUsecase.GetChefCertifications(ctx, topChef.ID)
		if err != nil {
			chefInfo["badges"] = []string{}
		} else {
			certNames := make([]string, len(certs))
			for i, cert := range certs {
				certNames[i] = cert.Name
			}
			chefInfo["badges"] = certNames
		}

		chefInfo["avatar_url"] = topChef.SmallImageURL
		// Макс кол-во отзывов на странице
		limit := 1
		rawReviews, err := h.reviewsUsecase.GetFullReviewsByChefID(ctx, topChef.ID, limit)

		var review map[string]interface{}
		if err != nil || len(rawReviews) == 0 {
			review = map[string]interface{}{
				"name": "",
				"text": "",
			}
		} else {
			rv := rawReviews[0]
			review = map[string]interface{}{
				"name": rv.UserName,
				"text": rv.Comment,
			}
		}
		chefDishesMap := make([]map[string]interface{}, 0)
		chefDishes, err := h.dishesUsecase.GetDishesByChefID(ctx, topChef.ID, 6)
		if err != nil {
			chefDishesMap = make([]map[string]interface{}, 0)
			fmt.Println(err)
			continue
		}
		for _, dish := range chefDishes {
			if dish.IsDeleted {
				continue
			}
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
	activeOrders, err := h.orderUsecase.GetActiveOrdersByUserID(ctx, userID)
	if err != nil {
		fmt.Println(err)
	} else {
		ordersResp := make([]map[string]interface{}, 0, len(activeOrders))
		for _, order := range activeOrders {
			ordersResp = append(ordersResp, map[string]interface{}{
				"id":     order.ID,
				"status": order.Status,
			})
		}
		response["active_orders"] = ordersResp
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
