package log

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientUnaryInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	const fieldUserIdSnake = "user_id"
	const fieldUserIdCamel = "UserId"
	usrID := ctx.Value(userID)
	if usrID != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, strings.ToLower(userID), usrID.(string))
	} else {
		// let's try to get it from req
		immutable := reflect.ValueOf(req)
		var field reflect.Value
		if immutable.Kind() == reflect.Ptr {
			if structValue := reflect.Indirect(immutable); structValue.IsValid() {
				field = structValue.FieldByName(fieldUserIdSnake)
				if field == (reflect.Value{}) {
					field = structValue.FieldByName(fieldUserIdCamel)
				}
			}
		} else {
			field = immutable.FieldByName(fieldUserIdSnake)
			if field == (reflect.Value{}) {
				field = immutable.FieldByName(fieldUserIdCamel)
			}
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field != (reflect.Value{}) {
			if field.Kind() == reflect.String {
				ctx = metadata.AppendToOutgoingContext(ctx, strings.ToLower(userID), field.String())
			}
		}
	}

	if serviceNameCtx := ctx.Value(serviceNameStruct{}); serviceNameCtx != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, serviceNameStruct{}.name(), serviceNameCtx.(string))

	}

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

// ServerUnaryInterceptor interceptor runs postprocessing of error returned from gRPC handler.
func ServerUnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
	mdUser := metadata.ValueFromIncomingContext(ctx, strings.ToLower(userID))

	if len(mdUser) > 0 {
		ctx = context.WithValue(ctx, userID, mdUser[0])
	}

	return h(ctx, req)
}

func ServerUnaryNamedInterceptor(serverNameI string) func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
	serverName := strings.ToUpper(serverNameI)

	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if userMd, ok := md[userID]; ok && len(userMd) > 0 {
				ctx = context.WithValue(ctx, userID, userMd[0])

			}

			if serviceNameMd, ok := md[serviceNameStruct{}.name()]; ok && len(serviceNameMd) > 0 {
				ctx = context.WithValue(ctx, serviceNameStruct{}, fmt.Sprintf("%s.%s", serviceNameMd[0], serverName))

			} else {
				ctx = context.WithValue(ctx, serviceNameStruct{}, serverName)

			}
		}

		return h(ctx, req)
	}
}
