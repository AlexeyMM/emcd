package log

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerInit(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, userID, "usrID")
	err := Init(ctx)
	assert.NoError(t, err)

	Info(ctx, "hey there")
	Info(ctx, "hey there2")
	Info(ctx, "hey there3")
	Debug(ctx, "debug")
	Warn(ctx, "warn")
	Err(ctx, errors.New("new error"))
}

func TestLoggerInitNoValues(t *testing.T) {
	ctx := context.Background()
	err := Init(ctx)
	assert.NoError(t, err)

	Info(ctx, "hey there")
}

func TestLoggerNoInit(t *testing.T) {
	Info(context.Background(), "hey there")
}

func TestLoggerContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, userID, "usrID")
	Info(ctx, "hey")
}
