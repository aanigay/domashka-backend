// internal/app/app.go
package app

import (
	"database/sql"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v4"

	"domashka-backend/config"

	"domashka-backend/pkg/logger"
	smtpmail "domashka-backend/pkg/mail"
	"domashka-backend/pkg/postgres"
	"domashka-backend/pkg/redis"
	"domashka-backend/pkg/sms"

	"domashka-backend/internal/clients/s3"
	v1 "domashka-backend/internal/controller/http/v1"
	"domashka-backend/internal/controller/telegram"
	cartrepo "domashka-backend/internal/repositories/cart"
	chefsrepo "domashka-backend/internal/repositories/chefs"
	dishesrepo "domashka-backend/internal/repositories/dishes"
	favoritesrepo "domashka-backend/internal/repositories/favorites"
	geopgrepo "domashka-backend/internal/repositories/geo"
	notifpgrepo "domashka-backend/internal/repositories/notifications"
	ordersrepo "domashka-backend/internal/repositories/orders"
	reviewsrepo "domashka-backend/internal/repositories/reviews"
	shiftsrepo "domashka-backend/internal/repositories/shifts"
	userspgrepo "domashka-backend/internal/repositories/users"
	authusecase "domashka-backend/internal/usecase/auth"
	cartusecase "domashka-backend/internal/usecase/cart"
	chefsusecase "domashka-backend/internal/usecase/chefs"
	dishesusecase "domashka-backend/internal/usecase/dishes"
	favoritesusecase "domashka-backend/internal/usecase/favorites"
	geousecase "domashka-backend/internal/usecase/geo"
	jwtusecase "domashka-backend/internal/usecase/jwt"
	notifusecase "domashka-backend/internal/usecase/notifications"
	ordersusecase "domashka-backend/internal/usecase/order"
	reviewsusecase "domashka-backend/internal/usecase/reviews"
	shiftsusecase "domashka-backend/internal/usecase/shifts"
	"domashka-backend/internal/usecase/tg"
	usersusecase "domashka-backend/internal/usecase/users"
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

	s3client, err := s3.New(cfg.S3)
	if err != nil {
		log.Fatal(err)
	}
	smtpClient := smtpmail.New(cfg.SMTP)
	smsClient := sms.New()

	dishReviewsWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.URL),
		Topic:        "dish_reviews",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	}
	defer func() {
		err := dishReviewsWriter.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	chefReviewsWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.URL),
		Topic:        "chef_reviews",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	}
	defer func() {
		err := chefReviewsWriter.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	// Repositories
	notifPGRepo := notifpgrepo.New(pg)
	usersPGRepo := userspgrepo.New(pg)
	geoPGRepo := geopgrepo.New(pg)
	dishesPGRepo := dishesrepo.New(pg)
	chefsPGRepo := chefsrepo.New(pg)
	cartPGRepo := cartrepo.New(pg)
	ordersPGRepo := ordersrepo.New(pg)
	shiftsPGRepo := shiftsrepo.New(pg)
	reviewsPGRepo := reviewsrepo.New(pg)
	favoritesPGRepo := favoritesrepo.New(pg)

	// Use Cases (сервисы)
	userUseCase := usersusecase.New(usersPGRepo)
	dishesUsecase := dishesusecase.New(dishesPGRepo, s3client)
	chefsUsecase := chefsusecase.New(chefsPGRepo, geoPGRepo, s3client)
	jwtUseCase := jwtusecase.New(cfg.JWT)
	authUseCase := authusecase.New(usersPGRepo, redisClient, jwtUseCase, smsClient)
	geoUseCase := geousecase.New(geoPGRepo)
	notifUseCase := notifusecase.New(notifPGRepo, smtpClient)
	cartUsecase := cartusecase.New(cartPGRepo)
	shiftsUsecase := shiftsusecase.New(shiftsPGRepo)
	reviewsUsecase := reviewsusecase.New(reviewsPGRepo, usersPGRepo, ordersPGRepo, dishReviewsWriter, chefReviewsWriter)
	ordersUsecase := ordersusecase.New(geoUseCase, cartUsecase, shiftsPGRepo, ordersPGRepo, dishesUsecase, chefsUsecase, reviewsUsecase)
	favoritesUsecase := favoritesusecase.New(favoritesPGRepo)
	// TG bot

	if cfg.Telegram.IsEnabled {
		tgUsecase := tg.New(redisClient, usersPGRepo, jwtUseCase)
		bot, err := tele.NewBot(tele.Settings{
			Token: cfg.Telegram.Token,
		})

		if err != nil {
			log.Fatalf("Ошибка инициализации Telegram бота: %v", err)
		}
		telegram.NewBot(bot, tgUsecase)
		go bot.Start()
	}

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
		dishesUsecase,
		chefsUsecase,
		cartUsecase,
		ordersUsecase,
		shiftsUsecase,
		reviewsUsecase,
		favoritesUsecase,
	)

	err = handler.Run(fmt.Sprintf(":%s", cfg.HostConfig.Port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
