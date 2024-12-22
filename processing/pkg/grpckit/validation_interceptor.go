package grpckit

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"code.emcdtech.com/b2b/processing/model"
)

func ValidationUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo, //nolint:revive
		handler grpc.UnaryHandler,
	) (any, error) {
		message, isMessage := req.(proto.Message)
		if isMessage {
			if err := protovalidate.Validate(message); err != nil {
				var validationErr *protovalidate.ValidationError
				if errors.As(err, &validationErr) {
					err = &model.Error{
						Code:    model.ErrorCodeInvalidArgument,
						Message: validationErr.Error(),
						Inner:   validationErr,
					}
				}

				return nil, err
			}
		}

		return handler(ctx, req)
	}
}
