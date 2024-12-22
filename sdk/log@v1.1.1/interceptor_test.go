package log

import (
	"context"
	"strings"
	"testing"

	"code.emcdtech.com/emcd/sdk/log/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestClientUnaryInterceptor(t *testing.T) {
	testUserID := "test_user_id"
	request := proto.TestRequest{
		UserId: &testUserID,
	}
	request2 := proto.TestRequestStatic{
		UserId: testUserID,
	}

	trID := uuid.New().String()
	ctx := context.WithValue(context.Background(), traceID, trID)
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		md, _ := metadata.FromOutgoingContext(ctx)
		vals := md.Get(strings.ToLower(traceID))
		assert.Len(t, vals, 1)
		assert.Equal(t, vals[0], trID)

		vals = md.Get(strings.ToLower(userID))
		assert.Len(t, vals, 1)
		if len(vals) > 0 {
			assert.Equal(t, vals[0], testUserID)
		}

		return nil
	}

	err := ClientUnaryInterceptor(ctx, "user-method", &request, nil, nil, invoker)
	assert.NoError(t, err)

	err = ClientUnaryInterceptor(ctx, "user-method-2", request2, nil, nil, invoker)
	assert.NoError(t, err)

	err = ClientUnaryInterceptor(ctx, "user-method-3", &request2, nil, nil, invoker)
	assert.NoError(t, err)
}

func TestServerUnaryInterceptor(t *testing.T) {
	ctx := context.Background()

	data := make(map[string]string)
	data[strings.ToLower(traceID)] = "traceID"
	data[strings.ToLower(userID)] = "userID"

	ctx = metadata.NewIncomingContext(ctx, metadata.New(data))
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		trID := ctx.Value(traceID)
		assert.Equal(t, trID.(string), "traceID")
		usrID := ctx.Value(userID)
		assert.Equal(t, usrID.(string), "userID")
		return nil, nil
	}
	_, err := ServerUnaryInterceptor(ctx, "interface{}", nil, handler)
	assert.NoError(t, err)
}
