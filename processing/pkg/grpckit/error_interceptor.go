package grpckit

import (
	"context"

	"google.golang.org/grpc"
)

type ErrorConverter func(ctx context.Context, err error) error

func UnaryServerInterceptor(errorConverter ErrorConverter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo, //nolint:revive
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			err = errorConverter(ctx, err)
		}

		return resp, err
	}
}

func StreamServerInterceptor(errorConverter ErrorConverter) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo, //nolint:revive
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)
		if err != nil {
			err = errorConverter(ss.Context(), err)
		}

		return err
	}
}
