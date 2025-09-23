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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// ConfigureOTLP sets up tracing using the OTLP exporter (recommended). It honors standard OTEL_* env vars.
func ConfigureOTLP(ctx context.Context, appname string) (func(), error) {
	// Use default environment-based configuration for endpoint, headers, etc.
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		goutils.Error("ConfigureOTLP:otlptracegrpc.New", goutils.Err(err))
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(appname))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			goutils.Error("ConfigureOTLP:Shutdown", goutils.Err(err))
		}
	}, nil
}

// ConfigureJaeger is kept for backward compatibility, but Jaeger exporter is deprecated upstream.
// Prefer ConfigureOTLP with a Collector routing to Jaeger.
func ConfigureJaeger(ctx context.Context, appname string) (func(), error) {
	if os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT") == "" {
		return nil, nil
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		goutils.Error("ConfigureJaeger", goutils.Err(err))
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
			goutils.Error("ConfigureJaeger:Shutdown", goutils.Err(err))
		}
	}, nil
}
