package extension

import (
	"fmt"
	"log"
	"strings"
)

type Logger struct{}

func (l *Logger) Print(level string, format string, params ...any) {
	format = fmt.Sprintf("[%s]", level) + format
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	log.Printf(format, params...)
}

func (l *Logger) Error(format string, params ...any) {
	l.Print("error", format, params...)
}

func (l *Logger) Warn(format string, params ...any) {
	l.Print("warn", format, params...)
}

func (l *Logger) Info(format string, params ...any) {
	l.Print("info", format, params...)
}

func (l *Logger) Debug(format string, params ...any) {
	l.Print("debug", format, params...)
}

func (l *Logger) Trace(format string, params ...any) {
	l.Print("trace", format, params...)
}
