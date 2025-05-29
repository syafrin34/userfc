package logger

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var Logger *logrus.Logger

type Fields = logrus.Fields

func SetupLogger() {
	logging := logrus.New()
	logging.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logging.Info("logged initiated using logrus")
	Logger = logging
}

func LogWithTrace(ctx context.Context) logrus.Entry {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	return *logrus.WithField("trace_id", traceID)

}
