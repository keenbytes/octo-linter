// Package loglevel offers utility functions for working with log levels.
package loglevel

import (
	"log/slog"
)

// GetLogLevelFromString takes a string used in a CLI flag and returns the corresponding slog log level constant.
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
