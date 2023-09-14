package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
)

func TestNewSentry(t *testing.T) {
	tests := map[string]struct {
		cfg       *config.Config
		expectErr bool
	}{
		"empty dsn": {
			cfg: &config.Config{
				SentryDSN: "",
			},
			expectErr: false,
		},
		"valid dsn": {
			cfg: &config.Config{
				SentryDSN: "https://examplePublicKey@o0.ingest.sentry.io/0",
			},
			expectErr: false,
		},
		"invalid dsn": {
			cfg: &config.Config{
				SentryDSN: "invalid",
			},
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewSentry(tc.cfg)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSentryTracer_Start(t *testing.T) {
	tests := map[string]struct {
		spanName string
	}{
		"empty span name": {
			spanName: "",
		},
		"valid span name": {
			spanName: "testSpan",
		},
	}

	cfg := &config.Config{
		SentryDSN: "https://examplePublicKey@o0.ingest.sentry.io/0",
	}

	tracer, _ := NewSentry(cfg)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _ := tracer.Start(context.Background(), tc.spanName)
			if ctx == nil {
				t.Errorf("Expected context, got nil")
			}
		})
	}
}

func TestSentryTracer_Clean(t *testing.T) {
	tests := map[string]struct {
		timeout time.Duration
	}{
		"zero timeout": {
			timeout: 0,
		},
		"positive timeout": {
			timeout: 2 * time.Second,
		},
	}

	cfg := &config.Config{
		SentryDSN: "https://examplePublicKey@o0.ingest.sentry.io/0",
	}

	tracer, _ := NewSentry(cfg)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tracer.Clean(tc.timeout)
		})
	}
}
