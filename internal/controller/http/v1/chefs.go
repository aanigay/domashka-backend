package v1

import (
	chefEntity "domashka-backend/internal/entity/chefs"
	"log"

	dishEntity "domashka-backend/internal/entity/dishes"

	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"net/http"
)

func getExperienceString(years int) string {
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
	dishesUsecase dishesUsecase
	chefUsecase   chefUsecase
}

func NewChefsHandler(rg *gin.RouterGroup, dishesUsecase dishesUsecase, chefUsecase chefUsecase) {
	ch := chefsHendler{
		dishesUsecase: dishesUsecase,
		chefUsecase:   chefUsecase,
	}

	rg = rg.Group("/chefs")
	{
		rg.GET("/:chefId", ch.getChefDetail)
		rg.POST("/avatar/:chefId", ch.UploadChefAvatar)
	}
}

func (ch *chefsHendler) getChefDetail(c *gin.Context) {
	// 1. Контекст запроса
	ctx := c.Request.Context()

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
		chef   *chefEntity.Chef
		years  int
		certs  []chefEntity.Certification
		dishes []dishEntity.Dish
	)
	g, ctx1 := errgroup.WithContext(ctx)
	g.Go(func() error {
		var e error
		chef, e = ch.chefUsecase.GetChefByID(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		years, e = ch.chefUsecase.GetChefExperienceYears(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		certs, e = ch.chefUsecase.GetChefCertifications(ctx1, chefID)
		return e
	})
	g.Go(func() error {
		var e error
		dishes, e = ch.dishesUsecase.GetDishesByChefID(ctx1, chefID)
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
			price, err := ch.dishesUsecase.GetMinimalPriceByDishID(ctx2, d.ID)
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

	// 6. Моковый отзыв
	reviewComment := fmt.Sprintf(
		"Каждый вечер заказываем еду у %s. Ощущение, будто моя бабушка готовит нам.",
		chef.Name,
	)

	// 7. Выдаем финальный JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":               chef.ID,
			"name":             chef.Name,
			"avatar_url":       chef.ImageURL,
			"distance_km":      "4.7",
			"certifications":   certNames,
			"experience_years": expStr,
			"description":      chef.Description,
			"rating":           chef.Rating,
			"reviews_count":    chef.ReviewsCount,
			"reviews": []gin.H{
				{
					"user":    "Ульяна",
					"comment": reviewComment,
				},
			},
			"menu": gin.H{
				"categories_lists": []string{"Все блюда"},
				"categories": []gin.H{
					{
						"name":   "Все блюда",
						"dishes": dishItems,
					},
				},
			},
			"legal_info": chef.LegalInfo,
		},
	})
}

func (ch *chefsHendler) UploadChefAvatar(c *gin.Context) {
	// --- Парсинг multipart/form-data ---
	idStr := c.Param("chefId")
	chefID, err := strconv.ParseInt(idStr, 10, 64)
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

	// --- Вызов бизнес-логики ---
	avatarURL, err := ch.chefUsecase.UploadAvatar(c.Request.Context(), chefID, fileHeader)
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
			"chefId":     chefID,
			"avatarUrl":  avatarURL,
			"uploadedAt": time.Now().UTC().Format(time.RFC3339),
		},
	})
}
