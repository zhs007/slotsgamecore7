package opentelemetry

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/zhs007/goutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// chooseSampler selects a sampler based on environment variables with sane defaults.
// Priority:
// 1) OTEL_TRACES_SAMPLER = always_on|always_off|traceidratio (with OTEL_TRACES_SAMPLER_ARG)
// 2) SGC7_TRACE_SAMPLE_RATIO (0..1) -> ParentBased(TraceIDRatioBased(r))
// 3) If ENV/APP_ENV/GO_ENV indicates prod/production -> ParentBased(TraceIDRatioBased(0.1))
// 4) Default -> AlwaysSample (developer-friendly)
func chooseSampler() sdktrace.Sampler {
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER"))); v != "" {
		switch v {
		case "always_on":
			return sdktrace.AlwaysSample()
		case "always_off":
			return sdktrace.NeverSample()
		case "traceidratio":
			if arg := os.Getenv("OTEL_TRACES_SAMPLER_ARG"); arg != "" {
				if r, err := strconv.ParseFloat(arg, 64); err == nil && r >= 0 && r <= 1 {
					return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(r))
				}
			}
			// default ratio when traceidratio is selected without arg
			return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))
		}
	}

	if arg := os.Getenv("SGC7_TRACE_SAMPLE_RATIO"); arg != "" {
		if r, err := strconv.ParseFloat(arg, 64); err == nil && r >= 0 && r <= 1 {
			return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(r))
		}
	}

	env := strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	}
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(os.Getenv("GO_ENV")))
	}
	if env == "prod" || env == "production" {
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))
	}

	return sdktrace.AlwaysSample()
}

// ConfigureOTLP sets up tracing using the OTLP exporter (recommended). It honors standard OTEL_* env vars.
func ConfigureOTLP(ctx context.Context, appname string) (func(), error) {
	// Use default environment-based configuration for endpoint, headers, etc.
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		goutils.Error("ConfigureOTLP:otlptracegrpc.New", goutils.Err(err))
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(chooseSampler()),
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
		sdktrace.WithSampler(chooseSampler()),
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
