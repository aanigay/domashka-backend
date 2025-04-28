package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	chefEntity "domashka-backend/internal/entity/chefs"
	dishEntity "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/utils/pointers"
)

type dishesHandler struct {
	dishesUsecase dishesUsecase
	chefUsecase   chefUsecase
}

// NewDishesHandler
// TODO:
// (не крит) Разделить инициализацию и регистрацию на отдельные функции,
// то есть сделать отдельные функции для добавления роутов
func NewDishesHandler(rg *gin.RouterGroup, dishesUsecase dishesUsecase, chefUsecase chefUsecase) {
	d := dishesHandler{
		dishesUsecase: dishesUsecase,
		chefUsecase:   chefUsecase,
	}

	rg = rg.Group("/dish")
	{
		rg.GET("/:dishId", d.getDishDetail)
	}
}

type DishesDetailsSuccessResponse struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}
type data struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	ImageURL      string        `json:"image_url"`
	Rating        float32       `json:"rating"`
	Chef          Chef          `json:"chef"`
	Sizes         []size        `json:"sizes"`
	Ingredients   []ingredient  `json:"ingredients"`
	Nutrition     Nutrition     `json:"nutrition"`
	RelatedDishes []dishSnippet `json:"related_dishes"`
}

type Chef struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Rating       float32 `json:"rating"`
	ReviewsCount int32   `json:"reviews_count"`
}

type size struct {
	ID           int64  `json:"id"`
	Label        string `json:"label"`
	Weight       weight `json:"weight"`
	Price        Price  `json:"price"`
	Availability bool   `json:"availability"`
}

type weight struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

type Price struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}
type ingredient struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IsAllergen  bool   `json:"is_allergen"`
	IsRemovable bool   `json:"is_removable"`
}

type Nutrition struct {
	Composition string  `json:"composition"`
	Per100g     per100g `json:"per_100g"`
}

type per100g struct {
	Calories      int `json:"calories"`
	Fat           int `json:"fat"`
	Protein       int `json:"protein"`
	Carbohydrates int `json:"carbohydrates"`
}

type dishSnippet struct {
	ID       int64   `json:"dish_id"`
	Name     string  `json:"name"`
	ImageURL string  `json:"dish_image_url"`
	Price    string  `json:"price"`
	Rating   float32 `json:"rating"`
}

type errorResponse struct {
	Status string       `json:"status"`
	Err    errorMessage `json:"error"`
}
type errorMessage struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (h *dishesHandler) getDishDetail(c *gin.Context) {
	ctx := c.Request.Context()
	g, ctx1 := errgroup.WithContext(ctx)
	dishID, err := strconv.ParseInt(c.Param("dishId"), 10, 64)
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
	// сначала это параллельно
	// Получение dish по ID
	var dish *dishEntity.Dish
	g.Go(func() error {
		dish, err = h.dishesUsecase.GetDishByID(ctx1, dishID)
		if err != nil {
			return err
		}
		return nil
	})
	// Получение ингредиентов
	var ingredients []dishEntity.Ingredient
	g.Go(func() error {
		ingredients, err = h.dishesUsecase.GetIngredientsByDishID(ctx1, dishID)
		if err != nil {
			return err
		}
		return nil
	})
	// Получение nutrition
	var nutrition *dishEntity.Nutrition
	g.Go(func() error {
		nutrition, err = h.dishesUsecase.GetNutritionByDishID(ctx1, dishID)
		if err != nil {
			return err
		}
		return nil
	})
	// получение размеров
	var sizes []dishEntity.Size
	g.Go(func() error {
		sizes, err = h.dishesUsecase.GetDishSizesByDishID(ctx1, dishID)
		if err != nil {
			return err
		}
		return nil
	})
	if err = g.Wait(); err != nil {
		if errors.Is(err, dishEntity.ErrDishNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{
				Status: "error",
				Err: errorMessage{
					Code:    4001,
					Message: "Dish not found",
					Details: "Блюдо с указанным ID не найдено.",
				},
			})
			return
		}
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
	// потом это параллельно
	g, ctx2 := errgroup.WithContext(ctx)
	// получение шефа по dish.ChefID
	var chef *chefEntity.Chef
	g.Go(func() error {
		chef, err = h.chefUsecase.GetChefByDishID(ctx2, dish.ChefID)
		if err != nil {
			return err
		}
		return nil
	})
	var relatedDishes []dishEntity.Dish
	g.Go(func() error {
		relatedDishes, err = h.dishesUsecase.GetDishesByChefID(ctx2, dish.ChefID)
		if err != nil {
			return err
		}
		return nil
	})
	if err = g.Wait(); err != nil {

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

	relatedDishesResponse := make([]dishSnippet, 0, len(relatedDishes))
	for _, relatedDish := range relatedDishes {
		price, err := h.dishesUsecase.GetMinimalPriceByDishID(ctx, relatedDish.ID)
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
		relatedDishesResponse = append(relatedDishesResponse, dishSnippet{
			ID:       relatedDish.ID,
			Name:     relatedDish.Name,
			ImageURL: relatedDish.ImageURL,
			Price:    fmt.Sprintf("от %.0fр", price.Value),
			Rating:   pointers.From(relatedDish.Rating),
		})
	}

	sizesResponse := make([]size, 0, len(sizes))
	for _, val := range sizes {
		sizesResponse = append(sizesResponse, size{
			ID:    val.ID,
			Label: val.Label,
			Weight: weight{
				Value: val.WeightValue,
				Unit:  val.WeightUnit,
			},
			Price: Price{
				Value:    val.PriceValue,
				Currency: val.PriceCurrency,
			},
			Availability: true,
		})
	}
	ingredientsResponse := make([]ingredient, 0, len(ingredients))
	inrStrings := make([]string, 0, len(ingredients))
	for _, val := range ingredients {
		inrStrings = append(inrStrings, strings.ToLower(val.Name))
		ingredientsResponse = append(ingredientsResponse, ingredient{
			ID:          val.ID,
			Name:        val.Name,
			IsAllergen:  val.IsAllergen,
			IsRemovable: val.IsRemovable,
		})
	}

	response := DishesDetailsSuccessResponse{
		Status: "success",
		Data: data{
			ID:          dishID,
			Name:        dish.Name,
			Description: dish.Description,
			ImageURL:    dish.ImageURL,
			Rating:      pointers.From(dish.Rating),
			Chef: Chef{
				ID:           chef.ID,
				Name:         chef.Name,
				Rating:       pointers.From(chef.Rating),
				ReviewsCount: pointers.From(chef.ReviewsCount),
			},
			Sizes:       sizesResponse,
			Ingredients: ingredientsResponse,
			Nutrition: Nutrition{
				Per100g: per100g{
					Calories:      nutrition.Calories,
					Fat:           nutrition.Fat,
					Protein:       nutrition.Protein,
					Carbohydrates: nutrition.Carbohydrates,
				},
				Composition: fmt.Sprintf("Cостав: %s.", strings.Join(inrStrings, ", ")),
			},
			RelatedDishes: relatedDishesResponse,
		},
	}
	c.JSON(http.StatusOK, response)
}
