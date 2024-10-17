package middleware

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type bodyWriter struct {
	body    *bytes.Buffer
	maxSize int
}

type bodyReader struct {
	io.ReadCloser
	body    *bytes.Buffer
	maxSize int
	bytes   int
}

var (
	TraceIDKey          = "trace_id"
	SpanIDKey           = "span_id"
	RequestBodyMaxSize  = 64 * 1024
	ResponseBodyMaxSize = 64 * 1024
)

func (w *bodyWriter) Write(b []byte) (int, error) {
	if w.body.Len()+len(b) > w.maxSize {
		return w.body.Write(b[:w.maxSize-w.body.Len()])
	}
	return w.body.Write(b)
}

func newBodyWriter(maxSize int) *bodyWriter {
	return &bodyWriter{
		body:    bytes.NewBufferString(""),
		maxSize: maxSize,
	}
}

func newBodyReader(reader io.ReadCloser, maxSize int) *bodyReader {
	body := bytes.NewBufferString("")
	return &bodyReader{
		ReadCloser: reader,
		body:       body,
		maxSize:    maxSize,
		bytes:      0,
	}
}

// Logger returns a `func(http.Handler) http.Handler` (middleware) that logs requests using slog.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			query := r.URL.RawQuery

			br := newBodyReader(r.Body, RequestBodyMaxSize)
			r.Body = br

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			bw := newBodyWriter(ResponseBodyMaxSize)
			ww.Tee(bw) // Tee allows capturing the response body

			defer func() {
				end := time.Now()
				latency := end.Sub(start)
				status := ww.Status()

				params := make(map[string]string)
				for i, k := range chi.RouteContext(r.Context()).URLParams.Keys {
					params[k] = chi.RouteContext(r.Context()).URLParams.Values[i]
				}

				userAgent := r.UserAgent()
				ip := r.RemoteAddr
				referer := r.Referer()

				// Request attributes
				requestAttributes := []slog.Attr{
					slog.Time("time", start.UTC()),
					slog.String("method", r.Method),
					slog.String("host", r.Host),
					slog.String("path", path),
					slog.String("query", query),
					slog.Any("params", params),
					slog.String("ip", ip),
					slog.String("referer", referer),
					slog.Int("length", br.bytes),
					slog.String("body", br.body.String()),
					slog.String("user-agent", userAgent),
				}

				responseAttributes := []slog.Attr{
					slog.Time("time", end.UTC()),
					slog.Duration("latency", latency),
					slog.Int("status", status),
					slog.String("body", bw.body.String()),
				}
				attributes := append(
					[]slog.Attr{
						{Key: "request", Value: slog.GroupValue(requestAttributes...)},
						{Key: "response", Value: slog.GroupValue(responseAttributes...)},
					},
					extractTraceSpanID(r.Context())...,
				)

				level := slog.LevelInfo
				if status >= http.StatusInternalServerError {
					level = slog.LevelError
				}
				logger.LogAttrs(r.Context(), level, strconv.Itoa(status)+": "+http.StatusText(status), attributes...)
			}()
			next.ServeHTTP(ww, r)
		})
	}
}

// extractTraceSpanID extracts the trace and span IDs from the given context.
// It returns a slice of slog.Attr containing the trace and span IDs if they exist.
func extractTraceSpanID(ctx context.Context) []slog.Attr {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	var attrs []slog.Attr
	spanCtx := span.SpanContext()
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID().String()
		attrs = append(attrs, slog.String(TraceIDKey, traceID))
	}
	if spanCtx.HasSpanID() {
		spanID := spanCtx.SpanID().String()
		attrs = append(attrs, slog.String(SpanIDKey, spanID))
	}
	return attrs
}
