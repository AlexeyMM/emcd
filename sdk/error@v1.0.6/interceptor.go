package errors

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "code.emcdtech.com/emcd/sdk/error/proto"
)

// ClientUnaryInterceptor interceptor runs preprocessing of error returned from gRPC call.
func ClientUnaryInterceptor(
	ctx context.Context,
	method string,
	req any,
	reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	pbStatus, ok := status.FromError(err)
	if !ok {
		return err
	}

	for _, d := range pbStatus.Details() {
		if details, ok := d.(*pb.ErrorDetails); ok {
			return NewError(details.Code, details.Message)
		}
	}

	return err
}

// ServerUnaryInterceptor interceptor runs postprocessing of error returned from gRPC handler.
func ServerUnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
	res, err := h(ctx, req)
	if err == nil {
		return res, nil
	}

	if _, ok := status.FromError(err); ok {
		return nil, err
	}

	var e *Error
	if errors.As(err, &e) {
		return nil, NewGRPCError(e.Code, e.Message)
	}

	return nil, err
}
