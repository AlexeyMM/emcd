package server

import (
	"code.emcdtech.com/emcd/service/profile/internal/service"
	pb "code.emcdtech.com/emcd/service/profile/protocol/profile"
	"context"
)

const (
	opGetOrCreateSecretKey = "server.profile.GetOrCreateSecretKey"
	opCheckSignature       = "server.profile.CheckSignature"
)

func (p *Profile) GetOrCreateSecretKey(ctx context.Context, req *pb.GetOrCreateSecretKeyRequest) (*pb.SecretKeyResponse, error) {
	userUUID, parentUUID, err := p.validateApiKeyRequest(ctx, req.GetUserUuid(), req.GetParentUuid(), req.GetIp(), opGetOrCreateSecretKey)
	if err != nil {
		return nil, err
	}

	key, err := p.profileService.GetOrCreateSecretKey(ctx, service.GetOrCreateSecretKeyReq{
		UserUUID:   userUUID,
		ParentUUID: parentUUID,
		IP:         req.GetIp(),
	})
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.SecretKeyResponse{
		SecretKey: key,
	}, nil
}

func (p *Profile) CheckSignature(ctx context.Context, req *pb.Signature) (*pb.SignatureResponse, error) {
	userUUID, parentUUID, err := p.validateApiKeyRequest(ctx, req.GetUserUuid(), req.GetParentUuid(), req.GetIp(), opCheckSignature)
	if err != nil {
		return nil, err
	}

	check, err := p.profileService.CheckSignature(ctx, service.CheckSignatureReq{
		UserUUID:   userUUID,
		ParentUUID: parentUUID,
		Signature:  req.Signature,
		Nonce:      req.Nonce,
	})
	if err != nil {
		return nil, toProtoError(err)
	}

	return &pb.SignatureResponse{
		Check: check,
	}, nil
}
