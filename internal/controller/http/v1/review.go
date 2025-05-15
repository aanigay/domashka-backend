// package v1

package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	reviewEntity "domashka-backend/internal/entity/reviews"
)

// reviewsHandler отвечает за эндпоинты работы с отзывами.
type reviewsHandler struct {
	reviewsUsecase reviewsUsecase
}

// RegisterReviewHandlers регистрирует маршруты для отзывов.
func RegisterReviewHandlers(rg *gin.RouterGroup, reviewsUsecase reviewsUsecase) {
	h := &reviewsHandler{reviewsUsecase: reviewsUsecase}
	grp := rg.Group("/reviews")
	{
		grp.POST("", h.createReview)
		grp.GET("/chef/:chefId", h.getReviewsByChefID)
	}
}

// request/response DTOs

type createReviewRequest struct {
	ChefID  int64  `json:"chef_id"  binding:"required"`
	UserID  int64  `json:"user_id"  binding:"required"`
	Stars   int16  `json:"stars"    binding:"required,min=1,max=5"`
	Comment string `json:"comment"  binding:"max=255"`
	OrderID int64  `json:"order_id" binding:"required"`
}

type reviewResponse struct {
	ID              int64     `json:"id"`
	ChefID          int64     `json:"chef_id"`
	UserID          int64     `json:"user_id"`
	Stars           int16     `json:"stars"`
	Comment         string    `json:"comment,omitempty"`
	IsVerified      bool      `json:"is_verified"`
	IncludeInRating bool      `json:"include_in_rating"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	OrderID         int64     `json:"order_id"`
}

type createReviewSuccessResponse struct {
	Status string `json:"status"` // всегда "success"
}

type reviewsListSuccessResponse struct {
	Status string           `json:"status"` // всегда "success"
	Data   []reviewResponse `json:"data"`
}

// POST /v1/reviews
func (h *reviewsHandler) createReview(c *gin.Context) {
	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err:    errorMessage{Code: 4001, Message: "Invalid request", Details: err.Error()},
		})
		return
	}

	rv := reviewEntity.Review{
		ChefID:  req.ChefID,
		UserID:  req.UserID,
		Stars:   req.Stars,
		Comment: req.Comment,
		OrderID: req.OrderID,
	}

	if err := h.reviewsUsecase.CreateReview(c.Request.Context(), rv); err != nil {
		// можно различать типы ошибок, например ErrReviewNotFound, но тут при создании
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err:    errorMessage{Code: 5001, Message: "Failed to create review", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, createReviewSuccessResponse{Status: "success"})
}

// GET /v1/reviews/chef/:chefId
func (h *reviewsHandler) getReviewsByChefID(c *gin.Context) {
	chefID, err := strconv.ParseInt(c.Param("chefId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status: "error",
			Err:    errorMessage{Code: 4002, Message: "Invalid chef ID"},
		})
		return
	}

	list, err := h.reviewsUsecase.GetReviewsByChefID(c.Request.Context(), chefID)
	if err != nil {
		// если usecase может возвращать ErrReviewNotFound, можно проверить errors.Is(...)
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status: "error",
			Err:    errorMessage{Code: 5002, Message: "Failed to fetch reviews", Details: err.Error()},
		})
		return
	}

	resp := make([]reviewResponse, 0, len(list))
	for _, rv := range list {
		// пропускаем мягко удалённые, если usecase не делает этого сам
		if rv.IsDeleted {
			continue
		}
		resp = append(resp, reviewResponse{
			ID:              rv.ID,
			ChefID:          rv.ChefID,
			UserID:          rv.UserID,
			Stars:           rv.Stars,
			Comment:         rv.Comment,
			IsVerified:      rv.IsVerified,
			IncludeInRating: rv.IncludeInRating,
			CreatedAt:       rv.CreatedAt,
			UpdatedAt:       rv.UpdatedAt,
			OrderID:         rv.OrderID,
		})
	}

	c.JSON(http.StatusOK, reviewsListSuccessResponse{
		Status: "success",
		Data:   resp,
	})
}
