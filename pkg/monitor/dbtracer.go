package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// RowsAffectedKey is a standard attribute for the number of rows affected
var RowsAffectedKey = attribute.Key("pgx.rows_affected")

type ctxKey string

type traceData struct {
	start time.Time
	sql   string
}

const traceCtxKey = ctxKey("pg_trace_ctx_key")

type dbTracer struct {
	tracer          trace.Tracer
	log             logger.Log
	attrs           []attribute.KeyValue
	logSQLStatement bool
	includeParams   bool
}

// NewDBTracer create new db tracer
func NewDBTracer() pgx.QueryTracer {
	return &dbTracer{}
}

func (t dbTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}

	if t.logSQLStatement {
		opts = append(opts, trace.WithAttributes(semconv.DBStatementKey.String(data.SQL)))
		if t.includeParams {
			opts = append(opts, trace.WithAttributes(makeParamsAttribute(data.Args)))
		}
	}

	spanName := "prepare " + data.SQL

	ctx, _ = t.tracer.Start(ctx, spanName, opts...)

	return context.WithValue(ctx, traceCtxKey, traceData{
		start: time.Now(),
		sql:   data.SQL,
	})
}

func (t dbTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)
	span.End()

	v, ok := ctx.Value(traceCtxKey).(traceData)
	if !ok {
		return
	}
	duration := time.Since(v.start).Milliseconds()

	t.log.Infof("END SQL(%dms): %s. Result: %s, Err: %+v", duration, v.sql, data.CommandTag.String(), data.Err)
}

// TraceCopyFromStart is called at the beginning of CopyFrom calls. The
// returned context is used for the rest of the call and will be passed to
// TraceCopyFromEnd.
func (t dbTracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(attribute.String("db.table", data.TableName.Sanitize())),
	}

	if conn != nil {
		opts = append(opts, connectionAttributesFromConfig(conn.Config())...)
	}

	ctx, _ = t.tracer.Start(ctx, "copy_from "+data.TableName.Sanitize(), opts...)

	return ctx
}

// TraceCopyFromEnd is called at the end of CopyFrom calls.
func (t dbTracer) TraceCopyFromEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)

	span.SetAttributes(RowsAffectedKey.Int(int(data.CommandTag.RowsAffected())))

	span.End()
}

func recordError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

func makeParamsAttribute(args []any) attribute.KeyValue {
	ss := make([]string, len(args))
	for i := range args {
		ss[i] = fmt.Sprintf("%+v", args[i])
	}
	// Since there doesn't appear to be a standard key for this in semconv, prefix it to avoid
	// clashing with future standard attributes.
	return attribute.Key("pgx.query.parameters").StringSlice(ss)
}

// connectionAttributesFromConfig returns a slice of SpanStartOptions that contain
// attributes from the given connection config.
func connectionAttributesFromConfig(config *pgx.ConnConfig) []trace.SpanStartOption {
	if config != nil {
		return []trace.SpanStartOption{
			trace.WithAttributes(attribute.String(string(semconv.NetPeerNameKey), config.Host)),
			trace.WithAttributes(attribute.Int(string(semconv.NetPeerPortKey), int(config.Port))),
			trace.WithAttributes(attribute.String(string(semconv.DBUserKey), config.User)),
		}
	}
	return nil
}
