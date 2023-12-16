package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	ctx      context.Context
	traceKey string
	logger   *logrus.Logger
}

func NewLogger(traceKey *string) *Logger {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)

	key := "X-Trace-Id"
	if traceKey != nil {
		key = *traceKey
	}
	return &Logger{
		ctx:      context.Background(),
		traceKey: key,
		logger:   logrus.New(),
	}
}

func (logger *Logger) withTraceKey() *logrus.Entry {
	if logger.traceKey != "" {
		if traceId, ok := logger.ctx.Value(logger.traceKey).(string); ok {
			return logger.logger.WithField(logger.traceKey, traceId)
		}
	}
	return logger.logger.WithContext(logger.ctx)
}

func (logger *Logger) WithTraceId(ctx context.Context, traceId string) {
	logger.ctx = context.WithValue(ctx, logger.traceKey, traceId)
}

func (logger *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return logger.withTraceKey().WithField(key, value)
}

func (logger *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.withTraceKey().WithFields(fields)
}

func (logger *Logger) WithError(err error) *logrus.Entry {
	return logger.withTraceKey().WithError(err)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.withTraceKey().Debugf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.withTraceKey().Infof(format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.withTraceKey().Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.withTraceKey().Warnf(format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.withTraceKey().Warningf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.withTraceKey().Errorf(format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.withTraceKey().Fatalf(format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.withTraceKey().Panicf(format, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.withTraceKey().Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.withTraceKey().Info(args...)
}

func (logger *Logger) Print(args ...interface{}) {
	logger.withTraceKey().Print(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.withTraceKey().Warn(args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.withTraceKey().Warning(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.withTraceKey().Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.withTraceKey().Fatal(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.withTraceKey().Panic(args...)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.withTraceKey().Debugln(args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.withTraceKey().Infoln(args...)
}

func (logger *Logger) Println(args ...interface{}) {
	logger.withTraceKey().Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.withTraceKey().Warnln(args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	logger.withTraceKey().Warningln(args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.withTraceKey().Errorln(args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.withTraceKey().Fatalln(args...)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.withTraceKey().Panicln(args...)
}
