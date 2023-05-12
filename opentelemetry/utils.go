package opentelemetry

import (
	"context"
	"os"

	"github.com/zhs007/goutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func ConfigureJaeger(ctx context.Context) (func(), error) {
	if os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT") == "" {
		return nil, nil
	}

	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		goutils.Error("ConfigureJaeger",
			zap.Error(err))

		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			goutils.Error("ConfigureJaeger:Shutdown",
				zap.Error(err))
		}
	}, nil
}
