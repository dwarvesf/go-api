package monitor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestGetTraceID(t *testing.T) {
	testCases := map[string]struct {
		knownTraceID  trace.TraceID
		expectedTrace string
	}{
		"known trace ID": {
			knownTraceID:  trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c},
			expectedTrace: "0102030405060708090a0b0c00000000",
		},
		"zero trace ID": {
			knownTraceID:  trace.TraceID{},
			expectedTrace: "00000000000000000000000000000000",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Create a new trace.SpanContext with a known trace ID for testing
			spanContext := trace.SpanContext{}.WithTraceID(tc.knownTraceID)

			// Create a context with the known trace ID
			ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

			// Call the GetTraceID function with the context
			traceID := GetTraceID(ctx)

			// Assert that the returned traceID matches the expected trace ID in string format
			assert.Equal(t, tc.expectedTrace, traceID)
		})
	}
}
