package logger

import (
	"fmt"
	"strings"
)

type Level int32

type Config interface {
	LogType() string
	LogLevel() Level
	LogOutput() string
}

//goland:noinspection SpellCheckingInspection

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...any)

	Error(v ...interface{})
	Errorf(format string, v ...any)

	Info(v ...interface{})
	Infof(format string, v ...any)

	Debug(v ...interface{})
	Debugf(format string, v ...any)
	Close() error
}

const (
	fatalString = "FATAL"
	errorString = "ERROR"
	infoString  = "INFO"
	debugString = "DEBUG"
)

func (l *Level) String() string {
	return fmt.Sprintf("Level[%s]", l.Rep())
}

func (l *Level) MarshalText() ([]byte, error) {
	return []byte(l.Rep()), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	parsed, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*l = parsed
	return nil
}

func (l *Level) Rep() string {
	switch *l {
	case LevelFatal:
		return fatalString
	case LevelError:
		return errorString
	case LevelInfo:
		return infoString
	case LevelDebug:
		return debugString
	}
	panic("unknown level")
}

func ParseLevel(str string) (Level, error) {
	switch {
	case strings.EqualFold(str, fatalString):
		return LevelFatal, nil
	case strings.EqualFold(str, errorString):
		return LevelError, nil
	case strings.EqualFold(str, infoString):
		return LevelInfo, nil
	case strings.EqualFold(str, debugString):
		return LevelDebug, nil
	default:
		return LevelUnknown, ErrUnknownLevel
	}
}

func PrintLevel(l Level) string {
	return l.Rep()
}
