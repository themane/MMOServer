package constants

import (
	"log"
)

const (
	SEPARATOR = " - "
	WARN      = "WARN"
	INFO      = "INFO"
)

type LoggingUtils struct {
	loggerName string
	level      string
}

func NewLoggingUtils(loggerName string, level string) *LoggingUtils {
	return &LoggingUtils{
		loggerName: loggerName,
		level:      level,
	}
}

func (l *LoggingUtils) info(message string) {
	log.Println(l.loggerName + SEPARATOR + message)
}

func (l *LoggingUtils) warnInfo(message string) {
	if l.level == WARN {
		log.Println(l.loggerName + SEPARATOR + message)
	}
}

func (l *LoggingUtils) error(message string, err error) {
	log.Println(l.loggerName+SEPARATOR+message, err)
}

func (l *LoggingUtils) warn(message string, err error) {
	if l.level == WARN {
		log.Println(l.loggerName+SEPARATOR+message, err)
	}
}
