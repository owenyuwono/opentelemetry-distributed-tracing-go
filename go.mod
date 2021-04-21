module test-apm

go 1.16

require (
	github.com/labstack/echo/v4 v4.2.2
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.15.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.15.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.15.1
	go.opentelemetry.io/contrib/propagators v0.15.1
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
)
