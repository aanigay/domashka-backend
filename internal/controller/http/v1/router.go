package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "domashka-backend/docs"
)

// @title           Domashka
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9093
// @BasePath  /v1
func NewRouter(handler *gin.Engine, l logger, u usersUsecase, a authUsecase, jwt jwtUsecase, n notificationUsecase, g geoUsecase, cartUsecase cartUsecase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/v1")
	{
		newUsersHandler(h, l, u)
		newAuthHandler(h, a, jwt)
		authorized := h.Group("/")
		authorized.Use(AuthMiddleware(jwt))
		{
			newNotificationHandler(authorized, n)
			newGeoHandler(h, g)
		}
		registerCartHandler(h, l, cartUsecase)
	}
	h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
