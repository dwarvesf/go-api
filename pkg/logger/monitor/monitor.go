package monitor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Exporter is the interface for monitor
type Exporter interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	PanicAlarmMiddleware() gin.HandlerFunc
}

// SentryExporter is the monitor for sentry
type SentryExporter struct {
	Client *sentry.Client
}

// New return the sentry client
func NewSentry(cfg *config.Config) (*SentryExporter, error) {
	if cfg.SentryDSN == "" {
		log.Println("Sentry DSN not provided. Not using Sentry Error Reporting")
		return nil, nil
	}

	hub := sentry.CurrentHub()
	client, err := sentry.NewClient(
		sentry.ClientOptions{
			Dsn:                cfg.SentryDSN,
			Environment:        cfg.Env,
			EnableTracing:      true,
			SampleRate:         1.0,
			TracesSampleRate:   1.0,
			ProfilesSampleRate: 1.0,
			Debug:              true,
			AttachStacktrace:   true,
			ServerName:         cfg.ServerName,
			Release:            cfg.Version,
		},
	)
	if err != nil {
		return nil, err
	}
	hub.BindClient(client)
	log.Println("Sentry Error Reporter initialized")

	// Opentelemetry
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())
	log.Println("Sentry Opentelemetry initialized")

	return &SentryExporter{
		Client: client,
	}, nil
}

func (s *SentryExporter) NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer(name).Start(ctx, name)
}

func (s *SentryExporter) PanicAlarmMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Panic occurred:", err)
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 5)

				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
					"traceID": trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String(),
				})
			}
		}()

		c.Next()
	}
}
