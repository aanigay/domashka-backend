package v1

import (
	"domashka-backend/internal/custom_errors"
	"domashka-backend/internal/entity/auth"
	"domashka-backend/internal/utils/validation"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authUsecase authUsecase
	jwtUsecase  jwtUsecase
}

func newAuthHandler(rg *gin.RouterGroup, auth authUsecase, jwt jwtUsecase) {
	h := &AuthHandler{authUsecase: auth, jwtUsecase: jwt}

	rg = rg.Group("/auth")
	{
		rg.POST("/login", h.Auth)
		rg.POST("/verify", h.Verify)
		rg.GET("/user", h.ValidateToken)
		rg.POST("/tg", h.TelegramAuth)
	}

}

func (h *AuthHandler) Auth(c *gin.Context) {
	request := auth.Request{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	err := h.authUsecase.Auth(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Код подтверждения отправлен на указанный номер телефона.",
	})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		Phone string `json:"phone" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}
	// todo: Убрать костыль перед выгрузкой в прод
	if request.OTP == "0123" {
		userID, chefID, token, _ := h.authUsecase.Verify(ctx, request.Phone, request.OTP, "admin")
		c.JSON(http.StatusOK, gin.H{"status": "success",
			"message": "Номер телефона успешно подтверждён.",
			"token":   token,
			"user_id": userID,
			"chef_id": chefID,
		})
		return
	}

	if len(request.OTP) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Неверный формат OTP."})
		return
	}

	userID, chefID, token, err := h.authUsecase.Verify(ctx, request.Phone, request.OTP, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Неверный или просроченный OTP."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success",
		"message": "Номер телефона успешно подтверждён.",
		"token":   token,
		"user_id": userID,
		"chef_id": chefID,
	})
}

func (h *AuthHandler) ValidateToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Токен отсутствует."})
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtUsecase.ValidateJWT(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Токен недействителен."})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *AuthHandler) TelegramAuth(ctx *gin.Context) {
	var request struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if !validation.ValidatePhoneNumber(request.Phone) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid phone number"})
		return
	}

	err := h.authUsecase.AuthViaTg(ctx.Request.Context(), request.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "pending", "message": "Ожидается подтверждение через Telegram."})
}

func (h *AuthHandler) TelegramAuthStatus(ctx *gin.Context) {
	phone := ctx.Query("phone")
	if phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	token, err := h.authUsecase.AuthViaTgStatus(ctx, phone)
	if err != nil {
		if errors.Is(err, custom_errors.ErrConfirmationNotReceived) {
			ctx.JSON(http.StatusOK, gin.H{"status": "pending", "message": "Ожидается подтверждение через Telegram."})
			return
		}
		if errors.Is(err, custom_errors.ErrExpiredTTL) {
			ctx.JSON(http.StatusOK, gin.H{"status": "error", "message": "Подтверждение через Telegram не получено."})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}
