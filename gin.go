package logformatter

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type (
	// GinLogConfig represents the configuration to setup gin logger
	GinLogConfig struct {
		Severity Severity
	}
)

var defaultGinLogConfig = GinLogConfig{Severity: Debug}

// GinLogSeverity sets the severity for the GinLog
func GinLogSeverity(s Severity) func(*GinLogConfig) {
	return func(g *GinLogConfig) {
		g.Severity = s
	}
}

// NewGinLogFormatter takes zero or one GinOption function and applis to GinLog.
func NewGinLogFormatter(Options ...func(*GinLogConfig)) gin.LogFormatter {
	config := defaultGinLogConfig

	for _, o := range Options {
		o(&config)
	}

	return func(params gin.LogFormatterParams) string {
		sLog := StackdriverLog{
			HttpRequest: httpRequest{
				RequestMethod: params.Method,
				RequestUrl:    params.Request.URL.String(),
				Status:        params.StatusCode,
				UserAgent:     params.Request.UserAgent(),
				RemoteIp:      params.ClientIP,
				Latency:       fmt.Sprintf("%fs", params.Latency.Seconds()),
				Protocol:      params.Request.Proto,
				ResponseSize:  fmt.Sprintf("%d", params.BodySize),
			},
			Severity:  config.Severity.String(),
			Timestamp: params.TimeStamp.String(),
		}

		result, err := json.Marshal(sLog)

		if err != nil {
			return err.Error() + "\n"
		}

		return fmt.Sprintf("%s\n", string(result))
	}
}
