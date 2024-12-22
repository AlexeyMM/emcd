package server

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/sdk/log"

	pb "code.emcdtech.com/emcd/service/profile/protocol/profile"
)

func (p *Profile) GetNoPayStatus(ctx context.Context, req *pb.GetNoPayStatusRequest) (*pb.GetNoPayStatusResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		log.Error(ctx, "GetNoPayStatus: parse user uuid: %s. %v", req.GetUserUuid(), err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	noPayStatus, timeRemaining, err := p.profileService.GetNoPayStatus(ctx, userUUID)
	if err != nil {
		return nil, toProtoError(err)
	}

	resp := &pb.GetNoPayStatusResponse{
		Status: noPayStatus,
	}

	if noPayStatus {
		resp.DateBefore = timestamppb.New(timeRemaining)
	}

	return resp, nil
}

func (p *Profile) UpdateNoPay(ctx context.Context, req *pb.UpdateNoPayRequest) (*pb.UpdateNoPayResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		log.Error(ctx, "UpdateNoPay: parse user uuid: %s. %v", req.GetUserUuid(), err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = p.profileService.UpdateNoPayToFalse(ctx, userUUID)
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.UpdateNoPayResponse{}, nil
}

func (p *Profile) CancelNoPayJob(ctx context.Context, req *pb.CancelNoPayJobRequest) (*pb.CancelNoPayJobResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		log.Error(ctx, "CancelNoPayJob: parse user uuid: %s. %v", req.GetUserUuid(), err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = p.profileService.CancelJobOffNoPay(ctx, userUUID)
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.CancelNoPayJobResponse{}, nil
}
