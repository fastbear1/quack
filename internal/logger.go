package utils

import (
	"fmt"
	"strings"
)

const (
	INFO = iota
	DEBUG
)

type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
}

type LogSet struct {
	LogLevel int
}

// Gel logger once by his name
var loggingMap = map[string]*LogSet{}

func GetLogger(name string, level any) Logger {
	logName := strings.ToUpper(name)
	if lg, ok := loggingMap[logName]; ok {
		return lg
	}
	logLevel := level.(int)
	log := LogSet{LogLevel: logLevel}
	loggingMap[logName] = &log
	return &log
}

func (l *LogSet) formatMessage(message string, args ...any) string {
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}
	return message
}

func (l *LogSet) Info(message string, args ...any) {
	fmt.Println(message)
}

func (l *LogSet) Debug(message string, args ...any) {
	if l.LogLevel == DEBUG {
		fmt.Println(message)
	}
}
