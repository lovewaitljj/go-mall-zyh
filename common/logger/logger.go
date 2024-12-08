package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

type logger struct {
	_logger *zap.Logger
}

var log *logger

// 起名a_zap.go是由于在一个包内，多个 Go 文件的init执行顺序按照 文件名的词典序执行。

func init() {
	log = &logger{
		_logger: _logger,
	}
}

// kv 应该是成对的数据, 类似: name,张三,age,10,...
func (l *logger) log(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	// 保证要打印的日志信息成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 日志行信息中增加追踪参数
	kv = append(kv, "traceid", ctx.Value("traceid"), "spanid", ctx.Value("spanid"), "pspanid", ctx.Value("pspanid"))
	// 增加日志调用者信息, 方便查日志时定位程序位置
	funcName, file, line := l.getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)

}

func Debug(ctx context.Context, msg string, kv ...interface{}) {
	log.Info(ctx, msg, kv...)
}

func Info(ctx context.Context, msg string, kv ...interface{}) {
	log.Info(ctx, msg, kv...)
}

func Warn(ctx context.Context, msg string, kv ...interface{}) {
	log.Info(ctx, msg, kv...)
}
func Error(ctx context.Context, msg string, kv ...interface{}) {
	log.Info(ctx, msg, kv...)
}

func (l *logger) Debug(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.DebugLevel, msg, kv...)
}

func (l *logger) Info(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.InfoLevel, msg, kv...)
}

func (l *logger) Warn(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.WarnLevel, msg, kv...)
}

func (l *logger) Error(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.ErrorLevel, msg, kv...)
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func (l *logger) getLoggerCallerInfo() (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(4) // 回溯3层，正好是调用日志的地方
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
