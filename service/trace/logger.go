package trace

import (
	"fmt"

	"github.com/icon-project/goloop/common/log"
	"github.com/icon-project/goloop/module"
)

type Logger struct {
	log.Logger
	isTrace bool
	prefix  string
	onLog   func(lv module.TraceLevel, msg string)
}

func (l *Logger) IsTrace() bool {
	return l.isTrace
}

func (l *Logger) TLog(lv module.TraceLevel, a ...interface{}) {
	if l.isTrace {
		l.onLog(lv, l.prefix+fmt.Sprint(a...))
	}
}

func (l *Logger) TLogln(lv module.TraceLevel, a ...interface{}) {
	if l.isTrace {
		l.onLog(lv, l.prefix+fmt.Sprint(a...))
	}
}

func (l *Logger) TLogf(lv module.TraceLevel, f string, a ...interface{}) {
	if l.isTrace {
		l.onLog(lv, l.prefix+fmt.Sprintf(f, a...))
	}
}

func (l *Logger) TDebug(a ...interface{}) {
	l.TLog(module.TDebugLevel, a...)
}

func (l *Logger) TDebugln(a ...interface{}) {
	l.TLogln(module.TDebugLevel, a...)
}

func (l *Logger) TDebugf(f string, a ...interface{}) {
	l.TLogf(module.TDebugLevel, f, a...)
}

func (l *Logger) TTrace(a ...interface{}) {
	l.TLog(module.TTraceLevel, a...)
}

func (l *Logger) TTraceln(a ...interface{}) {
	l.TLogln(module.TTraceLevel, a...)
}

func (l *Logger) TTracef(f string, a ...interface{}) {
	l.TLogf(module.TTraceLevel, f, a...)
}

func (l *Logger) TSystem(a ...interface{}) {
	l.TLog(module.TSystemLevel, a...)
}

func (l *Logger) TSystemln(a ...interface{}) {
	l.TLogln(module.TSystemLevel, a...)
}

func (l *Logger) TSystemf(f string, a ...interface{}) {
	l.TLogf(module.TSystemLevel, f, a...)
}

func (l *Logger) WithFields(f log.Fields) log.Logger {
	return &Logger{
		Logger:  l.Logger.WithFields(f),
		isTrace: l.isTrace,
		prefix:  l.prefix,
		onLog:   l.onLog,
	}
}

func (l *Logger) TPrefix() string {
	return l.prefix
}

func (l *Logger) WithTPrefix(prefix string) *Logger {
	return &Logger{
		Logger:  l.Logger,
		isTrace: l.isTrace,
		prefix:  prefix,
		onLog:   l.onLog,
	}
}

func dummyLog(lv module.TraceLevel, msg string) {
	// do nothing
}

func NewLogger(l log.Logger, t module.TraceCallback) *Logger {
	if t != nil {
		return &Logger{
			Logger:  l,
			isTrace: true,
			onLog:   t.OnLog,
		}
	} else {
		return &Logger{
			Logger: l,
			onLog:  dummyLog,
		}
	}
}

func LoggerOf(l log.Logger) *Logger {
	if logger, ok := l.(*Logger); ok {
		return logger
	}
	return NewLogger(l, nil)
}
