package _client

//
// import (
// 	"context"
//
// 	"google.golang.org/grpc"
// 	"google.golang.org/protobuf/types/known/emptypb"
//
// 	"code.emcdtech.com/emcd/blockchain/address/protocol/address"
// )
//
// type AddressClient interface {
// 	GetHello(ctx context.Context) (string, error)
// 	GetHelloError(ctx context.Context) (string, error)
// }
//
// type addressImpl struct {
// 	client address.AddressServiceClient
// }
//
// func NewAddressClient(conn *grpc.ClientConn) *addressImpl {
// 	return &addressImpl{
// 		client: address.NewAddressServiceClient(conn),
// 	}
// }
//
// func (a *addressImpl) GetHello(ctx context.Context) (string, error) {
// 	resp, err := a.client.Hello(ctx, &emptypb.Empty{})
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.GetHello(), nil
// }
//
// func (a *addressImpl) GetHelloError(ctx context.Context) (string, error) {
// 	resp, err := a.client.HelloError(ctx, &emptypb.Empty{})
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.GetHello(), nil
// }
