package logger

import (
	"fmt"
	"log/syslog"
)

type syslogLogger struct {
	writer *syslog.Writer
}

func (s syslogLogger) Fatal(v ...any) {
	_ = s.writer.Crit(fmt.Sprint(v...))
}

//goland:noinspection SpellCheckingInspection

func (s syslogLogger) Fatalf(format string, v ...any) {
	_ = s.writer.Crit(fmt.Sprintf(format, v...))
}

func (s syslogLogger) Error(v ...any) {
	_ = s.writer.Err(fmt.Sprint(v...))
}

//goland:noinspection SpellCheckingInspection

func (s syslogLogger) Errorf(format string, v ...any) {
	_ = s.writer.Err(fmt.Sprintf(format, v...))
}

func (s syslogLogger) Info(v ...any) {
	_ = s.writer.Info(fmt.Sprint(v...))
}

//goland:noinspection SpellCheckingInspection

func (s syslogLogger) Infof(format string, v ...any) {
	_ = s.writer.Info(fmt.Sprintf(format, v...))
}

func (s syslogLogger) Debug(v ...any) {
	_ = s.writer.Debug(fmt.Sprint(v...))
}

func (s syslogLogger) Close() error {
	return s.writer.Close()
}

//goland:noinspection SpellCheckingInspection

func (s syslogLogger) Debugf(format string, v ...any) {
	_ = s.writer.Debug(fmt.Sprintf(format, v...))
}

func NewSyslogLogger(name string) (Logger, error) {
	writer, err := syslog.New(syslog.LOG_USER, name)
	if err != nil {
		return nil, err
	}
	return syslogLogger{writer: writer}, nil
}
