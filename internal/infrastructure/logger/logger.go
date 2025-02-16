package logger

import (
	"github.com/icefed/zlog"
	"log/slog"
)

func NewLogger(devMode bool, logLevel slog.Level) *zlog.Logger {
	h := zlog.NewJSONHandler(&zlog.Config{
		HandlerOptions: slog.HandlerOptions{
			Level: logLevel,
		},
		Development: devMode,
	})
	return zlog.New(h)
}
