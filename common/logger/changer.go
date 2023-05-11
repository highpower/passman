package logger

import (
	"fmt"
	"mime"
	"net/http"
	"strings"
)

type LevelChanger struct {
	Logger Logger
}

func (l *LevelChanger) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writerLogger, ok := l.Logger.(*WriterLogger)
	if !ok {
		http.Error(writer, "no logger that supports level change here", http.StatusConflict)
		return
	}
	value := strings.TrimSpace(request.URL.Query().Get("level"))
	level, err := ParseLevel(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writerLogger.SetLevel(level)
	result := fmt.Sprintf("level %s set successfully\n", level.String())
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("content-type", mime.FormatMediaType("text/plain", map[string]string{"charset": "utf-8"}))
	_, _ = writer.Write([]byte(result))
}
