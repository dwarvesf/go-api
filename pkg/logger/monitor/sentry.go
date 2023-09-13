package monitor

import (
	"context"
	"log"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// sentryTracer is the monitor for sentry
type sentryTracer struct {
	client *sentry.Client
}

// NewSentry return the sentry monitor
func NewSentry(cfg *config.Config) (Tracer, error) {
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

	return &sentryTracer{
		client: client,
	}, nil
}

func (s *sentryTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// TODO: validate when otel.Tracer is nil
	return otel.Tracer(spanName).Start(ctx, spanName, opts...)
}

func (s *sentryTracer) Clean(timeout time.Duration) {
	s.client.Flush(timeout)
}
