package logformatter

import (
	"encoding/json"
)

type (
	// Define the object format of log for stackdriver
	StackdriverLog struct {
		HttpRequest *httpRequest `json:"httpRequest"`
		Severity    string       `json:"severity"`
		Timestamp   string       `json:"timestamp"`
	}

	httpRequest struct {
		RequestMethod string `json:"requestMethod"`
		RequestUrl    string `json:"requestUrl"`
		Status        int    `json:"status"`
		UserAgent     string `json:"userAgent"`
		RemoteIp      string `json:"remoteIp"`
		Latency       string `json:"latency"`
		Protocol      string `json:"protocol"`
		ResponseSize  string `json:"responseSize"`
	}

	Severity int
)

const (
	Default Severity = iota
	Debug
	Info
	Notice
	Warning
	Error
	Critical
	Alert
	Emergency
)

var toString = map[Severity]string{
	Default:   "DEFAULT",
	Debug:     "DEBUG",
	Info:      "INFO",
	Notice:    "NOTICE",
	Warning:   "WARNING",
	Error:     "ERROR",
	Critical:  "CRITICAL",
	Alert:     "ALERT",
	Emergency: "EMERGENCY",
}

var toEnum = map[string]Severity{
	"DEFAULT":   Default,
	"DEBUG":     Debug,
	"INFO":      Info,
	"NOTICE":    Notice,
	"WARNING":   Warning,
	"ERROR":     Error,
	"CRITICAL":  Critical,
	"ALERT":     Alert,
	"EMERGENCY": Emergency,
}

func (s Severity) String() string {
	return toString[s]
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return []byte(`"` + toString[s] + `"`), nil
}

func (s *Severity) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	*s = toEnum[str]
	return nil
}
