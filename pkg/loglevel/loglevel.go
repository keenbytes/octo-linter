package loglevel

const (
	LogLevelNone = iota
	LogLevelOnlyErrors
	LogLevelErrorsAndWarnings
	LogLevelDebug
)

func GetLogLevelFromString(s string) int {
	switch s {
	case "NONE":
		return LogLevelNone
	case "ERR":
		return LogLevelOnlyErrors
	case "WARN":
		return LogLevelErrorsAndWarnings
	case "DEBUG":
		return LogLevelDebug
	default:
		return LogLevelOnlyErrors
	}
}
