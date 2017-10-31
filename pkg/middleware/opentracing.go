package middleware

import (
	"net/http"

	ddsContext "github.com/gregbiv/news-api/pkg/context"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	"github.com/palantir/stacktrace"
)

// OpenTracing middleware that will extract the spanContext from the request headers and set it into OpenTracingSpanContext.
func OpenTracing(tracer opentracing.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			spanOptions := []opentracing.StartSpanOption{
				opentracing.Tag{Key: "user_agent", Value: r.UserAgent()},
			}

			spanContext, err := tracer.Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header),
			)
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				ddsContext.Logger(r.Context()).Info(stacktrace.Propagate(err, "Failed to extract OT context from headers %v", r.Header))
			} else {
				spanOptions = append(spanOptions, opentracing.FollowsFrom(spanContext))
			}

			span := opentracing.GlobalTracer().StartSpan("mas_request", spanOptions...)
			defer span.Finish()

			ctx := r.Context()
			newSpanContext := span.Context()

			if basicSpanContext, ok := newSpanContext.(basictracer.SpanContext); ok {
				ctx = ddsContext.WithTraceID(ctx, basicSpanContext.TraceID)
			}

			ctx = opentracing.ContextWithSpan(ctx, span)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
