package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	tp, tpErr := jaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ctx := context.Background()
	ctx, span := otel.Tracer("").Start(ctx, "Service.BrandSafety")
	defer span.End()

	r := gin.Default()
	//gin OpenTelemetry instrumentation
	r.Use(otelgin.Middleware("todo-service"))
	r.GET("/todo", func(c *gin.Context) {
		//Make sure to pass c.Request.Context() as the context and not c itself
		_, span := otel.Tracer("").Start(c.Request.Context(), "Parse")
		time.Sleep(time.Second)
		span.End()
		time.Sleep(time.Second)

		c.JSON(http.StatusOK, gin.H{"msg": "OK"})
	})
	_ = r.Run(":8080")
}

func jaegerTraceProvider() (*sdktrace.TracerProvider, error) {

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger-all-in-one:14268/api/traces")))
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("order-service"),
			attribute.String("environment", "development"),
		)),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1.0)),
	)

	return tp, nil
}
