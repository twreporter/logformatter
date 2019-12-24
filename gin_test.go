package logformatter

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinLogFormatter(t *testing.T) {
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

			mockRequest, _ := http.NewRequest("GET", "http://test.url/", nil)

			outFormatString := formatter(gin.LogFormatterParams{
				Request: mockRequest})
			assert.Contains(t, outFormatString, fmt.Sprintf(`"severity":"%s"`, c.expect.String()))
		})
	}
}
