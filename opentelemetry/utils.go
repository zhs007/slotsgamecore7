package opentelemetry

import (
	"context"
	"os"

	"github.com/zhs007/goutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func ConfigureJaeger(ctx context.Context, appname string) (func(), error) {
	if os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT") == "" {
		return nil, nil
	}

	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		goutils.Error("ConfigureJaeger",
			goutils.Err(err))

		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(appname))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			goutils.Error("ConfigureJaeger:Shutdown",
				goutils.Err(err))
		}
	}, nil
}
