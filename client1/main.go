package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	jp "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	e := echo.New()
	os.Setenv("JAEGER_SERVICE_NAME", "client1")

	// Middleware
	e.Use(middleware.Logger())

	e.Use(middleware.Recover())
	exp, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
	)
	if err != nil {
		panic(err)
	}

	// Initialize new provider with the exporter
	// provider := xshtracing.InitProvider(exp)

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String("client1"),
		)),
	)
	propagator := jp.Jaeger{}
	otel.SetTextMapPropagator(propagator)

	e.Use(otelecho.Middleware(
		"http.request",
		otelecho.WithTracerProvider(tp),
	))
	// Routes
	e.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:2323", nil)
		if err != nil {
			panic(err)
		}
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		span := trace.SpanFromContext(ctx)
		fmt.Println("SPANFROMCONTEXT", span.SpanContext().TraceID.String())

		ctx, req = otelhttptrace.W3C(ctx, req)
		otelhttptrace.Inject(ctx, req)
		_, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
