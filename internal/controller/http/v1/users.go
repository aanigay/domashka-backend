package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"domashka-backend/internal/custom_errors"
	"domashka-backend/internal/utils/types"

	chefEntity "domashka-backend/internal/entity/chefs"
	dishesEntity "domashka-backend/internal/entity/dishes"
	orderEntity "domashka-backend/internal/entity/orders"
	usersEntity "domashka-backend/internal/entity/users"
)

type usersHandler struct {
	log             logger
	usecase         usersUsecase
	orderUsecase    orderUsecase
	reviewUsecase   reviewsUsecase
	favoriteUsecase favoritesUsecase
	dishesUsecase   dishesUsecase
	chefUsecase     chefUsecase
}

func newUsersHandler(rg *gin.RouterGroup,
	log logger,
	uu usersUsecase,
	ru reviewsUsecase,
	ou orderUsecase,
	fu favoritesUsecase,
	dishesUsecase dishesUsecase,
	chefUsecase chefUsecase,
) {
	u := usersHandler{
		log:             log,
		usecase:         uu,
		orderUsecase:    ou,
		reviewUsecase:   ru,
		favoriteUsecase: fu,
		dishesUsecase:   dishesUsecase,
		chefUsecase:     chefUsecase,
	}

	rg = rg.Group("/users")
	{
		rg.POST("/create", u.Create)
		rg.GET("/:id", u.GetProfile)
		rg.POST("/:id/update", u.Update)
		rg.DELETE("/:id", u.Delete)
		rg.POST("/:id/favorite/chef", u.AddFavoriteChef)
		rg.DELETE("/:id/favorite/chef", u.RemoveFavoriteChef)
		rg.POST("/:id/favorite/dish", u.AddFavoriteDish)
		rg.DELETE("/:id/favorite/dish", u.RemoveFavoriteDish)
	}
}

type orderDetail struct {
	OrderID int64                 `json:"order_id"`
	Label   string                `json:"label"`
	Dishes  []dishDetail          `json:"dishes"`
	Chef    chefSnippet           `json:"chef"`
	Review  reviewResponseProfile `json:"review"`
}

type dishDetail struct {
	DishID    int64    `json:"dish_id"`
	Title     string   `json:"title"`
	ImageURL  string   `json:"image_url"`
	Rating    *float32 `json:"rating,omitempty"`
	Details   *string  `json:"details,omitempty"`
	PriceFrom Price    `json:"price_from"`
}

type priceResponse struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

type chefSnippet struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	AvatarURL    string   `json:"avatar_url"`
	Rating       *float32 `json:"rating,omitempty"`
	ReviewsCount *int32   `json:"reviews_count,omitempty"`
}

type reviewResponseProfile struct {
	CanWriteReview bool   `json:"can_write_review"`
	Rating         int    `json:"rating"`
	Comment        string `json:"comment"`
}

// mapChefList конвертирует список шефов в слайс chefSnippet
func mapChefList(in []chefEntity.Chef) []chefSnippet {
	out := make([]chefSnippet, 0, len(in))
	for _, c := range in {
		out = append(out, chefSnippet{
			ID:           c.ID,
			Name:         c.Name,
			AvatarURL:    c.SmallImageURL,
			Rating:       c.Rating,
			ReviewsCount: c.ReviewsCount,
		})
	}
	return out
}

// mapDishList конвертирует список блюд в слайс dishDetail без цены и деталей
func (h *usersHandler) mapDishList(ctx context.Context, in []dishesEntity.Dish) []dishDetail {
	out := make([]dishDetail, 0, len(in))
	for _, d := range in {
		minPrice, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, d.ID)
		if err != nil {
			h.log.Error(ctx, "Ошибка получения минимальной цены блюда", err)
			continue
		}
		out = append(out, dishDetail{
			DishID:   d.ID,
			Title:    d.Name,
			ImageURL: d.ImageURL,
			Rating:   d.Rating,
			PriceFrom: Price{
				Value:    minPrice.Value,
				Currency: minPrice.Currency,
			},
		})
	}
	return out
}

func mapOrderProfile(p orderEntity.OrderProfile) *orderDetail {
	if p.Order == nil {
		return nil
	}

	// Формируем метку
	label := fmt.Sprintf("Заказ от %s", p.Order.CreatedAt.Format("02.01.2006, 15:04"))

	// Маппинг позиций
	dishes := make([]dishDetail, 0, len(p.Items))
	for _, it := range p.Items {
		// Преобразуем поле Notes в *string
		var details *string
		if it.Notes != "" {
			details = &it.Notes
		}

		dishes = append(dishes, dishDetail{
			DishID:   it.Dish.ID,
			Title:    it.Dish.Name,
			ImageURL: it.Dish.ImageURL,
			Rating:   it.Dish.Rating,
			Details:  details,
			PriceFrom: Price{
				Value:    it.Size.PriceValue,
				Currency: it.Size.PriceCurrency,
			},
		})
	}

	// Маппинг шефа
	chef := chefSnippet{
		ID:           p.Chef.ID,
		Name:         p.Chef.Name,
		AvatarURL:    p.Chef.SmallImageURL,
		Rating:       p.Chef.Rating,
		ReviewsCount: p.Chef.ReviewsCount,
	}

	// Маппинг отзыва
	var review reviewResponseProfile
	if p.Review != nil {
		review.CanWriteReview = p.Review.CanWriteReview
		if p.Review.Rating != nil {
			review.Rating = *p.Review.Rating
		}
		if p.Review.Comment != nil {
			review.Comment = *p.Review.Comment
		}
	} else {
		// Если ReviewDetail == nil, считаем, что отзыв не писали
		review.CanWriteReview = true
	}

	return &orderDetail{
		OrderID: p.Order.ID,
		Label:   label,
		Dishes:  dishes,
		Chef:    chef,
		Review:  review,
	}
}

func mapOrderProfilesList(in []orderEntity.OrderProfile) []orderDetail {
	out := make([]orderDetail, 0, len(in))
	for _, p := range in {
		if d := mapOrderProfile(p); d != nil {
			out = append(out, *d)
		}
	}
	return out
}

type userProfileData struct {
	UserID          int64         `json:"user_id,omitempty"`
	Email           *string       `json:"email,omitempty"`
	Name            string        `json:"first_name,omitempty"`
	LastName        *string       `json:"last_name,omitempty"`
	PhoneNumber     *string       `json:"phone_number"`
	LastOrder       *orderDetail  `json:"last_order,omitempty"`
	Orders          []orderDetail `json:"orders"`
	FavoriteChefs   any           `json:"favorite_chefs"`
	FavoritesDishes any           `json:"favorites_dishes"`
}

func (u *usersHandler) AddFavoriteChef(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный юзер айди", "message": "Пожалуйста, авторизуйтесь заново."})
		return
	}
	var body struct {
		ChefID int64 `json:"chef_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный формат запроса", "message": "Пожалуйста, проверьте входные данные."})
		return
	}
	if err := u.favoriteUsecase.AddFavoriteChef(ctx, userID, body.ChefID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Не удалось добавить шефа в избранное."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (u *usersHandler) RemoveFavoriteChef(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный юзер айди", "message": "Пожалуйста, авторизуйтесь заново."})
		return
	}
	var body struct {
		ChefID int64 `json:"chef_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный формат запроса", "message": "Пожалуйста, проверьте входные данные."})
		return
	}
	err = u.favoriteUsecase.RemoveFavoriteChef(ctx, userID, body.ChefID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Не удалось удалить шефа из избранного."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (u *usersHandler) AddFavoriteDish(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный юзер айди", "message": "Пожалуйста, авторизуйтесь заново."})
		return
	}
	var body struct {
		DishID int64 `json:"dish_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный формат запроса", "message": "Пожалуйста, проверьте входные данные."})
		return
	}
	if err := u.favoriteUsecase.AddFavoriteDish(ctx, userID, body.DishID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Не удалось добавить блюдо в избранное."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (u *usersHandler) RemoveFavoriteDish(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный юзер айди", "message": "Пожалуйста, авторизуйтесь заново."})
		return
	}
	var body struct {
		DishID int64 `json:"dish_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверный формат запроса", "message": "Пожалуйста, проверьте входные данные."})
		return
	}
	err = u.favoriteUsecase.RemoveFavoriteDish(ctx, userID, body.DishID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Не удалось удалить блюдо из избранного."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type CreateRequest struct {
	Username         string     `json:"username"`                // Уникальный идентификатор
	Alias            string     `json:"alias"`                   // Полное имя (ФИО)
	FirstName        string     `json:"first_name"`              // Имя
	SecondName       *string    `json:"second_name,omitempty"`   // Отчество (опционально)
	LastName         *string    `json:"last_name,omitempty"`     // Фамилия (опционально)
	Email            *string    `json:"email,omitempty"`         // Электронная почта (опционально)
	NumberPhone      *string    `json:"number_phone,omitempty"`  // Телефон (опционально)
	Status           int        `json:"status"`                  // Состояние (0 = не бан, 1 = бан)
	ExternalType     int        `json:"external_type"`           // Внешний статус (0, 1, 2)
	TelegramName     *string    `json:"telegram_name,omitempty"` // Telegram username (опционально)
	ExternalID       *string    `json:"external_id,omitempty"`   // Внешний ID (опционально)
	NotificationFlag int        `json:"notification_flag"`       // Уведомления (0 = включены, 1 = выключены)
	Role             string     `json:"role"`                    // Роль (client, chef, admin)
	Birthday         *time.Time `json:"birthday,omitempty"`      // День рождения (опционально)
}

func (u *usersHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	r := CreateRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.usecase.Create(c.Request.Context(), &usersEntity.User{
		Username:         r.Username,
		Alias:            r.Alias,
		FirstName:        r.FirstName,
		SecondName:       r.SecondName,
		LastName:         r.LastName,
		Email:            r.Email,
		NumberPhone:      r.NumberPhone,
		Status:           r.Status,
		ExternalType:     r.ExternalType,
		TelegramName:     r.TelegramName,
		ExternalID:       r.ExternalID,
		NotificationFlag: r.NotificationFlag,
		Role:             r.Role,
		Birthday:         r.Birthday,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (u *usersHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// TODO: handle error
	}
	user, err := u.usecase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			u.log.Error(ctx, err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch user", "error": err, "user": user})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateRequest struct {
	Alias            string     `json:"alias"`                   // Полное имя (ФИО)
	FirstName        string     `json:"first_name"`              // Имя
	SecondName       *string    `json:"second_name,omitempty"`   // Отчество (опционально)
	LastName         *string    `json:"last_name,omitempty"`     // Фамилия (опционально)
	Email            *string    `json:"email,omitempty"`         // Электронная почта (опционально)
	NumberPhone      *string    `json:"number_phone,omitempty"`  // Телефон (опционально)
	Status           int        `json:"status"`                  // Состояние (0 = не бан, 1 = бан)
	ExternalType     int        `json:"external_type"`           // Внешний статус (0, 1, 2)
	TelegramName     *string    `json:"telegram_name,omitempty"` // Telegram username (опционально)
	ExternalID       *string    `json:"external_id,omitempty"`   // Внешний ID (опционально)
	NotificationFlag int        `json:"notification_flag"`       // Уведомления (0 = включены, 1 = выключены)
	Role             string     `json:"role"`                    // Роль (client, chef, admin)
	Birthday         *time.Time `json:"birthday,omitempty"`      // День рождения (опционально)

}

func (u *usersHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// TODO: handle error
	}

	r := UpdateRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.usecase.Update(ctx, id, usersEntity.User{
		Alias:            r.Alias,
		FirstName:        r.FirstName,
		SecondName:       r.SecondName,
		LastName:         r.LastName,
		Email:            r.Email,
		NumberPhone:      r.NumberPhone,
		Status:           r.Status,
		ExternalType:     r.ExternalType,
		TelegramName:     r.TelegramName,
		ExternalID:       r.ExternalID,
		NotificationFlag: r.NotificationFlag,
		Role:             r.Role,
		Birthday:         r.Birthday,
	})

	if err != nil {
		u.log.Error(ctx, err)
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (u *usersHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// TODO: handle error
	}

	err = u.usecase.Delete(ctx, id)
	if err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *usersHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err:    errorMessage{Code: 4003, Message: "Invalid user ID", Details: "Передан некорректный ID пользователя."},
		})
		return
	}

	var (
		user      *usersEntity.User
		favDishes []dishesEntity.Dish
		favChefs  []chefEntity.Chef
		profiles  []orderEntity.OrderProfile
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		u2, err := h.usecase.GetByID(gctx, id)
		if err != nil {
			return err
		}
		user = u2
		return nil
	})
	g.Go(func() error {
		d, err := h.usecase.GetFavoritesDishesByUserID(gctx, id)
		if err != nil {
			return err
		}
		favDishes = d
		return nil
	})
	g.Go(func() error {
		cfs, err := h.usecase.GetFavoritesChefsByUserID(gctx, id)
		if err != nil {
			return err
		}
		favChefs = cfs
		return nil
	})
	g.Go(func() error {
		profiles, err = h.orderUsecase.GetOrdersByUserID(gctx, id)
		return err
	})
	if err := g.Wait(); err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err:    errorMessage{Code: 5001, Message: "Internal server error", Details: "Не удалось получить данные профиля."},
		})
		return
	}

	chefs := make([]map[string]interface{}, 0)
	for _, topChef := range favChefs {
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
		})
	}

	dishes := make([]map[string]interface{}, 0)
	for _, dish := range favDishes {
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

	// Разбиваем историю на последний заказ и все прочие
	var lastDetail *orderDetail
	var olderDetails []orderDetail
	olderDetails = []orderDetail{}
	if len(profiles) > 0 {
		lastDetail = mapOrderProfile(profiles[0])
		if len(profiles) > 1 {
			olderDetails = mapOrderProfilesList(profiles[1:])
		}
	}

	resp := userProfileData{
		UserID:          user.ID,
		Name:            user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     types.FormatPhoneNumber(user.NumberPhone),
		Email:           user.Email,
		LastOrder:       lastDetail,
		FavoriteChefs:   chefs,
		FavoritesDishes: dishes,
		Orders:          olderDetails,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
