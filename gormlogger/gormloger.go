package gormlogger

import (
	"context"
	"errors"
	"time"

	mylog "2026-FM247-BackEnd/logger"

	gormlog "gorm.io/gorm/logger"
)

type StdLogger struct {
	LogLevel gormlog.LogLevel
}

func NewStdLogger(level gormlog.LogLevel) *StdLogger {
	return &StdLogger{LogLevel: level}
}

func (l *StdLogger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *StdLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlog.Info {
		mylog.Log.Infof(msg, data...)
	}
}

func (l *StdLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlog.Warn {
		mylog.Log.Warnf(msg, data...)
	}
}

func (l *StdLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlog.Error {
		mylog.Log.Errorf(msg, data...)
	}
}

func (l *StdLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlog.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && !errors.Is(err, gormlog.ErrRecordNotFound) {
		mylog.Log.Errorf("数据库错误 | 耗时: %v | 行数: %d | SQL: %s | 错误: %v",
			elapsed, rows, sql, err)
	} else if elapsed > 200*time.Millisecond {
		mylog.Log.Warnf("慢查询 | 耗时: %v | 行数: %d | SQL: %s",
			elapsed, rows, sql)
	} else if l.LogLevel >= gormlog.Info {
		mylog.Log.Infof("数据库查询 | 耗时: %v | 行数: %d | SQL: %s",
			elapsed, rows, sql)
	}
}
