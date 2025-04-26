package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	timeLayout = "2006/01/02 15:04:05"
)

type Logger interface {
	Now() context
}

type context interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	WithField(key string, val interface{}) context
}

func NewLogger(info func(string), err func(string)) Logger {
	return &logger{
		info: info,
		err:  err,
	}
}
func NewDefaultLogger() Logger {
	return &logger{
		info: func(msg string) {
			_, _ = fmt.Fprintf(os.Stdout, msg)
		},
		err: func(msg string) {
			_, _ = fmt.Fprintf(os.Stderr, msg)
		},
	}
}

type logger struct {
	info func(string)
	err  func(string)
}

func (l *logger) Now() context {
	_, file, line, _ := runtime.Caller(1)
	return &loggerCtx{
		l:      l,
		now:    time.Now(),
		kvList: []*kv{{"caller", fmt.Sprintf("%s:%d", file, line)}},
	}
}

type kv struct {
	key string
	val interface{}
}

type loggerCtx struct {
	l      *logger
	now    time.Time
	kvList []*kv
}

func (c *loggerCtx) toLog(msg string) string {
	kvStrList := make([]string, 0, len(c.kvList))
	for _, kv := range c.kvList {
		b, _ := json.Marshal(kv.val)
		s := string(b)
		kvStrList = append(kvStrList, fmt.Sprintf("%s=%s", kv.key, s))
	}
	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		msg += "\n"
	}
	return fmt.Sprintf("%s|%s|%s", c.now.Format(timeLayout), strings.Join(kvStrList, "|"), msg)
}

func (c *loggerCtx) Info(format string, args ...interface{}) {
	if c.l.info != nil {
		c.l.info(c.toLog(fmt.Sprintf(format, args...)))
	}
}

func (c *loggerCtx) Error(format string, args ...interface{}) {
	if c.l.err != nil {
		c.l.err(c.toLog(fmt.Sprintf(format, args...)))
	}
}

func (c *loggerCtx) WithField(key string, val interface{}) context {
	c.kvList = append(c.kvList, &kv{key, val})
	return c
}
