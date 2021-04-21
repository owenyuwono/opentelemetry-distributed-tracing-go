package main

import (
	"fmt"
	"net/http/httputil"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
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
	os.Setenv("JAEGER_SERVICE_NAME", "client2")

	// Middleware
	e.Use(middleware.Logger())

	e.Use(middleware.Recover())
	exp, err := jaeger.NewRawExporter(jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"))
	if err != nil {
		panic(err)
	}

	// Initialize new provider with the exporter
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
		span := trace.SpanFromContext(c.Request().Context())
		fmt.Println(span, span.SpanContext().TraceID.String())
		x, _ := httputil.DumpRequest(c.Request(), false)
		fmt.Println(string(x))
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(":2323"))
}
