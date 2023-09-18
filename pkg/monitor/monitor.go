package monitor

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Tracer is the interface for monitor
type Tracer interface {
	// Start start a span
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	// Clean clean the monitor before shutdown the app to avoid memory leak
	Clean(timeout time.Duration)
}

// GetTraceID return the trace id
func GetTraceID(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
