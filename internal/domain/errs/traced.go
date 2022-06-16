package errs

import (
	"errors"
	"fmt"
	"runtime"
)

const stringType = "string"

// compatible with error interface
type ErrorWithTrace struct {
	interrupt bool
	text      string
}

func New(err interface{}) ErrorWithTrace {
	e := ErrorWithTrace{
		interrupt: false,
	}

	s, ok := err.(string)
	if ok {
		e.text = e.trace(errors.New(s)).Error()
	} else {
		er, ok := err.(error)
		if ok {
			e.text = e.trace(er).Error()
		}
	}

	return e
}

func (e ErrorWithTrace) Error() string {
	return e.text
}

func (e ErrorWithTrace) IsInterrupt() bool {
	return e.interrupt
}

func (e ErrorWithTrace) Interrupt() ErrorWithTrace {
	e.interrupt = true
	return e
}

func (e ErrorWithTrace) trace(err error) error {
	pc := make([]uintptr, 15)

	n := runtime.Callers(3, pc)
	f := runtime.CallersFrames(pc[:n])

	errF1, _ := f.Next()
	errF2, _ := f.Next()

	return errors.New(
		fmt.Sprintf(
			"\n%s\n%s\n======== Error: ========\nValue: %s\n========================",
			fmt.Sprintf(
				"======= TRACE: 1 =======\nLine: %d\nFile: %s,\nMethod: %s",
				errF2.Line,
				errF2.File,
				errF2.Function,
			),
			fmt.Sprintf(
				"======= TRACE: 2 =======\nLine: %d\nFile: %s,\nMethod: %s",
				errF1.Line,
				errF1.File,
				errF1.Function,
			),
			err.Error(),
		),
	)
}
