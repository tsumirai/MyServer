package util

import (
	"MyServer/middleware/logger"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func getLocalIP() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Error(logger.LogArgs{"err": err.Error(), "msg": "获得ip地址失败"})
		return ip
	}

	for _, v := range addrs {
		if ipNet, ok := v.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ip
}

func GenerateTraceID() string {
	ip := getLocalIP()
	now := time.Now()
	timeStamp := now.Unix()
	nano := now.Nanosecond()
	pid := os.Getpid()
	traceID := bytes.Buffer{}
	traceID.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))
	traceID.WriteString(fmt.Sprintf("%4x", timeStamp))
	traceID.WriteString(fmt.Sprintf("%4x", nano))
	traceID.WriteString(fmt.Sprintf("%4x", pid))
	traceID.WriteString(fmt.Sprintf("%4x", rand.Int31n(1<<24)))
	return traceID.String()
}
