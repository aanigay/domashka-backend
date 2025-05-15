package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type searchHandler struct {
	dishesUsecase  dishesUsecase
	chefUsecase    chefUsecase
	orderUsecase   orderUsecase
	reviewsUsecase reviewsUsecase
	usersUsecase   usersUsecase
}

func RegisterSearchHandler(
	rg *gin.RouterGroup,
	dishesUsecase dishesUsecase,
	chefUsecase chefUsecase,
	orderUsecase orderUsecase,
	reviewsUsecase reviewsUsecase,
	usersUsecase usersUsecase,
) {
	h := searchHandler{
		dishesUsecase:  dishesUsecase,
		chefUsecase:    chefUsecase,
		orderUsecase:   orderUsecase,
		reviewsUsecase: reviewsUsecase,
		usersUsecase:   usersUsecase,
	}

	rg.GET("/search", h.getAll)
}

func (h *searchHandler) getAll(c *gin.Context) {
	ctx := c.Request.Context()
	userIDstr := c.GetString("user_id")
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if userIDstr == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Неверный юзер айди",
			"message": "Пожалуйста, авторизуйтесь заново.",
		})
		return
	}
	// получаем блюда, заказанные ранее
	// получаем поваров, заказанных ранее
	dishes, chefs, err := h.orderUsecase.GetOrderedDishesAndChefsByUserID(ctx, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Не удалось получить данные",
			"message": "Пожалуйста, авторизуйтесь заново.",
		})
		return
	}
	type chefSnippet struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Rating    string `json:"rating"`
	}
	type dishSnippet struct {
		ID            int64    `json:"id"`
		ChefID        int64    `json:"chef_id"`
		Name          string   `json:"name"`
		ImageURL      string   `json:"image_url"`
		ChefAvatarURL string   `json:"chef_avatar_url"`
		Badges        []string `json:"badges"`
		MinPrice      string   `json:"price"`
		Rating        *float32 `json:"rating,omitempty"`
	}
	type chefHome struct {
	}

	orderedChefSnippets := make([]chefSnippet, 0, len(chefs))
	chefAvatars := map[int64]string{}
	for _, chef := range chefs {
		orderedChefSnippets = append(orderedChefSnippets, chefSnippet{
			ID:        chef.ID,
			Name:      chef.Name,
			AvatarURL: chef.SmallImageURL,
			Rating:    fmt.Sprintf("%.1f (%d)", *chef.Rating, *chef.ReviewsCount),
		})
		chefAvatars[chef.ID] = chef.SmallImageURL
	}
	orderedDishSnippets := make([]dishSnippet, 0, len(dishes))
	for _, dish := range dishes {
		minPrice, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, dish.ID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Не удалось получить данные",
				"message": "Пожалуйста, авторизуйтесь заново.",
			})
			return

		}
		category, err := h.dishesUsecase.GetCategoryTitleByDishID(ctx, dish.ID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Не удалось получить данные",
				"message": err,
			})
		}
		orderedDishSnippets = append(orderedDishSnippets, dishSnippet{
			ID:            dish.ID,
			ChefID:        dish.ChefID,
			Name:          dish.Name,
			ImageURL:      dish.ImageURL,
			ChefAvatarURL: chefAvatars[dish.ChefID],
			MinPrice:      fmt.Sprintf("от %.0fр", minPrice.Value),
			Badges:        []string{category},
			Rating:        dish.Rating,
		})
	}

	chefs, err = h.chefUsecase.GetAll(ctx)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Не удалось получить данные",
			"message": "Пожалуйста, авторизуйтесь заново.",
		})
		return
	}
	allChefs := []chefSnippet{}
	for _, chef := range chefs {
		cs := chefSnippet{
			ID:        chef.ID,
			Name:      chef.Name,
			AvatarURL: chef.SmallImageURL,
			Rating:    fmt.Sprintf("%.1f (%d)", *chef.Rating, *chef.ReviewsCount),
		}
		allChefs = append(allChefs, cs)
		chefAvatars[chef.ID] = chef.SmallImageURL
	}
	// получаем все блюда
	dishes, err = h.dishesUsecase.GetAll(ctx)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Не удалось получить данные",
			"message": "Пожалуйста, авторизуйтесь заново.",
		})
		return
	}
	allDishes := []dishSnippet{}
	chefsDishes := map[int64][]map[string]interface{}{}
	for _, dish := range dishes {
		minPrice, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, dish.ID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Не удалось получить данные",
				"message": "Пожалуйста, авторизуйтесь заново.",
			})
			return
		}
		category, err := h.dishesUsecase.GetCategoryTitleByDishID(ctx, dish.ID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Не удалось получить данные",
				"message": err,
			})
			return
		}
		allDishes = append(allDishes, dishSnippet{
			ID:            dish.ID,
			ChefID:        dish.ChefID,
			Name:          dish.Name,
			ImageURL:      dish.ImageURL,
			ChefAvatarURL: chefAvatars[dish.ChefID],
			MinPrice:      fmt.Sprintf("от %.0fр", minPrice.Value),
			Badges:        []string{category},
		})
		chefsDishes[dish.ChefID] = append(chefsDishes[dish.ChefID], map[string]interface{}{
			"dish_id":        dish.ID,
			"dish_image_url": dish.ImageURL,
			"name":           dish.Name,
			"price":          fmt.Sprintf("от %.0fр", minPrice.Value),
			"rating":         dish.Rating,
		})
	}
	var chefsResponse []map[string]interface{}
	for _, chef := range chefs {
		certifications, err := h.chefUsecase.GetChefCertifications(ctx, chef.ID)
		if err != nil {
			log.Println(err)
		}
		badges := []string{}
		for _, certification := range certifications {
			badges = append(badges, certification.Name)
		}
		r, err := h.reviewsUsecase.GetReviewsByChefID(ctx, chef.ID)
		if err != nil {
			log.Println(err)
		}
		var review map[string]interface{}
		if len(r) > 0 {
			user, err := h.usersUsecase.GetByID(ctx, r[0].UserID)
			if err != nil {
				log.Println(err)
				review = map[string]interface{}{
					"name": user.ID,
					"text": r[0].Comment,
				}
			}
			review = map[string]interface{}{
				"name": user.FirstName,
				"text": r[0].Comment,
			}
		}
		chefsResponse = append(chefsResponse, map[string]interface{}{
			"chef_info": map[string]interface{}{
				"id":            chef.ID,
				"name":          chef.Name,
				"avatar_url":    chef.SmallImageURL,
				"rating_string": fmt.Sprintf("%.1f (%d)", *chef.Rating, *chef.ReviewsCount),
				"badges":        badges,
			},
			"review": review,
			"dishes": chefsDishes[chef.ID],
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"dishes": allDishes,
			"chefs":  chefsResponse,
			"ordered": gin.H{
				"dishes": orderedDishSnippets,
				"chefs":  orderedChefSnippets,
			},
		},
	})
}
