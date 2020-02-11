package logformatter

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/genproto/googleapis/devtools/clouderrorreporting/v1beta1"
	"google.golang.org/genproto/googleapis/logging/v2"
)

type Stackdriver struct {
	logging.LogEntry
	clouderrorreporting.ErrorEvent
	Payload map[string]interface{}
}

func (s *Stackdriver) MarshalJSON() ([]byte, error) {
	var (
		m jsonpb.Marshaler
	)

	// unwrap aims to remove the double quote after the component of stackdriver marshalled.
	unwrap := func(m []byte) []byte {
		if m == nil {
			return nil
		}
		return m[1 : len(m)-1]
	}

	buffer := bytes.NewBuffer([]byte("{"))

	entry, err := m.MarshalToString(&s.LogEntry)
	if err != nil {
		return nil, err
	}

	if !isEmptyJson(entry) {
		buffer.WriteString(string(unwrap([]byte(entry))))
	}

	report, err := m.MarshalToString(&s.ErrorEvent)
	if err != nil {
		return nil, err
	}

	if !isEmptyJson(report) {
		buffer.WriteString("," + string(unwrap([]byte(report))))
	}

	if s.Payload != nil {
		p, err := json.Marshal(s.Payload)
		if err != nil {
			return nil, err
		}

		if !isEmptyJson(string(p)) {
			buffer.WriteString(",\"payload\":" + string(p))
		}
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func isEmptyJson(s string) bool {
	return bytes.Equal([]byte(s), []byte("{}"))
}
