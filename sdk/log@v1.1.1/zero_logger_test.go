package log

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/apmtest"
)

func TestGetFieldStartTracing(t *testing.T) {
	var fields map[string]interface{}
	tx, spans, _ := apmtest.WithTransaction(func(ctx context.Context) {
		span, ctx := apm.StartSpan(ctx, "name", "type")
		defer span.End()

		log := loggerImpl{}
		fields = log.getFields(ctx)
	})

	transactionID, ok := fields[ecsTransactionID]
	require.Equal(t, ok, true)
	require.Equal(t, transactionID, fmt.Sprintf("%x", tx.ID))

	spanId, ok := fields[ecsSpanID]
	require.Equal(t, ok, true)
	require.Equal(t, spanId, fmt.Sprintf("%x", spans[0].ID))

	traceId, ok := fields[ecsTraceID]
	require.Equal(t, ok, true)
	require.Equal(t, traceId, fmt.Sprintf("%x", tx.TraceID))
}

func TestGetFieldNonTracing(t *testing.T) {
	fields := loggerImpl{}.getFields(context.Background())

	_, ok := fields[ecsTransactionID]
	require.Equal(t, ok, false)
	_, ok = fields[ecsSpanID]
	require.Equal(t, ok, false)
	_, ok = fields[ecsTraceID]
	require.Equal(t, ok, false)
}
