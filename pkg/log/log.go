package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

// L is the default logger that will be return if none are found within a context
var L = logrus.NewEntry(logrus.StandardLogger())

// WithFields embeds a logger within the context to structed context based logging
func WithFields(ctx context.Context, fields logrus.Fields) context.Context {
	logger := ctx.Value(loggerKey{})
	if logger == nil {
		logger = L
	}
	return context.WithValue(ctx, loggerKey{}, logger.(*logrus.Entry).WithFields(fields))
}

// GetLogger extracts the embedd logger from the context will be used by many RPC calls
func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		logger = L
	}

	return logger.(*logrus.Entry)
}
