package utils

import (
	"fmt"
	"strings"
)

const (
	INFO = iota + 1
	DEBUG
)

type Logger interface {
	Info(string, ...any)
	Debug(any, ...any)
}

type LogSet struct {
	LogLevel int
}

// Gel logger once by his name
var LoggingMap map[string]*LogSet

func init() {
	LoggingMap = make(map[string]*LogSet, 0)
}

func GetLogger(name string, level any) Logger {
	logName := strings.ToUpper(name)
	if lg, ok := LoggingMap[logName]; ok {
		if level != nil {
			lg.LogLevel = level.(int)
		}
		return lg
	}
	logLevel := INFO
	if level != nil {
		logLevel = level.(int)
	}

	log := LogSet{LogLevel: logLevel}
	LoggingMap[logName] = &log
	return &log
}

func (l *LogSet) formatMessage(message string, args ...any) string {
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}
	return message
}

func (l *LogSet) Info(message string, args ...any) {
	logLine := l.formatMessage(message, args...)
	fmt.Println(logLine)
}

func (l *LogSet) Debug(message any, args ...any) {
	if l.LogLevel == DEBUG {
		if err, ok := message.(error); ok {
			fmt.Println(err)
		} else {
			msg, _ := message.(string)
			logLine := l.formatMessage(msg, args...)
			fmt.Println(logLine)
		}
	}
}
