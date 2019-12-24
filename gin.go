package logformatter

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type (
	// GinLog represents the configuration to setup gin logger
	GinLog struct {
		Severity Severity
	}
)

var defaultGinLog = GinLog{Severity: Info}

// GinLogSeverity sets the severity for the GinLog
func GinLogSeverity(s Severity) func(*GinLog) {
	return func(g *GinLog) {
		g.Severity = s
	}
}

// NewGinLogFormatter takes zero or one GinOption function and applis to GinLog.
func NewGinLogFormatter(Options ...func(*GinLog)) gin.LogFormatter {
	config := defaultGinLog

	for _, o := range Options {
		o(&config)
	}

	return func(params gin.LogFormatterParams) string {
		sLog := StackdriverLog{
			HttpRequest: &httpRequest{
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
