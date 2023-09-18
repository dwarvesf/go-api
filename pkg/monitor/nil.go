package monitor

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type nilSpan struct{}

func (nilSpan) End(options ...trace.SpanEndOption)                  {}
func (nilSpan) AddEvent(name string, options ...trace.EventOption)  {}
func (nilSpan) IsRecording() bool                                   { return false }
func (nilSpan) RecordError(err error, options ...trace.EventOption) {}
func (nilSpan) SpanContext() trace.SpanContext                      { return trace.SpanContext{} }
func (nilSpan) SetStatus(code codes.Code, description string)       {}
func (nilSpan) SetName(name string)                                 {}
func (nilSpan) SetAttributes(kv ...attribute.KeyValue)              {}
func (nilSpan) TracerProvider() trace.TracerProvider                { return nil }

type nilMonitor struct{}

func (s *nilMonitor) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, nilSpan{}
}

func (s *nilMonitor) Clean(timeout time.Duration) {
}

// TestMonitor return the nil monitor
func TestMonitor() Tracer {
	return &nilMonitor{}
}
