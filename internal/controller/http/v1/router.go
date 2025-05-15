package v1

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"

	"domashka-backend/internal/utils/telemetry"
)

func NewRouter(
	handler *gin.Engine,
	l logger,
	u usersUsecase,
	a authUsecase,
	jwt jwtUsecase,
	n notificationUsecase,
	g geoUsecase,
	dishesUsecase dishesUsecase,
	chefsUsecase chefUsecase,
	cartUsecase cartUsecase,
	orderUsecase orderUsecase,
	shiftsUsecase shiftsUsecase,
	reviewsUsecase reviewsUsecase,
	favoritesUsecase favoritesUsecase,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	tp := telemetry.TracerProvider()
	defer func() { _ = tp.Shutdown(context.Background()) }()

	exporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint("otel-collector:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(time.Second*5))),
	)
	otel.SetMeterProvider(provider)

	handler.Use(telemetry.OtelGinMiddleware())

	meter := otel.Meter("domashka-app")
	requestCounter, _ := meter.Int64Counter(
		"test",
	)
	handler.GET("/health", func(c *gin.Context) {
		requestCounter.Add(context.Background(), 400)
		c.JSON(200, gin.H{})
	})

	// Routers
	h := handler.Group("/v1")
	{

		newUsersHandler(h, l, u, reviewsUsecase, orderUsecase, favoritesUsecase, dishesUsecase, chefsUsecase)
		newAuthHandler(h, a, jwt)
		authorized := h.Group("/")
		authorized.Use(AuthMiddleware(jwt))
		{
			newNotificationHandler(authorized, n)
			RegisterGeoHandlers(h, g)
			NewChefsHandler(authorized, dishesUsecase, chefsUsecase, g, shiftsUsecase, u, reviewsUsecase, orderUsecase)
			RegisterSearchHandler(authorized, dishesUsecase, chefsUsecase, orderUsecase, reviewsUsecase, u)
			NewDishesHandler(authorized, dishesUsecase, chefsUsecase, u)
		}
		RegisterCartHandlers(h, chefsUsecase, cartUsecase, dishesUsecase, g)
		RegisterOrderHandlers(h, g, cartUsecase, orderUsecase, shiftsUsecase, chefsUsecase)
		NewHomeHandler(authorized, jwt, g, dishesUsecase, chefsUsecase, orderUsecase, reviewsUsecase)
		RegisterReviewHandlers(h, reviewsUsecase)
	}
}
