package loglevel

import (
	"log/slog"
)

func GetLogLevelFromString(s string) slog.Level {
	switch s {
	case "ERR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "DEBUG":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
