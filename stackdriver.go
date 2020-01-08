package logformatter

import (
	"google.golang.org/genproto/googleapis/devtools/clouderrorreporting/v1beta1"
	"google.golang.org/genproto/googleapis/logging/v2"
)

type Stackdriver struct {
	logging.LogEntry
	clouderrorreporting.ErrorEvent
	Payload map[string]interface{}
}
