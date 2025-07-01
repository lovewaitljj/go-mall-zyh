package dao

import (
	"context"
	"errors"
	"github.com/go-study-lab/go-mall/common/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 500 * time.Millisecond, // 超过500毫秒算慢查询, 如果需要按环境定制化, 这个放到配置文件中
	}
}

// LogMode log mode
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{}
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.Info(ctx, msg, "data", data)
}

// Warn print info
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.Warn(ctx, msg, "data", data)
}

// Error print info
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Error(ctx, msg, "data", data)
}

// Trace gorm的trace函数记录了sql以及错误，自定制往里加日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	duration := time.Since(begin).Milliseconds()
	// 获取 SQL 语句和返回条数
	sql, rows := fc()
	// Gorm 错误时记录错误日志
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(ctx, "SQL ERROR", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
	// 慢查询日志
	if duration > l.SlowThreshold.Milliseconds() {
		logger.Warn(ctx, "SQL SLOW", "sql", sql, "rows", rows, "dur(ms)", duration)
	} else {
		logger.Debug(ctx, "SQL DEBUG", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
}
