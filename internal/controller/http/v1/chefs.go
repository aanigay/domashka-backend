package v1

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	chefEntity "domashka-backend/internal/entity/chefs"
	dishEntity "domashka-backend/internal/entity/dishes"
	geoEntity "domashka-backend/internal/entity/geo"
	"domashka-backend/internal/utils/pointers"
)

func getExperienceString(years int) string {
	if years == 0 {
		return "Не указан"
	}
	mod100 := years % 100
	mod10 := years % 10

	var suffix string
	if mod100 >= 11 && mod100 <= 14 {
		suffix = "лет"
	} else {
		switch mod10 {
		case 1:
			suffix = "год"
		case 2, 3, 4:
			suffix = "года"
		default:
			suffix = "лет"
		}
	}
	return fmt.Sprintf("%d %s", years, suffix)
}

type chefsHendler struct {
	dishesUsecase  dishesUsecase
	chefUsecase    chefUsecase
	geoUsecase     geoUsecase
	shiftsUsecase  shiftsUsecase
	usersUsecase   usersUsecase
	reviewsUsecase reviewsUsecase
	orderUsecase   orderUsecase
}

func NewChefsHandler(
	rg *gin.RouterGroup,
	dishesUsecase dishesUsecase,
	chefUsecase chefUsecase,
	geoUsecase geoUsecase,
	shiftsUsecase shiftsUsecase,
	usersUsecase usersUsecase,
	reviewsUsecase reviewsUsecase,
	orderUsecase orderUsecase,
) {
	ch := chefsHendler{
		dishesUsecase:  dishesUsecase,
		chefUsecase:    chefUsecase,
		geoUsecase:     geoUsecase,
		shiftsUsecase:  shiftsUsecase,
		usersUsecase:   usersUsecase,
		reviewsUsecase: reviewsUsecase,
		orderUsecase:   orderUsecase,
	}

	rg = rg.Group("/chefs")
	{
		rg.GET("/:chefId", ch.getChefDetail)
		rg.POST("/avatar/:chefId", ch.UploadChefAvatar)
		rg.POST("/shifts/open", ch.openShift)
		rg.POST("/shifts/close", ch.closeShift)
		rg.POST("/dishes/create", ch.CreateDish)
		rg.DELETE("/dishes/delete", ch.DeleteDish)
		rg.GET("/menu", ch.getMenu)
		rg.GET("/dish/form", ch.createDishForm)
		rg.GET("/stats", ch.getStats)
	}
}

func (h *chefsHendler) getChefDetail(c *gin.Context) {
	// 1. Контекст запроса
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

	// 2. Парсим chefId
	chefID, err := strconv.ParseInt(c.Param("chefId"), 10, 64)
	if err != nil {
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

	// 3. Параллельно получаем шефа, опыт, сертификаты и блюда
	var (
		chef    *chefEntity.Chef
		years   int
		certs   []chefEntity.Certification
		dishes  []dishEntity.Dish
		address *geoEntity.Address
	)
	g, ctx1 := errgroup.WithContext(ctx)
	g.Go(func() error {
		var e error
		chef, e = h.chefUsecase.GetChefByID(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		years, e = h.chefUsecase.GetChefExperienceYears(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		certs, e = h.chefUsecase.GetChefCertifications(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		dishes, e = h.dishesUsecase.GetAllDishesByChefID(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		address, e = h.geoUsecase.GetLastUpdatedClientAddress(ctx1, userID)
		return e
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}

	if address == nil {
		address = &geoEntity.Address{
			Longitude: defaultLong,
			Latitude:  defaultLat,
		}
	}
	distance, err := h.chefUsecase.GetDistanceToChef(ctx, address.Latitude, address.Longitude, chefID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4003,
				Message: "Internal server error",
				Details: fmt.Sprintf("%v", err),
			},
		})
		return
	}

	// 4. Формируем expStr и список имён сертификатов
	expStr := getExperienceString(years)

	certNames := make([]string, len(certs))
	for i, cert := range certs {
		certNames[i] = cert.Name
	}

	// 5. Параллельно для каждого блюда запрашиваем минимальную цену
	dishItems := make([]gin.H, len(dishes))
	g2, ctx2 := errgroup.WithContext(ctx)
	for i, d := range dishes {
		i, d := i, d // замыкание
		g2.Go(func() error {
			price, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx2, d.ID)
			if err != nil {
				return err
			}
			category, err := h.dishesUsecase.GetCategoryTitleByDishID(ctx2, d.ID)
			if err != nil {
				return err
			}
			dishItems[i] = gin.H{
				"id":            d.ID,
				"name":          d.Name,
				"description":   d.Description,
				"image_url":     d.ImageURL,
				"rating":        d.Rating,
				"reviews_count": d.ReviewsCount,
				"price_from": gin.H{
					"value":    price.Value,
					"currency": price.Currency,
				},
				"category": category,
			}
			return nil
		})
	}
	if err := g2.Wait(); err != nil {
		log.Printf("GetMinimalPriceByDishID error: %v", err)
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4005,
				Message: "Cannot load price",
				Details: "Не удалось получить цену для блюда.",
			},
		})
		return
	}

	// Макс кол-во отзывов на странице
	limit := 5
	rawReviews, err := h.reviewsUsecase.GetFullReviewsByChefID(ctx, chefID, limit)
	reviewsList := make([]gin.H, 0)
	if err != nil {
		log.Printf("failed to load reviews for chef %d: %v", chefID, err)
	} else {
		for _, rv := range rawReviews {
			reviewsList = append(reviewsList, gin.H{
				"user":    rv.UserName,
				"comment": rv.Comment,
			})
		}
	}
	menu := map[string][]interface{}{}
	for _, dishItem := range dishItems {
		categoryTitle := dishItem["category"].(string)
		menu[categoryTitle] = append(menu[categoryTitle], dishItem)
	}
	var menuResponse []map[string]interface{}
	var catList []string
	for categoryTitle, dishes := range menu {
		catList = append(catList, categoryTitle)
		menuResponse = append(menuResponse, gin.H{
			"name":  categoryTitle,
			"items": dishes,
		})
	}
	isFavorite := false
	favs, err := h.usersUsecase.GetFavoritesChefsByUserID(ctx, userID)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if favs != nil {
			for _, fav := range favs {
				if fav.ID == chefID {
					isFavorite = true
					break
				}
			}
		}
	}
	// 7. Выдаем финальный JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":               chef.ID,
			"name":             chef.Name,
			"avatar_url":       chef.ImageURL,
			"distance_km":      fmt.Sprintf("%.1f км", distance/1000),
			"certifications":   certNames,
			"experience_years": expStr,
			"description":      chef.Description,
			"rating":           chef.Rating,
			"reviews_count":    chef.ReviewsCount,
			"reviews":          reviewsList,
			"menu": gin.H{
				"categories_lists": catList,
				"categories":       menuResponse,
			},
			"legal_info":  chef.LegalInfo,
			"is_favorite": isFavorite,
		},
	})
}

func (h *chefsHendler) UploadChefAvatar(c *gin.Context) {
	// --- Парсинг multipart/form-data ---
	idStr := c.Param("chefId")
	chefID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Неверный юзер айди",
			"message": "Пожалуйста, авторизуйтесь заново.",
		})
		return
	}
	imageSizeLabel := c.Query("size")
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  gin.H{"code": 4001, "message": "Avatar file is required"},
		})
		return
	}

	// --- Проверка размера (≤5 МБ) ---
	const maxAvatarSize = 5 << 20 // 5 MiB
	if fileHeader.Size > maxAvatarSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"status": "error",
			"error":  gin.H{"code": 4002, "message": "File size exceeds 5MB"},
		})
		return
	}

	// --- Проверка формата файла ---
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  gin.H{"code": 4001, "message": "Unsupported file format"},
		})
		return
	}

	var avatarURL string

	// --- Вызов бизнес-логики ---
	if imageSizeLabel == "m" {
		avatarURL, err = h.chefUsecase.UploadAvatar(c.Request.Context(), chefID, fileHeader)
	} else {
		avatarURL, err = h.chefUsecase.UploadSmallAvatar(c.Request.Context(), chefID, fileHeader)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  gin.H{"code": 5001, "message": "Could not save avatar"},
		})
		return
	}

	// --- Успешный ответ ---
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"chef_id":     chefID,
			"avatar_url":  avatarURL,
			"uploaded_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *chefsHendler) openShift(c *gin.Context) {
	ctx := c.Request.Context()
	chefID, err := strconv.ParseInt(c.Query("chef_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid chef ID",
				Details: "Передан некорректный ID повара.",
			},
		})
		return
	}
	err = h.shiftsUsecase.OpenShift(ctx, chefID)
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

func (h *chefsHendler) closeShift(c *gin.Context) {
	ctx := c.Request.Context()
	chefID, err := strconv.ParseInt(c.Query("chef_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid chef ID",
				Details: "Передан некорректный ID повара.",
			},
		})
		return
	}
	err = h.shiftsUsecase.CloseShift(ctx, chefID)
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

type CreateDishRequest struct {
	ChefID int64       `json:"chef_id"`
	Dish   DishDetails `json:"dish"`
}

type DishDetails struct {
	ID            *int64  `json:"id,omitempty"`
	Title         string  `json:"title"`
	ImageURL      *string `json:"image_url,omitempty"`
	CategoryID    int64   `json:"category_id"`
	CategoryTitle *string `json:"category_title,omitempty"`
	Description   string  `json:"description"`
	Sizes         []struct {
		Label  string `json:"label"`
		Weight weight `json:"weight"`
		Price  Price  `json:"price"`
	} `json:"sizes"`
	Ingredients []struct {
		ID          int64   `json:"id"`
		Title       *string `json:"title,omitempty"`
		IsRemovable bool    `json:"is_removable"`
	} `json:"ingredients"`
	Nutrition *per100g `json:"nutrition,omitempty"`
}

func (h *chefsHendler) CreateDish(context *gin.Context) {
	ctx := context.Request.Context()
	var req CreateDishRequest
	if err := context.BindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4001,
				Message: err.Error(),
				Details: err.Error()}},
		)
		return
	}
	sizes := make([]dishEntity.Size, 0, len(req.Dish.Sizes))
	for _, size := range req.Dish.Sizes {
		sizes = append(sizes, dishEntity.Size{
			Label:         size.Label,
			WeightUnit:    size.Weight.Unit,
			WeightValue:   size.Weight.Value,
			PriceValue:    size.Price.Value,
			PriceCurrency: size.Price.Currency,
		})
	}
	var nutrition *dishEntity.Nutrition
	if req.Dish.Nutrition != nil {
		nutrition = &dishEntity.Nutrition{
			Calories:      req.Dish.Nutrition.Calories,
			Protein:       req.Dish.Nutrition.Protein,
			Fat:           req.Dish.Nutrition.Fat,
			Carbohydrates: req.Dish.Nutrition.Carbohydrates,
		}
	}
	ingredients := make([]dishEntity.Ingredient, 0, len(req.Dish.Ingredients))
	for _, ingredient := range req.Dish.Ingredients {
		ingredients = append(ingredients, dishEntity.Ingredient{
			ID:          ingredient.ID,
			IsRemovable: ingredient.IsRemovable,
		})
	}
	var dishID int64
	var err error
	if req.Dish.ID != nil {
		dishID, err = h.dishesUsecase.Update(
			ctx,
			&dishEntity.Dish{
				ID:          *req.Dish.ID,
				ChefID:      req.ChefID,
				Name:        req.Dish.Title,
				CategoryID:  req.Dish.CategoryID,
				Description: req.Dish.Description,
			},
			nutrition,
			sizes,
			ingredients)
	} else {
		dishID, err = h.dishesUsecase.Create(
			ctx,
			&dishEntity.Dish{
				ChefID:      req.ChefID,
				Name:        req.Dish.Title,
				CategoryID:  req.Dish.CategoryID,
				Description: req.Dish.Description,
			},
			nutrition,
			sizes,
			ingredients,
		)
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"dish_id": dishID,
		},
	})
}

func (h *chefsHendler) getMenu(context *gin.Context) {
	ctx := context.Request.Context()
	chefID, err := strconv.ParseInt(context.Query("chef_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid chef ID",
				Details: "Передан некорректный ID повара.",
			},
		})
		return
	}
	dishes, err := h.dishesUsecase.GetDishesByChefID(ctx, chefID, 10000)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	dishesResp := make([]map[string]interface{}, 0, len(dishes))
	for _, dish := range dishes {
		if dish.IsDeleted {
			continue
		}
		dishesResp = append(dishesResp, map[string]interface{}{
			"title": dish.Name,
			"id":    dish.ID,
			"image": dish.ImageURL,
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"dishes": dishesResp,
		},
	})
}

func (h *chefsHendler) createDishForm(context *gin.Context) {
	ctx := context.Request.Context()
	ingredients, err := h.dishesUsecase.GetAllIngredients(ctx)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	categories, err := h.dishesUsecase.GetAllCategories(ctx)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}

	if context.Query("dish_id") == "" {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"ingredients": ingredients,
				"categories":  categories,
			},
		})
		return
	}
	dishID, err := strconv.ParseInt(context.Query("dish_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid dish ID",
				Details: "Передан некорректный ID блюда.",
			},
		})
		return
	}
	dish, err := h.dishesUsecase.GetDishByID(ctx, dishID)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
	}
	if dish.IsDeleted {
		context.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Dish is deleted",
				Details: "Данное блюдо удалено.",
			},
		})
		return
	}
	catTitle, err := h.dishesUsecase.GetCategoryTitleByDishID(ctx, dishID)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	sizes, err := h.dishesUsecase.GetDishSizesByDishID(ctx, dishID)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",

			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}

	dishResp := DishDetails{
		Title:         dish.Name,
		CategoryID:    dish.CategoryID,
		ImageURL:      pointers.To(dish.ImageURL),
		CategoryTitle: &catTitle,
		Description:   dish.Description,
		Sizes:         nil,
		Ingredients:   nil,
		Nutrition:     nil,
	}

	for _, size := range sizes {
		dishResp.Sizes = append(dishResp.Sizes, struct {
			Label  string `json:"label"`
			Weight weight `json:"weight"`
			Price  Price  `json:"price"`
		}{
			Label: size.Label,
			Weight: weight{
				Unit:  size.WeightUnit,
				Value: size.WeightValue,
			},
			Price: Price{
				Currency: size.PriceCurrency,
				Value:    size.PriceValue,
			},
		})
	}

	dishIngredients, err := h.dishesUsecase.GetIngredientsByDishID(ctx, dishID)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	for _, dishIngredient := range dishIngredients {
		dishResp.Ingredients = append(dishResp.Ingredients, struct {
			ID          int64   `json:"id"`
			Title       *string `json:"title,omitempty"`
			IsRemovable bool    `json:"is_removable"`
		}{
			ID:          dishIngredient.ID,
			Title:       &dishIngredient.Name,
			IsRemovable: dishIngredient.IsRemovable,
		})
	}

	dishNutrition, err := h.dishesUsecase.GetNutritionByDishID(ctx, dishID)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Details: err.Error(),
			},
		})
		return
	}
	if dishNutrition != nil {
		dishResp.Nutrition = &per100g{
			Protein:       dishNutrition.Protein,
			Fat:           dishNutrition.Fat,
			Carbohydrates: dishNutrition.Carbohydrates,
			Calories:      dishNutrition.Calories,
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"ingredients": ingredients,
			"categories":  categories,
			"dish":        dishResp,
		},
	})
	return
}

func (h *chefsHendler) DeleteDish(c *gin.Context) {
	ctx := c.Request.Context()
	dishID, err := strconv.ParseInt(c.Query("dish_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid dish ID",
				Details: "Передан некорректный ID блюда.",
			},
		})
		return
	}
	err = h.dishesUsecase.Delete(ctx, dishID)
	if err != nil {
		log.Println(err)
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

type DailyIncome struct {
	Date   string  `json:"date"`
	Profit float32 `json:"profit"`
}
type Monthly struct {
	Month       string        `json:"month"`
	DailyIncome []DailyIncome `json:"daily"`
}
type Income struct {
	Monthly []Monthly `json:"monthly"`
}

type ReviewCount struct {
	Rating int16 `json:"rating"`
	Count  int32 `json:"count"`
}

type Review struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Rating int16  `json:"rating"`
}
type Rating struct {
	ChefRating    *float32      `json:"chef_rating"`
	ReviewsCounts []ReviewCount `json:"reviews_counts"`
	Reviews       []Review      `json:"reviews"`
}

type TopDish struct {
	DishID      int64    `json:"dish_id"`
	ImageURL    string   `json:"dish_image_url"`
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	Rating      *float32 `json:"rating"`
	OrdersCount int32    `json:"orders_count"`
}
type chefStats struct {
	Income    Income    `json:"income"`
	Rating    Rating    `json:"rating"`
	TopDishes []TopDish `json:"top_dishes"`
}

func (h *chefsHendler) getStats(c *gin.Context) {
	ctx := c.Request.Context()
	chefID, err := strconv.ParseInt(c.Query("chef_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err: errorMessage{
				Code:    4002,
				Message: "Invalid chef ID",
				Details: err.Error(),
			},
		})
		return
	}
	profits, err := h.shiftsUsecase.GetDailyProfits(ctx, chefID)
	if err != nil {
		log.Println(err)
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
	profitsMap := map[string]map[string]float32{}
	for _, profit := range profits {
		if _, ok := profitsMap[profit.Month]; !ok {
			profitsMap[profit.Month] = map[string]float32{}
		}
		profitsMap[profit.Month][profit.Date] = profit.Profit
	}
	monthlyIncomes := make([]Monthly, 0)
	for month, profs := range profitsMap {
		dailyIncomes := make([]DailyIncome, 0)
		for date, prof := range profs {
			dailyIncomes = append(dailyIncomes, DailyIncome{
				Date:   date,
				Profit: prof,
			})
		}
		monthlyIncomes = append(monthlyIncomes, Monthly{
			Month:       month,
			DailyIncome: dailyIncomes,
		})
	}
	rating := Rating{
		Reviews: make([]Review, 0),
	}
	ratingStats, err := h.reviewsUsecase.GetRatingStats(ctx, chefID)
	if err != nil {
		log.Println(err)
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
	for stars, r := range ratingStats {
		rating.ReviewsCounts = append(rating.ReviewsCounts, ReviewCount{
			Rating: stars,
			Count:  r,
		})
	}
	reviews, err := h.reviewsUsecase.GetFullReviewsByChefID(ctx, chefID, 1000)
	if err != nil {
		log.Println(err)
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
	reviewsResp := make([]Review, 0)
	for _, review := range reviews {
		reviewsResp = append(reviewsResp, Review{
			Name:   review.UserName,
			Text:   review.Comment,
			Rating: review.Stars,
		})
	}
	chef, err := h.chefUsecase.GetChefByID(ctx, chefID)
	if err != nil {
		log.Println(err)
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
	rating.ChefRating = chef.Rating
	rating.Reviews = reviewsResp
	topDishes := make([]TopDish, 0)
	chefDishes, err := h.dishesUsecase.GetDishesByChefID(ctx, chefID, 1000)
	if err != nil {
		log.Println(err)
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
	for _, dish := range chefDishes {
		if dish.IsDeleted {
			continue
		}
		minPrice, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, dish.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cnt, err := h.orderUsecase.CountDishInOrders(ctx, dish.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		topDishes = append(topDishes, TopDish{
			DishID:      dish.ID,
			ImageURL:    dish.ImageURL,
			Name:        dish.Name,
			Price:       fmt.Sprintf("от %.0fр", minPrice.Value),
			Rating:      dish.Rating,
			OrdersCount: int32(cnt),
		})
	}
	stats := chefStats{
		Income:    Income{monthlyIncomes},
		Rating:    rating,
		TopDishes: topDishes,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}
