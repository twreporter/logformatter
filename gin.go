package logformatter

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/genproto/googleapis/logging/v2"
)

type (
	// GinLog represents the configuration to setup gin logger
	GinLog struct {
		Severity ltype.LogSeverity
	}
)

var defaultGinLog = GinLog{Severity: ltype.LogSeverity_INFO}

// GinLogSeverity sets the severity for the GinLog
func GinLogSeverity(s ltype.LogSeverity) func(*GinLog) {
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
		sLog := logging.LogEntry{
			HttpRequest: &ltype.HttpRequest{
				RequestMethod: params.Method,
				RequestUrl:    params.Request.URL.String(),
				Status:        int32(params.StatusCode),
				UserAgent:     params.Request.UserAgent(),
				RemoteIp:      params.ClientIP,
				Protocol:      params.Request.Proto,
				ResponseSize:  int64(params.BodySize),
			},
			Severity: config.Severity,
		}

		if !params.TimeStamp.IsZero() {
			ts, _ := ptypes.TimestampProto(params.TimeStamp)
			sLog.Timestamp = ts
		}

		if int64(params.Latency) != 0 {
			sLog.HttpRequest.Latency = ptypes.DurationProto(params.Latency)
		}

		result, err := json.Marshal(sLog)

		if err != nil {
			return err.Error() + "\n"
		}

		return string(result) + "\n"
	}
}
