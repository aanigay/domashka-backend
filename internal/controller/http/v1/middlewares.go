package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt jwtUsecase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Токен отсутствует."})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Некорректный формат токена."})
			ctx.Abort()
			return
		}

		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Токен недействителен или истек."})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Set("user_id", claims["user_id"])
		ctx.Set("role", claims["role"])

		ctx.Next()
	}
}
