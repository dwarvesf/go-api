package logger

import (
	"log"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/getsentry/sentry-go"
)

// NewSentry return the sentry client
func NewSentry(cfg *config.Config) (*sentry.Client, error) {
	if cfg.SentryDSN == "" {
		log.Println("Sentry DSN not provided. Not using Sentry Error Reporting")
		return nil, nil
	}

	client, err := sentry.NewClient(
		sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			AttachStacktrace: true,
			SampleRate:       1,
			ServerName:       cfg.ServerName,
			Release:          cfg.Version,
			Environment:      cfg.Env,
			EnableTracing:    true,
		},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Sentry Error Reporter initialized")

	return client, nil
}
