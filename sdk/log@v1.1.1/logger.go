package log

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	oldUserId  = "user_id"
	userID     = "user.id"
	ecsSpanID  = "span.id"
	ecsTraceID = "trace.id"
)

type serviceNameStruct struct{} // for context value using

func (serviceNameStruct) name() string { // for other

	return "service.name"
}

// logger is the global logger, init by default
var logger Logger

type Log struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

// _init - init logger with custom parameters
func _init() error {
	if logger != nil {

		return nil
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	env, exists := os.LookupEnv("LOG_LEVEL")
	if !exists {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		logLevel, err := zerolog.ParseLevel(env)
		if err != nil {
			return errors.WithStack(err)
		}

		zerolog.SetGlobalLevel(logLevel)
	}

	logger = newZeroLogger(os.Stderr)

	return nil
}

func init() {
	if err := _init(); err != nil {
		log.Fatal(fmt.Errorf("failed to initialize logger: %w", err))
	}
}

type Logger interface {
	Err(ctx context.Context, err error)
	Debug(ctx context.Context, format string, v ...interface{})
	Info(ctx context.Context, format string, v ...interface{})
	Warn(ctx context.Context, format string, v ...interface{})
	Error(ctx context.Context, format string, v ...interface{})
	Panic(ctx context.Context, format string, v ...interface{})
	Fatal(ctx context.Context, format string, v ...interface{})
	// Methods for structured logging
	SDebug(ctx context.Context, msg string, fields map[string]any)
	SInfo(ctx context.Context, msg string, fields map[string]any)
	SWarn(ctx context.Context, msg string, fields map[string]any)
	SError(ctx context.Context, msg string, fields map[string]any)
	SPanic(ctx context.Context, msg string, fields map[string]any)
	SFatal(ctx context.Context, msg string, fields map[string]any)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
func Err(ctx context.Context, err error) {
	logger.Err(ctx, err)
}

// Debug starts a new message with debug level.
func Debug(ctx context.Context, format string, v ...interface{}) {
	logger.Debug(ctx, format, v...)
}

// Info starts a new message with info level.
func Info(ctx context.Context, format string, v ...interface{}) {
	logger.Info(ctx, format, v...)
}

// Warn starts a new message with warn level.
func Warn(ctx context.Context, format string, v ...interface{}) {
	logger.Warn(ctx, format, v...)
}

// Error starts a new message with error level.
func Error(ctx context.Context, format string, v ...interface{}) {
	logger.Error(ctx, format, v...)
}

// Panic starts a new message with error level.
func Panic(ctx context.Context, format string, v ...interface{}) {
	logger.Panic(ctx, format, v...)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
func Fatal(ctx context.Context, format string, v ...interface{}) {
	logger.Fatal(ctx, format, v...)
}

func SDebug(ctx context.Context, msg string, fields map[string]any) {
	logger.SDebug(ctx, msg, fields)
}
func SInfo(ctx context.Context, msg string, fields map[string]any) {
	logger.SInfo(ctx, msg, fields)
}
func SWarn(ctx context.Context, msg string, fields map[string]any) {
	logger.SWarn(ctx, msg, fields)
}
func SError(ctx context.Context, msg string, fields map[string]any) {
	logger.SError(ctx, msg, fields)
}
func SPanic(ctx context.Context, msg string, fields map[string]any) {
	logger.SPanic(ctx, msg, fields)
}
func SFatal(ctx context.Context, msg string, fields map[string]any) {
	logger.SFatal(ctx, msg, fields)
}
