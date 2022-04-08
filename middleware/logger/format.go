package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (l *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry == nil {
		return make([]byte, 0), nil
	}

	file := ""
	ifile, exit := entry.Data["file"]
	if exit && ifile != nil {
		file = fmt.Sprintf("%s", ifile)
	}

	extMsg := ""
	itraceID, exit := entry.Data["trace_id"]
	if exit && itraceID != nil {
		if extMsg != "" {
			extMsg += "||"
		}
		extMsg += fmt.Sprintf("trace_id=%s", itraceID)
	}

	reqMethod, exit := entry.Data["req_method"]
	if exit && reqMethod != nil {
		if extMsg != "" {
			extMsg += "||"
		}
		extMsg += fmt.Sprintf("req_method=%s", reqMethod)
	}

	reqUri, exit := entry.Data["req_uri"]
	if exit && reqUri != nil {
		if extMsg != "" {
			extMsg += "||"
		}
		extMsg += fmt.Sprintf("req_uri=%s", reqUri)
	}

	statusCode, exit := entry.Data["status_code"]
	if exit && statusCode != nil {
		if extMsg != "" {
			extMsg += "||"
		}
		extMsg += fmt.Sprintf("status=%d", statusCode)
	}

	latencyTime, exit := entry.Data["latency_time"]
	if exit && latencyTime != nil {
		if extMsg != "" {
			extMsg += "||"
		}
		extMsg += fmt.Sprintf("latency_time=%s", latencyTime)
	}

	clientIp, exit := entry.Data["client_ip"]
	if exit && clientIp != nil {
		if extMsg != "" {
			extMsg += "||"
		}
		extMsg += fmt.Sprintf("client_ip=%s", clientIp)
	}

	if extMsg != "" {
		extMsg += "||"
	}

	msg := fmt.Sprintf("[%s][%s][%s]%s%s\n", entry.Level, entry.Time.Format("2006-01-02 15:04:05"), file, extMsg, entry.Message)
	return []byte(msg), nil
}
