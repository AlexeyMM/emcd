package log

import (
	"context"
	"io"

	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

var _ Logger = loggerImpl{}

type loggerImpl struct {
	Logger
	zl zerolog.Logger
}

func (s loggerImpl) Err(ctx context.Context, err error) {
	fields := s.getFields(ctx)
	s.zl.Error().Fields(fields).CallerSkipFrame(2).Err(err).Send()
}

func (s loggerImpl) Info(ctx context.Context, format string, v ...interface{}) {
	fields := s.getFields(ctx)
	s.zl.Info().Fields(fields).CallerSkipFrame(2).Msgf(format, v...)
}

func (s loggerImpl) Debug(ctx context.Context, format string, v ...interface{}) {
	fields := s.getFields(ctx)
	s.zl.Debug().Fields(fields).CallerSkipFrame(2).Msgf(format, v...)
}

func (s loggerImpl) Warn(ctx context.Context, format string, v ...interface{}) {
	fields := s.getFields(ctx)
	s.zl.Warn().Fields(fields).CallerSkipFrame(2).Msgf(format, v...)
}

func (s loggerImpl) Error(ctx context.Context, format string, v ...interface{}) {
	fields := s.getFields(ctx)
	s.zl.Error().Fields(fields).CallerSkipFrame(2).Msgf(format, v...)
}

func (s loggerImpl) Panic(ctx context.Context, format string, v ...interface{}) {
	fields := s.getFields(ctx)
	s.zl.Panic().Fields(fields).CallerSkipFrame(2).Msgf(format, v...)
}

func (s loggerImpl) Fatal(ctx context.Context, format string, v ...interface{}) {
	fields := s.getFields(ctx)
	s.zl.Fatal().Fields(fields).CallerSkipFrame(2).Msgf(format, v...)
}

func (s loggerImpl) SDebug(ctx context.Context, msg string, fields map[string]any) {
	s.zl.Debug().
		Fields(s.getFields(ctx)).
		Fields(fields).
		CallerSkipFrame(2).
		Msg(msg)
}

func (s loggerImpl) SInfo(ctx context.Context, msg string, fields map[string]any) {
	s.zl.Info().
		Fields(s.getFields(ctx)).
		Fields(fields).
		CallerSkipFrame(2).
		Msg(msg)
}

func (s loggerImpl) SWarn(ctx context.Context, msg string, fields map[string]any) {
	s.zl.Warn().
		Fields(s.getFields(ctx)).
		Fields(fields).
		CallerSkipFrame(2).
		Msg(msg)
}

func (s loggerImpl) SError(ctx context.Context, msg string, fields map[string]any) {
	s.zl.Error().
		Fields(s.getFields(ctx)).
		Fields(fields).
		CallerSkipFrame(2).
		Msg(msg)
}

func (s loggerImpl) SPanic(ctx context.Context, msg string, fields map[string]any) {
	s.zl.Panic().
		Fields(s.getFields(ctx)).
		Fields(fields).
		CallerSkipFrame(2).
		Msg(msg)
}

func (s loggerImpl) SFatal(ctx context.Context, msg string, fields map[string]any) {
	s.zl.Fatal().
		Fields(s.getFields(ctx)).
		Fields(fields).
		CallerSkipFrame(2).
		Msg(msg)
}

func (s loggerImpl) getFields(ctx context.Context) map[string]interface{} {
	res := make(map[string]interface{})

	usrID := ctx.Value(userID)
	if usrID != nil {
		res[userID] = usrID.(string)
	} else {
		usrID = ctx.Value(oldUserId)
		if usrID != nil {
			res[userID] = usrID.(string)
		}
	}

	if serviceNameCtx := ctx.Value(serviceNameStruct{}); serviceNameCtx != nil {
		res[serviceNameStruct{}.name()] = serviceNameCtx.(string)
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		res[ecsTraceID] = spanContext.TraceID().String()
		res[ecsSpanID] = spanContext.SpanID().String()
	}
	return res
}

func newZeroLogger(out io.Writer) loggerImpl {
	return loggerImpl{
		zl: zerolog.New(out).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger(),
	}
}
