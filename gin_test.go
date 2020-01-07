package logformatter

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/genproto/googleapis/logging/v2"
)

func TestNewGinLogFormatter(t *testing.T) {
	for _, c := range [...]struct {
		name   string
		expect logging.LogEntry
	}{
		{
			name: "Test gin log contains required entries",
			expect: logging.LogEntry{
				// Fill up default zero entries
				HttpRequest: &ltype.HttpRequest{
					Protocol:   "HTTP/1.1",
					RequestUrl: "http://test.url/",
				},
				Severity: defaultGinLog.Severity,
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
		expect ltype.LogSeverity
	}{
		{
			name:   "Test default severity(Info)",
			option: nil,
			expect: ltype.LogSeverity_INFO,
		},
		{
			name:   "Test set gin log severity",
			option: GinLogSeverity(ltype.LogSeverity_WARNING),
			expect: ltype.LogSeverity_WARNING,
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
			assert.Contains(t, out, fmt.Sprintf(`"severity":%d`, c.expect))
		})
	}
}
