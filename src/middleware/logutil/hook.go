package logutil

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type contextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

func NewContextHook(levels ...logrus.Level) logrus.Hook {
	hook := contextHook{
		Field:  "line",
		Skip:   0,
		levels: levels,
	}

	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return &hook
}

func (hook contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook contextHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		fmt.Println(file)
		if !strings.HasPrefix(file, "logrus") && !strings.HasPrefix(file, "logutil") && !strings.HasPrefix(file, "gin") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
