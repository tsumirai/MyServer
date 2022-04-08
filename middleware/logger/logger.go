package logger

import (
	config "MyServer/conf"
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type LogModel struct {
	logModel *logrus.Logger
}

func NewLogModel() *LogModel {
	return &LogModel{
		logModel: log,
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		fmt.Println("PathExists: not exist")
		return false, nil
	}
	return false, err
}

func getCurrentPathByCaller() string {
	var execPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		execPath = filepath.Dir(filename)
	}
	return execPath
}

func (l *LogModel) InitLogger() {
	logFilePath := getCurrentPathByCaller() + config.Config.GetString("log.log_file_path")
	logFileName := config.Config.GetString("log.log_file_name")
	errLogFileName := config.Config.GetString("log.err_log_file_name")

	exist, err := pathExists(logFilePath)
	if err != nil {
		fmt.Println("InitLogger Failed: ", err.Error())
		panic(err)
	}

	if !exist {
		err = os.Mkdir(logFilePath, os.ModePerm)
		if err != nil {
			fmt.Println("InitLogger Failed: mkdir failed! ", err.Error())
			panic(err)
		} else {
			fmt.Println("InitLogger mkdir success!")
		}
	} else {
		fmt.Println("InitLogger path exist!")
	}

	fileName := path.Join(logFilePath, logFileName)
	errFileName := path.Join(logFilePath, errLogFileName)

	logContent, err := rotatelogs.New(
		fileName+"-%Y-%m-%d-%H",
		rotatelogs.WithLinkName(fileName),
		// MaxAge and RotationCount cannot be both set 两者不能同时设置
		//rotatelogs.WithMaxAge(5*time.Minute),
		rotatelogs.WithRotationCount(5),        // number 默认7份，大于7份或者到了清理时间，开始清理
		rotatelogs.WithRotationTime(time.Hour), // rotate最小为1分钟轮询。默认60s，低于1分钟按照1分钟来
	)
	if err != nil {
		fmt.Println("InitLogger Failed: ", err.Error())
		panic(err)
	}

	errLogContent, err := rotatelogs.New(
		errFileName+"-%Y-%m-%d-%H",
		rotatelogs.WithLinkName(errFileName),
		// MaxAge and RotationCount cannot be both set 两者不能同时设置
		//rotatelogs.WithMaxAge(5*time.Minute),
		rotatelogs.WithRotationCount(5),        // number 默认7份，大于7份或者到了清理时间，开始清理
		rotatelogs.WithRotationTime(time.Hour), // rotate最小为1分钟轮询。默认60s，低于1分钟按照1分钟来
	)
	if err != nil {
		fmt.Println("InitErrLogger Failed: ", err.Error())
		panic(err)
	}

	//src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm)
	//if err != nil {
	//	fmt.Println("LoggerToFile Open file failed: ", err.Error())
	//	panic(err)
	//}

	//logger := logrus.New()
	// 设置输出
	// log.SetOutput(logContent)

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(config.Config.GetString("log.log_level"))
	if err != nil {
		fmt.Println("LoggerToFile Parse logLevel failed: ", err.Error())
		panic(err)
	}
	log.SetLevel(logLevel)
	log.Hooks.Add(NewContentHook(logContent, logrus.InfoLevel))
	log.Hooks.Add(NewContentHook(errLogContent, logrus.ErrorLevel, logrus.PanicLevel, logrus.WarnLevel, logrus.FatalLevel))

	//log.SetOutput(logContent)
	// 设置行号和文件名
	//log.SetReportCaller(true)

	// 设置日志格式
	// log.SetFormatter(&logrus.JSONFormatter{
	// 	//ForceColors:     true,
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// 	//CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
	// 	//	// 处理文件名
	// 	//	fileShortName := path.Base(frame.File)
	// 	//	return frame.Function + " : " + strconv.Itoa(frame.Line), fileShortName
	// 	//},
	// })
	log.SetFormatter(new(LogFormatter))
	log.Info("Init Log Success")
}

// 将日志输出到文件
func (l *LogModel) LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqURI := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// trace_id
		traceID := c.Value("trace_id")

		//log.Hooks.Add(NewContextHook())

		log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqURI,
			"trace_id":     traceID,
		}).Info()

		//日志格式
		//logger.Infof("| %3d | %l3v | %15s | %s | %s |",
		//	statusCode,
		//	latencyTime,
		//	clientIP,
		//	reqMethod,
		//	reqURI,
		//)
	}
}

type LogArgs map[string]interface{}

func (a LogArgs) String() string {
	b := bytes.Buffer{}
	keys := make([]string, 0, len(a))
	for k := range a {
		if k != "trace_id" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	newKeys := make([]string, 0, len(a))
	if _, exit := a["trace_id"]; exit {
		newKeys = append(newKeys, "trace_id")
	}
	newKeys = append(newKeys, keys...)
	for i, k := range newKeys {
		if i < len(newKeys)-1 {
			b.WriteString(fmt.Sprintf("%+v=%+v||", k, a[k]))
		} else {
			b.WriteString(fmt.Sprintf("%+v=%+v", k, a[k]))
		}
	}
	return b.String()
}

func (a LogArgs) addTraceID(ctx context.Context) {
	traceID := ctx.Value("trace_id")
	if traceID != nil {
		a["trace_id"] = traceID
	}
}

// todo 日志记录到 MongoDB
func (l *LogModel) LoggerToMongoDB() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// todo 日志记录到ES
func (l *LogModel) LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// TODO 日志记录到MQ
func (l *LogModel) LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func Infof(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Infof(format, args.String())
}

func Warnf(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Warnf(format, args.String())
}
func Debugf(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Debugf(format, args.String())
}
func Errorf(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Errorf(format, args.String())
}
func Tracef(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Tracef(format, args.String())
}

func Panicf(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Panicf(format, args.String())
}

func Printf(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Printf(format, args.String())
}

func Fatalf(ctx context.Context, format string, args LogArgs) {
	args.addTraceID(ctx)
	log.Fatalf(format, args.String())
}

func Trace(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.TraceLevel, args.String())
}

func Debug(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.DebugLevel, args.String())
}

func Info(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.InfoLevel, args.String())
}

func Warn(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.WarnLevel, args.String())
}

func Warning(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Warn(args.String())
}

func Error(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.ErrorLevel, args.String())
}

func Fatal(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.FatalLevel, args.String())
	log.Exit(1)
}

func Panic(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Log(logrus.PanicLevel, args.String())
}

func Traceln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.TraceLevel, args.String())
}

func Debugln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.DebugLevel, args.String())
}

func Infoln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.InfoLevel, args.String())
}

func Warnln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.WarnLevel, args.String())
}

func Warningln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Warnln(args.String())
}

func Errorln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.ErrorLevel, args.String())
}

func Fatalln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.FatalLevel, args.String())
	log.Exit(1)
}

func Panicln(ctx context.Context, args LogArgs) {
	args.addTraceID(ctx)
	log.Logln(logrus.PanicLevel, args.String())
}
