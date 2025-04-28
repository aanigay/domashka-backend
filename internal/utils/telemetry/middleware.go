package telemetry

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func OtelGinMiddleware() gin.HandlerFunc {
	meter := otel.Meter("domashka-gin")

	requestsCounter, _ := meter.Int64Counter(
		"http_server_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)

	requestDuration, _ := meter.Float64Histogram(
		"http_server_requests_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
	)

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()), // автоматически берёт роут вида /ping, /users/:id
			attribute.Int("http.status_code", c.Writer.Status()),
		}

		ctx := context.Background()

		requestsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		requestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
	}
}
