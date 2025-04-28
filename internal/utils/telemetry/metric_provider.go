package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"log"
	"time"
)

func MeterProvider() *metric.MeterProvider {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx,
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
	meter := otel.Meter("domashka-app")
	requestCounter, _ := meter.Int64Counter(
		"http_server_requests_seconds_count",
	)
	requestCounter.Add(ctx, 400)

	return provider
}
