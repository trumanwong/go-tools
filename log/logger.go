package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

// Logger is a struct that encapsulates the logging functionality.
// It uses the logrus library for logging and allows for context-based logging.
// The context is used to store and retrieve values across API boundaries and between processes.
// The traceKey is a string that is used as a key to retrieve the trace ID from the context.
// The logger is an instance of the logrus.Logger struct, which is used to perform the actual logging.
type Logger struct {
	ctx      context.Context // The context in which the logger operates.
	traceKey string          // The key used to retrieve the trace ID from the context.
	logger   *logrus.Logger  // The underlying logrus logger.
}

// NewLogger is a function that creates a new instance of the Logger struct.
// It takes an optional traceKey as an argument. If no traceKey is provided, it defaults to "X-Trace-Id".
// The function configures the logrus library to log as JSON and output to stdout.
// It then creates a new Logger instance with a background context, the provided or default traceKey, and a new logrus logger.
// The newly created Logger instance is then returned.
func NewLogger(traceKey *string, formatter logrus.Formatter) *Logger {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if formatter != nil {
		logrus.SetFormatter(formatter)
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Default traceKey is "X-Trace-Id". If a traceKey is provided, use it instead.
	key := "X-Trace-Id"
	if traceKey != nil {
		key = *traceKey
	}

	// Create a new Logger instance with a background context, the provided or default traceKey, and a new logrus logger.
	return &Logger{
		ctx:      context.Background(),
		traceKey: key,
		logger:   logrus.New(),
	}
}

// withTraceKey is a method on the Logger struct.
// It checks if the traceKey of the logger is not an empty string.
// If the traceKey is not empty, it retrieves the traceId from the logger's context using the traceKey.
// If the traceId is successfully retrieved, it returns a new logrus Entry with the traceKey and traceId as a field.
// If the traceKey is empty or the traceId could not be retrieved, it returns a new logrus Entry with the logger's context.
func (logger *Logger) withTraceKey() *logrus.Entry {
	if logger.traceKey != "" {
		if traceId, ok := logger.ctx.Value(logger.traceKey).(string); ok {
			return logger.logger.WithField(logger.traceKey, traceId)
		}
	}
	return logger.logger.WithContext(logger.ctx)
}

// WithTraceId is a method on the Logger struct.
// It takes a context and a traceId as arguments.
// The method sets the logger's context to a new context with the logger's traceKey and the provided traceId as a value.
func (logger *Logger) WithTraceId(ctx context.Context, traceId string) {
	logger.ctx = context.WithValue(ctx, logger.traceKey, traceId)
}

// WithField is a method on the Logger struct.
// It takes a key and a value as arguments.
// The method returns a new logrus Entry with the provided key and value as a field.
// The Entry also includes the logger's traceKey and traceId if they are available.
func (logger *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return logger.withTraceKey().WithField(key, value)
}

// WithFields is a method on the Logger struct.
// It takes a map of keys and values (logrus.Fields) as an argument.
// The method returns a new logrus Entry with the provided fields.
// The Entry also includes the logger's traceKey and traceId if they are available.
func (logger *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.withTraceKey().WithFields(fields)
}

// WithError is a method on the Logger struct.
// It takes an error as an argument.
// The method returns a new logrus Entry with the provided error.
// The Entry also includes the logger's traceKey and traceId if they are available.
func (logger *Logger) WithError(err error) *logrus.Entry {
	return logger.withTraceKey().WithError(err)
}

// Debugf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs a debug message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.withTraceKey().Debugf(format, args...)
}

// Infof is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs an info message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.withTraceKey().Infof(format, args...)
}

// Printf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs a message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.withTraceKey().Printf(format, args...)
}

// Warnf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs a warning message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.withTraceKey().Warnf(format, args...)
}

// Warningf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs a warning message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.withTraceKey().Warningf(format, args...)
}

// Errorf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs an error message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.withTraceKey().Errorf(format, args...)
}

// Fatalf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs a fatal message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
// After logging the message, the method calls os.Exit(1) to terminate the program.
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.withTraceKey().Fatalf(format, args...)
}

// Panicf is a method on the Logger struct.
// It takes a format string and a variadic number of arguments.
// The method logs a panic message with the provided format and arguments.
// The message includes the logger's traceKey and traceId if they are available.
// After logging the message, the method calls panic() to stop the ordinary flow of a goroutine.
func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.withTraceKey().Panicf(format, args...)
}

// Debug is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a debug message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Debug(args ...interface{}) {
	logger.withTraceKey().Debug(args...)
}

// Info is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs an info message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Info(args ...interface{}) {
	logger.withTraceKey().Info(args...)
}

// Print is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Print(args ...interface{}) {
	logger.withTraceKey().Print(args...)
}

// Warn is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a warning message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Warn(args ...interface{}) {
	logger.withTraceKey().Warn(args...)
}

// Warning is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a warning message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Warning(args ...interface{}) {
	logger.withTraceKey().Warning(args...)
}

// Error is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs an error message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Error(args ...interface{}) {
	logger.withTraceKey().Error(args...)
}

// Fatal is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a fatal message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
// After logging the message, the method calls os.Exit(1) to terminate the program.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.withTraceKey().Fatal(args...)
}

// Panic is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a panic message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
// After logging the message, the method calls panic() to stop the ordinary flow of a goroutine.
func (logger *Logger) Panic(args ...interface{}) {
	logger.withTraceKey().Panic(args...)
}

// Debugln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a debug message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Debugln(args ...interface{}) {
	logger.withTraceKey().Debugln(args...)
}

// Infoln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs an info message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Infoln(args ...interface{}) {
	logger.withTraceKey().Infoln(args...)
}

// Println is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Println(args ...interface{}) {
	logger.withTraceKey().Println(args...)
}

// Warnln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a warning message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Warnln(args ...interface{}) {
	logger.withTraceKey().Warnln(args...)
}

// Warningln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a warning message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Warningln(args ...interface{}) {
	logger.withTraceKey().Warningln(args...)
}

// Errorln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs an error message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
func (logger *Logger) Errorln(args ...interface{}) {
	logger.withTraceKey().Errorln(args...)
}

// Fatalln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a fatal message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
// After logging the message, the method calls os.Exit(1) to terminate the program.
func (logger *Logger) Fatalln(args ...interface{}) {
	logger.withTraceKey().Fatalln(args...)
}

// Panicln is a method on the Logger struct.
// It takes a variadic number of arguments.
// The method logs a panic message with the provided arguments.
// The message includes the logger's traceKey and traceId if they are available.
// After logging the message, the method calls panic() to stop the ordinary flow of a goroutine.
func (logger *Logger) Panicln(args ...interface{}) {
	logger.withTraceKey().Panicln(args...)
}
