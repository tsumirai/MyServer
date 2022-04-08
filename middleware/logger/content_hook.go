package logger

import (
	"fmt"
	"runtime"
	"strings"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type contentHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
	output *rotatelogs.RotateLogs
}

func NewContentHook(output *rotatelogs.RotateLogs, levels ...logrus.Level) logrus.Hook {
	hook := contentHook{
		Field:  "file",
		Skip:   5,
		levels: levels,
		output: output,
	}

	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return &hook
}

func (hook contentHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook contentHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	content, err := entry.Bytes()
	if err != nil {
		panic(err)
	}
	hook.output.Write(content)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 20; i++ {
		file, line = getCaller(skip + i)
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
