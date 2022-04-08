package logutil

import (
	"MyServer/src/config"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func PathExists(path string) (bool, error) {
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

func InitLogger() {
	logFilePath := getCurrentPathByCaller() + config.Config.GetString("log.log_file_path")
	logFileName := config.Config.GetString("log.log_file_name")

	exist, err := PathExists(logFilePath)
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

	//src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm)
	//if err != nil {
	//	fmt.Println("LoggerToFile Open file failed: ", err.Error())
	//	panic(err)
	//}

	//logger := logrus.New()
	// 设置输出
	log.SetOutput(logContent)

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(config.Config.GetString("log.log_level"))
	if err != nil {
		fmt.Println("LoggerToFile Parse logLevel failed: ", err.Error())
		panic(err)
	}
	//log.SetLevel(logLevel)
	log.Hooks.Add(NewContextHook(logLevel))

	// 设置行号和文件名
	//log.SetReportCaller(true)

	// 设置日志格式
	log.SetFormatter(&logrus.JSONFormatter{
		//ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		//CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		//	// 处理文件名
		//	fileShortName := path.Base(frame.File)
		//	return frame.Function + " : " + strconv.Itoa(frame.Line), fileShortName
		//},
	})
	log.Info("Init Log Success")
}

// 将日志输出到文件
func LoggerToFile() gin.HandlerFunc {
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

		//log.Hooks.Add(NewContextHook())

		log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqURI,
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

// todo 日志记录到 MongoDB
func LoggerToMongoDB() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// todo 日志记录到ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// TODO 日志记录到MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args)
}
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}
func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func Trace(args ...interface{}) {
	log.Log(logrus.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	log.Log(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	log.Log(logrus.InfoLevel, args...)
}

func Warn(args ...interface{}) {
	log.Log(logrus.WarnLevel, args...)
}

func Warning(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Log(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	log.Log(logrus.FatalLevel, args...)
	log.Exit(1)
}

func Panic(args ...interface{}) {
	log.Log(logrus.PanicLevel, args...)
}

func Traceln(args ...interface{}) {
	log.Logln(logrus.TraceLevel, args...)
}

func Debugln(args ...interface{}) {
	log.Logln(logrus.DebugLevel, args...)
}

func Infoln(args ...interface{}) {
	log.Logln(logrus.InfoLevel, args...)
}

func Warnln(args ...interface{}) {
	log.Logln(logrus.WarnLevel, args...)
}

func Warningln(args ...interface{}) {
	log.Warnln(args...)
}

func Errorln(args ...interface{}) {
	log.Logln(logrus.ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	log.Logln(logrus.FatalLevel, args...)
	log.Exit(1)
}

func Panicln(args ...interface{}) {
	log.Logln(logrus.PanicLevel, args...)
}
