// internal/app/app.go
package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "domashka-backend/docs"

	"domashka-backend/config"
	v1 "domashka-backend/internal/controller/http/v1"
	cartpgrepo "domashka-backend/internal/repositories/cart"
	geopgrepo "domashka-backend/internal/repositories/geo"
	notifpgrepo "domashka-backend/internal/repositories/notifications"
	userspgrepo "domashka-backend/internal/repositories/users"
	authusecase "domashka-backend/internal/usecase/auth"
	cartusecase "domashka-backend/internal/usecase/cart"
	geousecase "domashka-backend/internal/usecase/geo"
	jwtusecase "domashka-backend/internal/usecase/jwt"
	notifusecase "domashka-backend/internal/usecase/notifications"
	usersusecase "domashka-backend/internal/usecase/users"
	"domashka-backend/pkg/logger"
	smtpmail "domashka-backend/pkg/mail"
	"domashka-backend/pkg/postgres"
	"domashka-backend/pkg/redis"
)

type Application struct {
	Config *config.Config
	DB     *sql.DB
	Redis  *redis.Redis
}

func Run(cfg *config.Config) {
	l := logger.New()
	pg, err := postgres.New(cfg.DB.GetDSN(), postgres.MaxPoolSize(cfg.DB.PoolCapacity))
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer pg.Close()

	redisClient, err := redis.New(cfg.Redis)
	if err != nil || redisClient.Ping() != nil {
		log.Fatalf("Ошибка инициализации Redis: %v", err)
	}

	smtpClient := smtpmail.New(cfg.SMTP)

	// Repositories
	notifPGRepo := notifpgrepo.New(pg)
	usersPGRepo := userspgrepo.New(pg)
	geoPGRepo := geopgrepo.New(pg)
	cartPGRepo := cartpgrepo.New(pg)

	// Use Cases (сервсисы)
	userUseCase := usersusecase.New(usersPGRepo, cartPGRepo)
	jwtUseCase := jwtusecase.New(cfg.JWT)
	authUseCase := authusecase.New(usersPGRepo, redisClient, jwtUseCase)
	geoUseCase := geousecase.New(geoPGRepo)
	notifUseCase := notifusecase.New(notifPGRepo, smtpClient)
	cartUseCase := cartusecase.New(cartPGRepo)

	// Http Server
	handler := gin.New()
	v1.NewRouter(
		handler,
		l,
		userUseCase,
		authUseCase,
		jwtUseCase,
		notifUseCase,
		geoUseCase,
		cartUseCase,
	)

	err = handler.Run(fmt.Sprintf(":%s", cfg.HostConfig.Port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
