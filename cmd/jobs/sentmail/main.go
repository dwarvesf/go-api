package main

import (
	"context"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/controller/user"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	cfg := config.LoadConfig(config.DefaultConfigLoaders())

	l := logger.NewLogByConfig(cfg)
	l.Infof("Cronjob starting")

	_, err := db.Init(*cfg)
	if err != nil {
		l.Fatal(err, "failed to init db")
	}

	sentryMonitor, err := monitor.NewSentry(cfg)
	if err != nil {
		l.Fatal(err, "failed to init sentry")
	}

	defer sentryMonitor.Clean(2 * time.Second)

	// new span for sentmail job
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("job", "SentMail")),
		trace.WithSpanKind(trace.SpanKindUnspecified),
	}
	spanName := "main"
	ctx, span := sentryMonitor.Start(context.Background(), spanName, opts...)
	defer span.End()

	// new controler
	c := user.NewUserController(*cfg, repository.NewRepo(), sentryMonitor)
	c.SentMail(ctx)
}
