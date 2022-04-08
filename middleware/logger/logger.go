package logger

import (
	config "MyServer/conf"
	"bytes"
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
		// ticket
		ticket := c.Value("ticket")

		//log.Hooks.Add(NewContextHook())

		log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqURI,
			"ticket":       ticket,
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
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i < len(keys)-1 {
			b.WriteString(fmt.Sprintf("%+v=%+v||", k, a[k]))
		} else {
			b.WriteString(fmt.Sprintf("%+v=%+v", k, a[k]))
		}
	}
	return b.String()
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

func Infof(format string, args LogArgs) {
	log.Infof(format, args.String())
}

func Warnf(format string, args LogArgs) {
	log.Warnf(format, args.String())
}
func Debugf(format string, args LogArgs) {
	log.Debugf(format, args.String())
}
func Errorf(format string, args LogArgs) {
	log.Errorf(format, args.String())
}
func Tracef(format string, args LogArgs) {
	log.Tracef(format, args.String())
}

func Panicf(format string, args LogArgs) {
	log.Panicf(format, args.String())
}

func Printf(format string, args LogArgs) {
	log.Printf(format, args.String())
}

func Fatalf(format string, args LogArgs) {
	log.Fatalf(format, args.String())
}

func Trace(args LogArgs) {
	log.Log(logrus.TraceLevel, args.String())
}

func Debug(args LogArgs) {
	log.Log(logrus.DebugLevel, args.String())
}

func Info(args LogArgs) {
	log.Log(logrus.InfoLevel, args.String())
}

func Warn(args LogArgs) {
	log.Log(logrus.WarnLevel, args.String())
}

func Warning(args LogArgs) {
	log.Warn(args.String())
}

func Error(args LogArgs) {
	log.Log(logrus.ErrorLevel, args.String())
}

func Fatal(args LogArgs) {
	log.Log(logrus.FatalLevel, args.String())
	log.Exit(1)
}

func Panic(args LogArgs) {
	log.Log(logrus.PanicLevel, args.String())
}

func Traceln(args LogArgs) {
	log.Logln(logrus.TraceLevel, args.String())
}

func Debugln(args LogArgs) {
	log.Logln(logrus.DebugLevel, args.String())
}

func Infoln(args LogArgs) {
	log.Logln(logrus.InfoLevel, args.String())
}

func Warnln(args LogArgs) {
	log.Logln(logrus.WarnLevel, args.String())
}

func Warningln(args LogArgs) {
	log.Warnln(args.String())
}

func Errorln(args LogArgs) {
	log.Logln(logrus.ErrorLevel, args.String())
}

func Fatalln(args LogArgs) {
	log.Logln(logrus.FatalLevel, args.String())
	log.Exit(1)
}

func Panicln(args LogArgs) {
	log.Logln(logrus.PanicLevel, args.String())
}
