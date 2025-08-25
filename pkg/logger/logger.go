package logger

import (
    "log"
    "strings"
)

type Logger struct {
    level int
}

const (
    DEBUG = 0
    INFO  = 1
    WARN  = 2
    ERROR = 3
)

func New(level string) *Logger {
    l := &Logger{}
    switch strings.ToLower(level) {
    case "debug":
        l.level = DEBUG
    case "warn":
        l.level = WARN
    case "error":
        l.level = ERROR
    default:
        l.level = INFO
    }
    return l
}

func (l *Logger) Debug(msg string, args ...any) {
    if l.level <= DEBUG {
        log.Printf("[DEBUG] "+msg, args...)
    }
}

func (l *Logger) Info(msg string, args ...any) {
    if l.level <= INFO {
        log.Printf("[INFO] "+msg, args...)
    }
}

func (l *Logger) Warn(msg string, args ...any) {
    if l.level <= WARN {
        log.Printf("[WARN] "+msg, args...)
    }
}

func (l *Logger) Error(msg string, args ...any) {
    if l.level <= ERROR {
        log.Printf("[ERROR] "+msg, args...)
    }
}