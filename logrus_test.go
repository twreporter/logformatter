package logformatter

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/logging/type"
)

type frame struct {
	funcName string
	file     string
	line     int
}

func TestFormatStack(t *testing.T) {
	pwd, _ := os.Getwd()

	for _, c := range [...]struct {
		name   string
		err    error
		expect *frame
	}{
		{
			name: "Nil error object",
		},
		{
			name: "Error without stack information",
		},
		{
			name: "Test error with github.com/pkg/errors wrap object",
			err:  errors.Errorf("rrrr"),
			expect: &frame{
				funcName: "github.com/twreporter/logformatter.TestFormatStack",
				file:     pwd + "/logrus_test.go",
				line:     40,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			out := FormatStack(c.err)

			if c.expect == nil {
				assert.Equal(t, string(out), "")
			} else {
				// Validate if the format adheres to runtime.Stack
				rs := out
				goroutineState := rs[:bytes.Index(rs, []byte("\n"))]
				// Validate first line contains goroutine state information
				assert.Regexp(t, regexp.MustCompile(`goroutine [1-9]+ \[running\]:`), string(goroutineState))

				// Next, validate if the stack frame format follow
				// [function]\n\t[file]:[line] +[function address]
				segments := bytes.Split(rs[bytes.Index(rs, []byte("\n"))+1:], []byte("\n"))
				firstStack := []byte(string(segments[0]) + "\n" + string(segments[1]))
				s := bytes.Split(firstStack, []byte(" "))

				assert.Equal(t, fmt.Sprintf("%s()\n\t%s:%d", c.expect.funcName, c.expect.file, c.expect.line), string(s[0]))

				assert.Equal(t, strings.HasPrefix(string(s[1]), "+0x"), true)
				_, err := strconv.ParseInt(string(s[1][3:]), 16, 64)
				assert.Nil(t, err)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	for _, c := range [...]struct {
		name           string
		entry          logrus.Entry
		expectSeverity ltype.LogSeverity
		expectMessage  string
	}{
		{
			name: "Test info log",
			entry: logrus.Entry{
				Level:   logrus.InfoLevel,
				Message: "Mock info log",
			},
			expectSeverity: ltype.LogSeverity_INFO,
			expectMessage:  "Mock info log",
		},
		{
			name: "Test error log",
			entry: logrus.Entry{
				Level:   logrus.ErrorLevel,
				Message: "Mock error log",
			},
			expectSeverity: ltype.LogSeverity_ERROR,
			expectMessage:  "Mock error log",
		},
	} {
		f := NewStackdriverFormatter("Mock service", "test")
		t.Run(c.name, func(t *testing.T) {
			out, _ := f.Format(&c.entry)
			assert.Contains(t, string(out), fmt.Sprintf(`"severity":%d`, c.expectSeverity))
			assert.Contains(t, string(out), fmt.Sprintf(`"message":"%s"`, c.expectMessage))
		})
	}
}
