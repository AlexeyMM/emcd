package repository

import (
	"context"

	wlPb "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"google.golang.org/grpc"
)

//go:generate mockery --name=WhiteLabel --structname=MockWhiteLabel --outpkg=mocks --output ./mocks --filename $GOFILE

type WhiteLabel interface {
	GetCoins(ctx context.Context, in *wlPb.GetCoinsRequest, opts ...grpc.CallOption) (*wlPb.GetCoinsResponse, error)
}
