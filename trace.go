package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("MyGinApp"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func TraceFunc(ctx context.Context, spanName string, fn func(context.Context) error) error {
	tracer := otel.Tracer("dns_say_what")
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	return fn(ctx)
}
