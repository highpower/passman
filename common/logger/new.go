package logger

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrLoggerType = errors.New("bad logger type")

func New(config Config) (Logger, error) {
	t := strings.TrimSpace(config.LogType())
	switch {
	case strings.EqualFold(t, "syslog"):
		return NewSyslogLogger(config.LogOutput())
	case strings.EqualFold(t, "writer"):
		return NewOutputLogger(config.LogOutput(), config.LogLevel())
	case strings.EqualFold(t, "discard"):
		return NewDiscardLogger(), nil
	default:
		return nil, fmt.Errorf("%w '%s'", ErrLoggerType, t)
	}
}

func NewDiscardLogger() Logger {
	return NewWriterLogger(closerAdapter{writer: io.Discard}, LevelError)
}
