package context

import (
	"context"
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
)

type indexContext int

const (
	traceIDKey indexContext = iota
)

var (
	baseLogger logrus.FieldLogger
)

func init() {
	baseLogger = logrus.New()
}

// SetLogger receives a logger and stores it as the base logger for the package.
func SetLogger(logger logrus.FieldLogger) {
	baseLogger = logger
}

// WithTraceID returns a copy of the parent context but with the TraceID stored on it.
func WithTraceID(ctx context.Context, traceID uint64) context.Context {
	return context.WithValue(ctx, traceIDKey, strconv.FormatUint(traceID, 16))
}

// Logger returns an logger with all values from context loaded on it.
func Logger(ctx context.Context) logrus.FieldLogger {
	if id, err := traceID(ctx); err == nil {
		return baseLogger.WithField("traceID", id)
	}

	return baseLogger
}

func traceID(ctx context.Context) (string, error) {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "", errors.New("No TraceID in the context")
	}

	return traceID, nil
}
