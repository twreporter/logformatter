package logformatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/devtools/clouderrorreporting/v1beta1"
	"google.golang.org/genproto/googleapis/logging/type"
)

type (
	// Interface for inspecting error objects recursively
	// Ref: https://github.com/pkg/errors/blob/7f95ac13edff643b8ce5398b6ccab125f8a20c1a/errors.go#L30-L53
	causer interface {
		Cause() error
	}

	// Interface for retrieving stack frame for each error object
	// Ref: https://github.com/pkg/errors/blob/7f95ac13edff643b8ce5398b6ccab125f8a20c1a/errors.go#L68-L73
	stackTracer interface {
		StackTrace() errors.StackTrace
	}
)

// Adapt from https://github.com/googleapis/google-cloud-go/issues/1084#issuecomment-474565019
func FormatStack(err error) (buffer []byte) {
	if err == nil {
		return nil
	}

	// find the inner most error with a stack
	inner := err
	for inner != nil {
		if cause, ok := inner.(causer); ok {
			inner = cause.Cause()
			if _, ok := inner.(stackTracer); ok {
				err = inner
			}
		} else {
			break
		}
	}

	if stackTrace, ok := err.(stackTracer); ok {
		buf := bytes.Buffer{}
		buf.WriteString(getGoroutineState() + "\n")

		// format each frame of the stack to match runtime.Stack's format
		var lines []string
		for _, frame := range stackTrace.StackTrace() {
			pc := uintptr(frame) - 1
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				file, line := fn.FileLine(pc)
				lines = append(lines, fmt.Sprintf("%s()\n\t%s:%d +%#x", fn.Name(), file, line, fn.Entry()))
			}
		}
		buf.WriteString(strings.Join(lines, "\n"))

		buffer = buf.Bytes()
	}
	return
}

func NewStackdriverFormatter(service, version string) *Stackdriver {
	return &Stackdriver{
		ErrorEvent: clouderrorreporting.ErrorEvent{
			ServiceContext: &clouderrorreporting.ServiceContext{
				Service: service,
				Version: version,
			},
		},
	}
}

func (s *Stackdriver) Format(entry *logrus.Entry) ([]byte, error) {
	// Copy customized fields
	s.Payload = make(logrus.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			s.Payload[k] = v.Error()
		default:
			s.Payload[k] = v
		}
	}

	s.Message = entry.Message
	s.Severity = convertLevelToLogSeverity(entry.Level)

	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = new(bytes.Buffer)
	}

	encoder := json.NewEncoder(b)

	if err := encoder.Encode(s); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %+v", err)
	}

	return b.Bytes(), nil
}

func convertLevelToLogSeverity(lvl logrus.Level) ltype.LogSeverity {
	switch lvl {
	case logrus.InfoLevel:
		return ltype.LogSeverity_INFO
	case logrus.DebugLevel:
		return ltype.LogSeverity_DEBUG
	default:
		// Omit intentionally
	}
	return ltype.LogSeverity_ERROR
}

func getGoroutineState() string {
	stack := make([]byte, 64)
	stack = stack[:runtime.Stack(stack, false)]
	stack = stack[:bytes.Index(stack, []byte("\n"))]

	return string(stack)
}
