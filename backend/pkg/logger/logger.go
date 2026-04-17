package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"webGL-720yun/config"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	log      = logrus.New()
	logMutex sync.RWMutex
	logChan  = make(chan *logrus.Entry, 100)
)

func init() {
	go processLogEntries()
}

func processLogEntries() {
	for entry := range logChan {
		logMutex.RLock()
		entry.Info()
		logMutex.RUnlock()
	}
}

// GinZapLogger 自定义日志中间件
func GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		contentType := c.GetHeader("Content-Type")

		isUpload := strings.HasPrefix(path, "/api/v1/upload/chunk") || contentType == "multipart/form-data"

		var requestBody string
		if method != "GET" && method != "HEAD" && method != "OPTIONS" && !isUpload {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil && len(body) > 0 {
				requestBody = string(body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		fields := logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format("2006-01-02 15:04:05"),
			"latency":    latency,
		}

		if !isUpload {
			fields["body"] = requestBody
		}

		entry := log.WithFields(fields)

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}

		logChan <- entry
	}
}

// Setup 初始化日志系统
func Setup(logConf *config.LoggerConfig) {
	level, err := logrus.ParseLevel(logConf.Level)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.SetOutput(io.MultiWriter(
		&lumberjack.Logger{
			Filename:   logConf.Filename,
			MaxSize:    logConf.MaxSize, // megabytes
			MaxBackups: logConf.MaxBackups,
			MaxAge:     logConf.MaxAge,   // days
			Compress:   logConf.Compress, // disabled by default
		},
		os.Stdout, // 输出到控制台
	))

	log.SetReportCaller(true)
}

// Debug 在标准记录器上记录级别为Debug的消息
func Debug(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Debug(args...)
}

// Info 在标准记录器上记录级别为Info的消息
func Info(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Info(args...)
}

// Warn 在标准记录器上记录级别为Warn的消息
func Warn(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Warn(args...)
}

// Error 在标准记录器上记录级别为Error的消息
func Error(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Error(args...)
}

// Fatal 在标准记录器上记录级别为Fatal的消息，然后进程将以状态1退出
func Fatal(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Fatal(args...)
}

// Panic 在标准记录器上记录级别为Panic的消息
func Panic(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Panic(args...)
}

// Debugf 在标准记录器上记录级别为Debug的消息，支持格式化
func Debugf(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Debugf(format, args...)
}

// Infof 在标准记录器上记录级别为Info的消息，支持格式化
func Infof(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Infof(format, args...)
}

// Warnf 在标准记录器上记录级别为Warn的消息，支持格式化
func Warnf(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Warnf(format, args...)
}

// Errorf 在标准记录器上记录级别为Error的消息，支持格式化
func Errorf(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Errorf(format, args...)
}

// Fatalf 在标准记录器上记录级别为Fatal的消息，然后进程将以状态1退出，支持格式化
func Fatalf(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Fatalf(format, args...)
}

// Panicf 在标准记录器上记录级别为Panic的消息，支持格式化
func Panicf(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Panicf(format, args...)
}

// Debugln 在标准记录器上记录级别为Debug的消息，结尾添加换行符
func Debugln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Debugln(args...)
}

// Infoln 在标准记录器上记录级别为Info的消息，结尾添加换行符
func Infoln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Infoln(args...)
}

// Warnln 在标准记录器上记录级别为Warn的消息，结尾添加换行符
func Warnln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Warnln(args...)
}

// Errorln 在标准记录器上记录级别为Error的消息，结尾添加换行符
func Errorln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Errorln(args...)
}

// Fatalln 在标准记录器上记录级别为Fatal的消息，然后进程将以状态1退出，结尾添加换行符
func Fatalln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Fatalln(args...)
}

// Panicln 在标准记录器上记录级别为Panic的消息，结尾添加换行符
func Panicln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Panicln(args...)
}

// Trace 在标准记录器上记录级别为Trace的消息
func Trace(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Trace(args...)
}

// Tracef 在标准记录器上记录级别为Trace的消息，支持格式化
func Tracef(format string, args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Tracef(format, args...)
}

// Traceln 在标准记录器上记录级别为Trace的消息，结尾添加换行符
func Traceln(args ...interface{}) {
	logMutex.RLock()
	defer logMutex.RUnlock()
	log.Traceln(args...)
}

// WithField 返回一个带有指定字段的新条目
func WithField(key string, value interface{}) *logrus.Entry {
	logMutex.RLock()
	defer logMutex.RUnlock()
	return log.WithField(key, value)
}

// WithFields 返回一个带有指定字段的新条目
func WithFields(fields logrus.Fields) *logrus.Entry {
	logMutex.RLock()
	defer logMutex.RUnlock()
	return log.WithFields(fields)
}

// WithError 返回一个带有指定错误的新条目
func WithError(err error) *logrus.Entry {
	logMutex.RLock()
	defer logMutex.RUnlock()
	return log.WithError(err)
}
