package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type WriterLogger struct {
	level    Level
	writer   io.WriteCloser
	printers [4]*log.Logger
}

type closerAdapter struct {
	writer io.Writer
}

const (
	LevelFatal Level = iota
	LevelError
	LevelInfo
	LevelDebug
	LevelUnknown Level = -1
)

var ErrUnknownLevel = errors.New("unknown level")

func (l *WriterLogger) Level() Level {
	return Level(atomic.LoadInt32((*int32)(&l.level)))
}

func (l *WriterLogger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&l.level), int32(level))
}

func (l *WriterLogger) levelAllowed(level Level) bool {
	return l.Level() >= level
}

func (l *WriterLogger) Fatal(v ...any) {
	l.log(LevelFatal, v...)
}

//goland:noinspection SpellCheckingInspection

func (l *WriterLogger) Fatalf(format string, v ...any) {
	l.logf(LevelFatal, format, v...)
}

func (l *WriterLogger) Error(v ...any) {
	l.log(LevelError, v...)
}

//goland:noinspection SpellCheckingInspection

func (l *WriterLogger) Errorf(format string, v ...any) {
	l.logf(LevelError, format, v...)
}

func (l *WriterLogger) Info(v ...any) {
	l.log(LevelInfo, v...)
}

//goland:noinspection SpellCheckingInspection

func (l *WriterLogger) Infof(format string, v ...any) {
	l.logf(LevelInfo, format, v...)
}

func (l *WriterLogger) Debug(v ...any) {
	l.log(LevelDebug, v...)
}

//goland:noinspection SpellCheckingInspection

func (l *WriterLogger) Debugf(format string, v ...any) {
	l.logf(LevelDebug, format, v...)
}

func (l *WriterLogger) Close() error {
	return l.writer.Close()
}

func (l *WriterLogger) log(level Level, v ...any) {
	switch {
	case level == LevelFatal:
		l.printers[LevelFatal].Fatal(v...)
	case l.levelAllowed(level):
		l.printers[int(level)].Print(v...)
	}
}

func (l *WriterLogger) logf(level Level, format string, v ...any) {
	switch {
	case level == LevelFatal:
		l.printers[LevelFatal].Fatalf(format, v...)
	case l.levelAllowed(level):
		l.printers[int(level)].Printf(format, v...)
	}
}

func (c closerAdapter) Close() error {
	return nil
}

func (c closerAdapter) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}

func NewWriterLogger(writer io.WriteCloser, level Level) *WriterLogger {
	result := WriterLogger{
		level:  level,
		writer: writer,
		printers: [4]*log.Logger{
			log.New(writer, fmt.Sprintf("[%s] ", fatalString), log.LstdFlags),
			log.New(writer, fmt.Sprintf("[%s] ", errorString), log.LstdFlags),
			log.New(writer, fmt.Sprintf("[%s] ", infoString), log.LstdFlags),
			log.New(writer, fmt.Sprintf("[%s] ", debugString), log.LstdFlags)}}

	return &result
}

func NewOutputLogger(output string, level Level) (*WriterLogger, error) {
	o := strings.TrimSpace(output)
	switch {
	case strings.EqualFold(o, "stderr"):
		return NewWriterLogger(closerAdapter{writer: os.Stderr}, level), nil
	case strings.EqualFold(o, "stdout") || o == "-":
		return NewWriterLogger(closerAdapter{writer: os.Stdout}, level), nil
	default:
		file, err := os.OpenFile(filepath.Clean(output), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return nil, err
		}
		return NewWriterLogger(file, level), nil
	}
}

func MustNewOutputLogger(output string, level Level) *WriterLogger {
	writer, err := NewOutputLogger(output, level)
	if err != nil {
		panic(err)
	}
	return writer
}
