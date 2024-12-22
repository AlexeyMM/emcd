package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/service"
	pb "code.emcdtech.com/emcd/service/profile/protocol/profile"
)

const (
	opGetNewAPIKey = "server.profile.GetNewAPIKey"
	opDeleteAPIKey = "server.profile.DeleteAPIKey"
	opGetAPIKey    = "server.profile.GetAPIKey"
)

func (p *Profile) CreateAPIKey(ctx context.Context, req *pb.CreateAPIKeyRequest) (*pb.CreateAPIKeyResponse, error) {
	userUUID, parentUUID, err := p.validateApiKeyRequest(ctx, req.GetUserUuid(), req.GetParentUuid(), req.GetIp(), opGetNewAPIKey)
	if err != nil {
		return nil, err
	}

	apiKey, err := p.profileService.CreateNewAPIKey(ctx, service.CreateNewAPIKeyReq{
		UserUUID:   userUUID,
		ParentUUID: parentUUID,
		IP:         req.GetIp(),
	})
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.CreateAPIKeyResponse{
		ApiKey: apiKey,
	}, nil
}

func (p *Profile) GetAPIKey(ctx context.Context, req *pb.GetAPIKeyRequest) (*pb.GetAPIKeyResponse, error) {
	userUUID, parentUUID, err := p.validateApiKeyRequest(ctx, req.GetUserUuid(), req.GetParentUuid(), req.GetIp(), opGetAPIKey)
	if err != nil {
		return nil, err
	}

	apiKey, err := p.profileService.GetAPIKey(ctx, service.GetAPIKeyReq{
		UserUUID:   userUUID,
		ParentUUID: parentUUID,
		IP:         req.GetIp(),
	})
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.GetAPIKeyResponse{
		ApiKey: apiKey,
	}, nil
}

func (p *Profile) DeleteAPIKey(ctx context.Context, req *pb.DeleteAPIKeyRequest) (*pb.DeleteAPIKeyResponse, error) {
	userUUID, parentUUID, err := p.validateApiKeyRequest(ctx, req.GetUserUuid(), req.GetParentUuid(), req.GetIp(), opDeleteAPIKey)
	if err != nil {
		return nil, err
	}

	err = p.profileService.DeleteAPIKey(ctx, service.DeleteAPIKeyReq{
		UserUUID:   userUUID,
		ParentUUID: parentUUID,
		IP:         req.GetIp(),
	})
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.DeleteAPIKeyResponse{}, nil
}

func (p *Profile) validateApiKeyRequest(
	ctx context.Context,
	userUUID, parentUUID, ip, logStr string,
) (userID, parentID uuid.UUID, err error) {
	if userUUID == "" {
		err = status.Error(codes.InvalidArgument, "no user uuid")
		return
	}

	if ip == "" {
		err = status.Error(codes.InvalidArgument, "no ip")
		return
	}

	userID, err = uuid.Parse(userUUID)
	if err != nil {
		log.Error(ctx, "%s: parse user uuid from: %s: %s", logStr, userUUID, err.Error())
		err = status.Error(codes.InvalidArgument, "parse user uuid")
		return
	}

	if parentUUID != "" {
		parentID, err = uuid.Parse(parentUUID)
		if err != nil {
			log.Error(ctx, "%s: parse parent uuid from: %s: %s", logStr, parentUUID, err.Error())
			err = status.Error(codes.InvalidArgument, "parse parent uuid")
		}
	}
	return
}

func toProtoError(err error) error {
	switch {
	case errors.Is(err, service.ErrAlreadyExist):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
