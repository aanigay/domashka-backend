package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"domashka-backend/internal/custom_errors"
)

type AuthHandler struct {
	authUsecase authUsecase
	jwtUsecase  jwtUsecase
}

func newAuthHandler(rg *gin.RouterGroup, auth authUsecase, jwt jwtUsecase) {
	h := &AuthHandler{authUsecase: auth, jwtUsecase: jwt}

	rg = rg.Group("/auth")
	{
		rg.POST("/register", h.Register)
		rg.POST("/verify", h.Verify)
		rg.POST("/login", h.Login)
		rg.GET("/user", h.ValidateToken)
	}

}

func (h *AuthHandler) Register(c *gin.Context) {

	ctx := c.Request.Context()

	var request struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	err := h.authUsecase.Register(ctx, request.Phone)
	if errors.Is(err, custom_errors.ErrUserExists) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Пользователь с таким номером телефона уже зарегистрирован."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
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

	token, err := h.authUsecase.Verify(ctx, request.Phone, request.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Неверный или просроченный OTP."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success",
		"message": "Номер телефона успешно подтверждён.",
		"token":   token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	err := h.authUsecase.Login(ctx, request.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Пользователь с таким номером телефона еще не зарегистрирован."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Код подтверждения отправлен на указанный номер телефона.",
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
