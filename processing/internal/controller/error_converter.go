package controller

import (
	"context"
	"errors"

	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/grpckit"
	sdkErrorProto "code.emcdtech.com/emcd/sdk/error/proto"
	"code.emcdtech.com/emcd/sdk/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errorCodeToGRPCCode = map[model.ErrorCode]codes.Code{
	model.ErrorCodeNoSuchMerchant:    codes.NotFound,
	model.ErrorCodeNoSuchInvoiceForm: codes.NotFound,
	model.ErrorCodeInternal:          codes.Internal,
	model.ErrorCodeInvalidArgument:   codes.InvalidArgument,
	model.ErrorCodeNoSuchInvoice:     codes.NotFound,
}

func NewErrorConverter() grpckit.ErrorConverter {
	return func(ctx context.Context, err error) error {
		if _, ok := status.FromError(err); ok {
			return err
		}

		var modelErr *model.Error
		if !errors.As(err, &modelErr) {
			log.SError(ctx, "unhandled error", map[string]any{"error": err})

			modelErr = model.ErrorWithDefaultMessage(&model.Error{Code: model.ErrorCodeInternal})
		}

		log.SError(ctx, "request finished with error", map[string]any{"error": modelErr})

		grpcCode := codes.Unknown
		if code, ok := errorCodeToGRPCCode[modelErr.Code]; ok {
			grpcCode = code
		}

		// TODO: use sdk/error
		st := status.New(grpcCode, modelErr.Message)
		st, _ = st.WithDetails(&sdkErrorProto.ErrorDetails{
			Code:    string(modelErr.Code),
			Message: modelErr.Message,
		})

		return st.Err()
	}
}
