package logger_test

import (
	"bytes"
	"fmt"
	"github.com/khanh-nguyen-code/go_util/pkg/logger"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	l := logger.NewDefaultLogger()
	l.Now().Info("hello")
	l.Now().WithField("foo", 123).Error("good bye")

	b := bytes.NewBuffer(make([]byte, 0))
	info := func(msg string) {
		b.Write([]byte(msg))
		_, _ = fmt.Fprintf(os.Stdout, msg)
	}
	err := func(msg string) {
		b.Write([]byte(msg))
		_, _ = fmt.Fprintf(os.Stderr, msg)
	}
	l = logger.NewLogger(info, err)
	l.Now().Info("hello")
	l.Now().WithField("foo", 123).Error("good bye")
}
