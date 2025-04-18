package pkg

import "log"

type Logger interface {
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
}
type LogLevel uint32

const (
	LogLevelDebug = LogLevel(0)
	LogLevelInfo  = LogLevel(1)
	LogLevelWarn  = LogLevel(2)
	LogLevelError = LogLevel(3)
)

var DefaultLogger Logger = &defaultLogger{
	logLevel: LogLevelInfo,
}

var DebugLogger Logger = &defaultLogger{
	logLevel: LogLevelDebug,
}

type defaultLogger struct {
	logLevel LogLevel
}

func (l *defaultLogger) Debugf(format string, a ...any) {
	if l.logLevel > LogLevelDebug {
		return
	}
	log.Printf("[Debug] "+format+"\n", a...)
}

func (l *defaultLogger) Infof(format string, a ...any) {
	if l.logLevel > LogLevelInfo {
		return
	}
	log.Printf("[Info] "+format+"\n", a...)
}

func (l *defaultLogger) Warnf(format string, a ...any) {
	if l.logLevel > LogLevelWarn {
		return
	}
	log.Printf("[Warn] "+format+"\n", a...)
}

func (l *defaultLogger) Errorf(format string, a ...any) {
	if l.logLevel > LogLevelError {
		return
	}
	log.Printf("[Error] "+format+"\n", a...)
}
