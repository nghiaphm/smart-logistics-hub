package logger

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}

func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(slog.String("request_id", requestID))
}
