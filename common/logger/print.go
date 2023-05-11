package logger

import "fmt"

func PrintConfig(config Config) string {
	return fmt.Sprintf("logger[type=%s,level=%s,output=%s]",
		config.LogType(), PrintLevel(config.LogLevel()), config.LogOutput())
}
