package logger

import (
	"MyServer/base"
	"bytes"
	"context"
	"fmt"
	"io"
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

// InitLogger 初始化日志配置
func (l *LogModel) InitLogger() {
	logFilePath := getCurrentPathByCaller() + base.Config.GetString("log.log_file_path")
	logFileName := base.Config.GetString("log.log_file_name")
	errLogFileName := base.Config.GetString("log.err_log_file_name")

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
	logLevel, err := logrus.ParseLevel(base.Config.GetString("log.log_level"))
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

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 日志中间件。
func (l *LogModel) LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用自定义 ResponseWriter
		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw

		// 打印请求信息
		reqBody, _ := c.GetRawData()
		// trace_id
		traceID := c.Value("trace_id")
		// 请求IP
		clientIP := c.ClientIP()

		log.WithFields(logrus.Fields{
			"req_method": c.Request.Method,
			"req_uri":    c.Request.RequestURI,
			"req_body":   reqBody,
			"trace_id":   traceID,
			"client_ip":  clientIP,
		}).Info()

		//fmt.Printf("[INFO] Request: %s %s %s\n", c.Request.Method, c.Request.RequestURI, reqBody)
		// 请求包体写回。
		if len(reqBody) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 开始时间
		startTime := time.Now()

		// 执行请求处理程序和其他中间件函数
		c.Next()

		respBody := string(crw.body.Bytes())
		// 状态码
		statusCode := c.Writer.Status()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//fmt.Printf("[INFO] Response: %s %s %s (%v)\n", c.Request.Method, c.Request.RequestURI, respBody, latency)

		log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   c.Request.Method,
			"req_uri":      c.Request.RequestURI,
			"resp_body":    respBody,
			"trace_id":     traceID,
		}).Info()
	}
}

// // LoggerToFile 将日志输出到文件
// func (l *LogModel) LoggerToFile() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 开始时间
// 		startTime := time.Now()
// 		// 处理请求
// 		c.Next()
// 		// 结束时间
// 		endTime := time.Now()
// 		// 执行时间
// 		latencyTime := endTime.Sub(startTime)
// 		// 请求方式
// 		reqMethod := c.Request.Method
// 		// 请求路由
// 		reqURI := c.Request.RequestURI
// 		// 状态码
// 		statusCode := c.Writer.Status()
// 		// 请求IP
// 		clientIP := c.ClientIP()
// 		// trace_id
// 		traceID := c.Value("trace_id")

// 		//log.Hooks.Add(NewContextHook())

// 		log.WithFields(logrus.Fields{
// 			"status_code":  statusCode,
// 			"latency_time": latencyTime,
// 			"client_ip":    clientIP,
// 			"req_method":   reqMethod,
// 			"req_uri":      reqURI,
// 			"trace_id":     traceID,
// 		}).Info()

// 		//日志格式
// 		//logger.Infof("| %3d | %l3v | %15s | %s | %s |",
// 		//	statusCode,
// 		//	latencyTime,
// 		//	clientIP,
// 		//	reqMethod,
// 		//	reqURI,
// 		//)
// 	}
// }

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

// formatMsg 格式化输出日志
func formatMsg(ctx context.Context, module string, args LogArgs) string {
	traceID := ""
	if ctx.Value("trace_id") != nil {
		traceID = ctx.Value("trace_id").(string)
	}

	msg := ""
	if traceID != "" {
		msg += "trace_id=" + traceID + "||"
	}

	if module != "" {
		msg += "module=" + module + "||"
	}

	msg += args.String()

	return msg
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

func Trace(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.TraceLevel, msg)
}

func Debug(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.DebugLevel, msg)
}

func Info(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.InfoLevel, msg)
}

func Warn(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.WarnLevel, msg)
}

func Error(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.ErrorLevel, msg)
}

func Fatal(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.FatalLevel, msg)
	log.Exit(1)
}

func Panic(ctx context.Context, module string, args LogArgs) {
	msg := formatMsg(ctx, module, args)
	log.Log(logrus.PanicLevel, msg)
}
