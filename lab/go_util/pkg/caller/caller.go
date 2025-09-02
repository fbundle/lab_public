package caller

import (
	"errors"
	"fmt"
	"runtime"
)

type Caller struct {
	Name string
	File string
	Line int
}

func (caller Caller) String() string {
	return fmt.Sprintf("%s\t%s:%d", caller.Name, caller.File, caller.Line)
}

func CallStack(skip int) (callers []Caller) {
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		name := "<unknown>"
		if fn != nil {
			name = fn.Name()
		}
		callers = append(callers, Caller{
			Name: name,
			File: file,
			Line: line,
		})
	}
	return callers
}

func CallStackError(skip int) error {
	callers := CallStack(skip + 1)
	msg := ""
	for _, c := range callers {
		msg += c.String() + "\n"
	}
	return errors.New(msg)
}
