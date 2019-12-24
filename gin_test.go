package logformatter

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewGinLogFormatter(t *testing.T) {
	for _, c := range [...]struct {
		name   string
		expect StackdriverLog
	}{
		{
			name: "Test gin log contains required entries",
			expect: StackdriverLog{
				// Fill up default zero entries
				HttpRequest: &httpRequest{
					Protocol:     "HTTP/1.1",
					Latency:      fmt.Sprintf("%fs", float64(0)),
					RequestUrl:   "http://test.url/",
					ResponseSize: fmt.Sprintf("%d", 0),
				},
				Severity:  defaultGinLog.Severity.String(),
				Timestamp: time.Time{}.String(), // default empty value
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			formatter := NewGinLogFormatter()

			mockRequest := httptest.NewRequest("GET", "http://test.url/", nil)

			out := formatter(gin.LogFormatterParams{Request: mockRequest})

			expectJson, _ := json.Marshal(&c.expect)
			assert.JSONEq(t, out, string(expectJson))
		})
	}
}

func TestSetGinLogSeverity(t *testing.T) {
	cases := []struct {
		name   string
		option func(*GinLog)
		expect Severity
	}{
		{
			name:   "Test default severity(Info)",
			option: nil,
			expect: Info,
		},
		{
			name:   "Test set gin log severity",
			option: GinLogSeverity(Warning),
			expect: Warning,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var formatter gin.LogFormatter
			if c.option != nil {
				formatter = NewGinLogFormatter(c.option)
			} else {
				formatter = NewGinLogFormatter()
			}

			mockRequest := httptest.NewRequest("GET", "http://test.url/", nil)

			out := formatter(gin.LogFormatterParams{
				Request: mockRequest})
			assert.Contains(t, out, fmt.Sprintf(`"severity":"%s"`, c.expect.String()))
		})
	}
}
