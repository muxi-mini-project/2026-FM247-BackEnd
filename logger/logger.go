package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// 定义日志等级
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger struct {
	level  Level
	logger *log.Logger
}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(levelStr string) Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN", "WARNING":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

func NewLogger(level Level, out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(out, prefix, flag),
	}
}

// 检查日志等级是否足够重要
func (l *Logger) shouldLog(level Level) bool {
	return level >= l.level
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.shouldLog(DebugLevel) {
		l.logger.Output(2, "[DEBUG] "+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.shouldLog(InfoLevel) {
		l.logger.Output(2, "[INFO] "+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.shouldLog(WarnLevel) {
		l.logger.Output(2, "[WARN] "+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.shouldLog(ErrorLevel) {
		l.logger.Output(2, "[ERROR] "+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Output(2, "[FATAL] "+fmt.Sprintf(format, v...))
	os.Exit(1)
}

var Log *Logger

func InitLogger(levelStr string) {
	// 解析日志级别
	level := ParseLevel(levelStr)

	// 设置输出
	//后续加写入文件功能
	multiWriter := io.MultiWriter(os.Stdout)

	Log = NewLogger(level, multiWriter, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}
