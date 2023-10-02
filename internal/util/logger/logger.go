package logger

import (
	"context"
	"log"

	config "github.com/mochammadshenna/arch-pba-template/config"
	"github.com/mochammadshenna/arch-pba-template/internal/state"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

var lf loggerField

func Init() {
	lf = newLoggerField()

	logConfig := config.Get().Log

	Logger.SetFormatter(
		NewFormatter(
			WithStackSkip("github.com/mochammadshenna/arch-pba-template/internal/util/logger"),
			WithStackSkip("github.com/mochammadshenna/arch-pba-template/internal/util/helper"),
		),
	)

	if logConfig.Level == "" {
		logConfig.Level = "debug"
	}
	logLevel, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		log.Printf("error logger %v", err)
		panic(err)
	}
	Logger.SetLevel(logLevel)
}

type loggerField struct {
	// Custom field
	RequestId     string `json:"requestId"`
	RequestMethod string `json:"requestMethod"`
	Resource      string `json:"resource"`
	Status        string `json:"status"`

	// Field handle by logger
	Message        string `json:"message"`
	Severity       string `json:"severity"`
	Timestamp      string `json:"timestamp"`
	SourceLocation string `json:"sourceLocation"`
}

func newLoggerField() loggerField {
	return loggerField{
		RequestId:      "requestId",
		RequestMethod:  "requestMethod",
		Resource:       "resource",
		Status:         "status",
		Message:        "message",
		Severity:       "severity",
		Timestamp:      "timestamp",
		SourceLocation: "sourceLocation",
	}
}

func LoggerField() loggerField {
	return lf
}

func Trace(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Trace(args...)
		return
	}
	Logger.Trace(args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Tracef(format, args...)
		return
	}
	Logger.Tracef(format, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Debug(args...)
		return
	}
	Logger.Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Debugf(format, args...)
		return
	}
	Logger.Debugf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Info(args...)
		return
	}
	Logger.Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Infof(format, args...)
		return
	}
	Logger.Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Warn(args...)
		return
	}
	Logger.Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Warnf(format, args...)
		return
	}
	Logger.Warnf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Error(args...)
		return
	}
	Logger.Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Errorf(format, args...)
		return
	}
	Logger.Errorf(format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Fatal(args...)
		return
	}
	Logger.Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Fatalf(format, args...)
		return
	}
	Logger.Fatalf(format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Panic(args...)
		return
	}
	Logger.Panic(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	fields := withFields(ctx)
	if len(fields) > 0 {
		Logger.WithFields(fields).Panicf(format, args...)
		return
	}
	Logger.Panicf(format, args...)
}

func withFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}

	requestId := ctx.Value(state.HttpHeaders().RequestId)
	if requestId != "" {
		fields[LoggerField().RequestId] = requestId
	}

	return fields
}
