package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	businessErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	pb "code.emcdtech.com/emcd/service/profile/protocol/profile"
)

var (
	errInconsistentUserInfo    = businessErr.NewError("pr-00001", service.ErrInconsistentUserInfo.Error())
	errAddressChangeNotAllowed = businessErr.NewError("pr-00002", "address change not allowed")
	errBothAddressIsTheSame    = businessErr.NewError("pr-00003", "both address is the same")
	errMinPayNotValid          = businessErr.NewError("pr-00004", "min pay value is not valid")
	errCoinNotFound            = businessErr.NewError("pr-00005", service.ErrCoinNotFound.Error())
	errTokenNotFound           = businessErr.NewError("pr-00006", service.ErrTokenNotFound.Error())
	errEmailIsEmpty            = businessErr.NewError("pr-00007", "email is empty")
)

const (
	defaultPoolType = "emcd"
)

type oldIdResolver interface {
	GetIDs(ctx context.Context, ids []uuid.UUID) ([]*model.ID, error)
}

type Profile struct {
	pb.UnimplementedProfileServiceServer

	profileService service.Profile
	oldIDResolver  oldIdResolver
}

func NewProfile(
	profileSrv service.Profile,
	oldIDResolver oldIdResolver,
) *Profile {
	return &Profile{
		profileService: profileSrv,
		oldIDResolver:  oldIDResolver,
	}
}

func (p *Profile) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	pr, err := p.parseProto(req.Profile)
	if err != nil {
		log.Error(ctx, "profile: create: %v", err)
		return nil, fmt.Errorf("profile: create: %w", err)
	}
	id, err := p.profileService.Create(ctx, pr)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.CreateResponse{
		UserID: id,
	}, nil
}

func (p *Profile) SaveV3(ctx context.Context, req *pb.SaveV3Request) (*pb.SaveV3Response, error) {
	const op = "server.profile.SaveV3"
	profile, err := p.parseProtoV3(req.Profile)
	if err != nil {
		log.Error(ctx, "%s: %v", op, err)
		return nil, fmt.Errorf("profile: save: %w", err)
	}
	id, err := p.profileService.SaveV3(ctx, profile)
	if err != nil {
		log.Error(ctx, "%s: profile: %+v: %v", op, profile, err)
		return nil, err
	}
	return &pb.SaveV3Response{
		UserID: id,
	}, nil
}

func (p *Profile) parseProto(pr *pb.Profile) (*model.Profile, error) {
	userID, err := uuid.Parse(pr.User.ID)
	if err != nil {
		return nil, fmt.Errorf("parse proto: parse user id %v: %w", pr.User.ID, err)
	}
	wlID, err := uuid.Parse(pr.User.WhiteLabelID)
	if err != nil {
		return nil, fmt.Errorf("parse proto: parse wl id %v: %w", pr.User.WhiteLabelID, err)
	}
	var newRefID uuid.UUID
	if pr.User.NewRefId != "" {
		newRefID, err = uuid.Parse(pr.User.NewRefId)
		if err != nil {
			return nil, fmt.Errorf("parse new_ref_id: %w %s", err, pr.User.NewRefId)
		}
	}
	poolType := defaultPoolType
	if pr.User != nil && pr.User.PoolType != nil {
		poolType = *pr.User.PoolType
		// ? looks strange, but *bug*
		// if client using old library version it pass empty string instead of nil pointer
		if len(poolType) == 0 {
			poolType = defaultPoolType
		}
	}
	return &model.Profile{
		User: &model.User{
			ID:           userID,
			Username:     pr.User.Username,
			Vip:          pr.User.Vip,
			SegmentID:    int(pr.User.SegmentID),
			RefID:        int(pr.User.RefID),
			Email:        pr.User.Email,
			Password:     pr.User.Password,
			CreatedAt:    pr.User.CreatedAt.AsTime(),
			WhiteLabelID: wlID,
			ApiKey:       pr.User.ApiKey,
			NewRefID:     newRefID,
			PoolType:     poolType,
			OldParentID:  pr.User.OldParentId,
		},
	}, nil
}

func (p *Profile) parseProtoV3(pr *pb.ProfileV4) (*model.Profile, error) {
	userID, err := uuid.Parse(pr.User.ID)
	if err != nil {
		return nil, fmt.Errorf("parse user_id %v: %w", pr.User.ID, err)
	}
	wlID, err := uuid.Parse(pr.User.WhiteLabelID)
	if err != nil {
		return nil, fmt.Errorf("parse whiteLabelID %v: %w", pr.User.WhiteLabelID, err)
	}
	var parentID uuid.UUID
	if pr.User.ParentId != "" {
		parentID, err = uuid.Parse(pr.User.ParentId)
		if err != nil {
			return nil, fmt.Errorf("parse parent_id %v: %w", pr.User.ParentId, err)
		}
	}
	var newRefID uuid.UUID
	if pr.User.NewRefId != "" {
		newRefID, err = uuid.Parse(pr.User.NewRefId)
		if err != nil {
			return nil, fmt.Errorf("parse new_ref_id: %w %s", err, pr.User.NewRefId)
		}
	}
	poolType := defaultPoolType
	if pr.User != nil && pr.User.PoolType != nil {
		poolType = *pr.User.PoolType
		// ? looks strange, but *bug*
		// if client using old library version it pass empty string instead of nil pointer
		if len(poolType) == 0 {
			poolType = defaultPoolType
		}

	}

	return &model.Profile{
		User: &model.User{
			ID:           userID,
			Username:     pr.User.Username,
			Vip:          pr.User.Vip,
			SegmentID:    int(pr.User.SegmentID),
			RefID:        int(pr.User.RefID),
			ParentID:     parentID,
			Email:        pr.User.Email,
			Password:     pr.User.Password,
			CreatedAt:    pr.User.CreatedAt.AsTime(),
			WhiteLabelID: wlID,
			ApiKey:       pr.User.ApiKey,
			IsActive:     pr.User.IsActive,
			AppleID:      pr.User.AppleId,
			NewRefID:     newRefID,
			PoolType:     poolType,
			OldParentID:  pr.User.OldParentId,
		},
	}, nil
}

func (p *Profile) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	if req.Email == "" {
		return nil, errEmailIsEmpty
	}
	wlID := uuid.Nil
	err := p.profileService.UpdatePassword(ctx, req.Email, req.Password, wlID)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.UpdatePasswordResponse{}, nil
}

func (p *Profile) GetUserByEmailAndWl(
	ctx context.Context,
	req *pb.GetUserByEmailAndWlRequest,
) (*pb.GetUserByEmailAndWlResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "profile: get user by email and wl: parsing wl id %s: %v", req.WhiteLabelID, err)
		return nil, fmt.Errorf("parse: %w", err)
	}
	u, err := p.profileService.GetUserByEmailAndWl(ctx, req.Email, wlID)
	if err != nil {
		log.Error(ctx, "GetUserByEmailAndWl %v. email: %s", err, req.Email)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	return &pb.GetUserByEmailAndWlResponse{
		User: p.convertToProtoUser(u),
	}, nil
}

func (p *Profile) GetOldUserByEmailAndWl(
	ctx context.Context,
	req *pb.GetOldUserByEmailAndWlRequest,
) (*pb.GetOldUserByEmailAndWlResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "profile: get old user by email and wl: parsing wl id %s: %v", req.WhiteLabelID, err)
		return nil, fmt.Errorf("profile: get old user by email and wl: parsing wl id: %w", err)
	}
	u, err := p.profileService.GetOldUserByEmailAndWl(ctx, req.Email, wlID)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.GetOldUserByEmailAndWlResponse{
		User: p.convertToProtoUser(u),
	}, nil
}

func (p *Profile) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, totalCount, err := p.profileService.GetAllUsers(
		ctx,
		int(req.Skip),
		int(req.Take),
		req.Sort.Field,
		req.Sort.Asc,
		req.SearchEmail,
	)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	protoUsers := make([]*pb.User, len(users))
	for i := range users {
		protoUsers[i] = p.convertToProtoUser(users[i])
	}
	return &pb.GetAllUsersResponse{
		Users:      protoUsers,
		TotalCount: int32(totalCount),
	}, nil
}

func (p *Profile) GetAllUsersByWlID(
	ctx context.Context,
	req *pb.GetAllUsersByWlIDRequest,
) (*pb.GetAllUsersByWlIDResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "profile: get all users by wl id %v: %v", req.WhiteLabelID, err)
		return nil, fmt.Errorf("profile: get all users by wl id: %w", err)
	}
	users, totalCount, err := p.profileService.GetAllUsersByWlID(
		ctx,
		int(req.Skip),
		int(req.Take),
		req.Sort.Field,
		req.Sort.Asc,
		req.SearchEmail,
		wlID,
	)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	protoUsers := make([]*pb.User, len(users))
	for i := range users {
		protoUsers[i] = p.convertToProtoUser(users[i])
	}
	return &pb.GetAllUsersByWlIDResponse{
		Users:      protoUsers,
		TotalCount: int32(totalCount),
	}, nil
}

func (p *Profile) GetUserByApiKey(ctx context.Context, req *pb.GetUserByApiKeyRequest) (*pb.GetUserByApiKeyResponse, error) {
	u, err := p.profileService.GetUserByApiKey(ctx, req.ApiKey)
	if err != nil {
		log.Error(ctx, err.Error())
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	return &pb.GetUserByApiKeyResponse{
		User: p.convertToProtoUser(u),
	}, nil
}

func (p *Profile) convertToProtoUser(u *model.User) *pb.User {
	if u == nil {
		return nil
	}

	poolType := u.PoolType
	// ? looks strange, but *bug*
	// if client using old library version it pass empty string instead of nil pointer
	if len(poolType) == 0 {
		poolType = defaultPoolType
	}
	return &pb.User{
		ID:           u.ID.String(),
		Username:     u.Username,
		Vip:          u.Vip,
		SegmentID:    int32(u.SegmentID),
		RefID:        int32(u.RefID),
		ParentId:     u.ParentID.String(),
		Email:        u.Email,
		Password:     u.Password,
		CreatedAt:    timestamppb.New(u.CreatedAt),
		WhiteLabelID: u.WhiteLabelID.String(),
		ApiKey:       u.ApiKey,
		OldID:        u.OldID,
		Suspended:    u.Suspended,
		IsActive:     u.IsActive,
		NewRefId:     u.NewRefID.String(),
		Language:     u.Language,
		PoolType:     &poolType,
		TgId:         u.TgID,
		TgUsername:   u.TgUsername,
		IsAmbassador: u.IsAmbassador,
	}
}

func (p *Profile) convertToProtoUserV2(u *model.User) *pb.UserV2 {
	if u == nil {
		return nil
	}

	poolType := u.PoolType
	// ? looks strange, but *bug*
	// if client using old library version it pass empty string instead of nil pointer
	if len(poolType) == 0 {
		poolType = defaultPoolType
	}

	return &pb.UserV2{
		ID:           u.ID.String(),
		Username:     u.Username,
		Vip:          u.Vip,
		SegmentID:    int32(u.SegmentID),
		RefID:        int32(u.RefID),
		Email:        u.Email,
		Password:     u.Password,
		CreatedAt:    timestamppb.New(u.CreatedAt),
		WhiteLabelID: u.WhiteLabelID.String(),
		ApiKey:       u.ApiKey,
		OldID:        u.OldID,
		IsActive:     u.IsActive,
		AppleId:      u.AppleID,
		ParentId:     u.ParentID.String(),
		NewRefId:     u.NewRefID.String(),
		PoolType:     &poolType,
		Language:     u.Language,
		TgId:         u.TgID,
		TgUsername:   u.TgUsername,
	}
}

func (p *Profile) GetByUserID(ctx context.Context, req *pb.GetByUserIDRequest) (*pb.GetByUserIDResponse, error) {
	const op = "server.profile.GetByUserID"
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error(ctx, "%s: parse user id %s: %v", op, req.UserID, err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: user id: %s", op, req.UserID))
	}
	pr, err := p.profileService.GetByUserID(ctx, userID)
	if err != nil {
		log.Error(ctx, "%s: get user by id: %s: %v", op, req.UserID, err)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	if pr == nil {
		log.Warn(ctx, "%s: requested unexisted user id: %s: %v", op, req.UserID, err)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("users id %s", req.UserID))
	}
	return &pb.GetByUserIDResponse{
		Profile: &pb.Profile{
			User: p.convertToProtoUser(pr.User),
		},
	}, nil
}

func (p *Profile) GetByUserIDV2(ctx context.Context, req *pb.GetByUserIDV2Request) (*pb.GetByUserIDV2Response, error) {
	const op = "server.profile.GetByUserIDV2"
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error(ctx, "%s: parse user id %s: %v", op, req.UserID, err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: user id: %s", op, req.UserID))
	}
	pr, err := p.profileService.GetByUserID(ctx, userID)
	if err != nil {
		log.Error(ctx, "%s: get user by id: %s: %v", op, req.UserID, err)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	if pr == nil {
		log.Warn(ctx, "%s: requested unexisted user id: %s: %v", op, req.UserID, err)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("users id %s", req.UserID))
	}
	return &pb.GetByUserIDV2Response{
		Profile: &pb.ProfileV2{
			User: p.convertToProtoUserV2(pr.User),
		},
	}, nil
}

func (p *Profile) GetUserIDByOldID(ctx context.Context, req *pb.GetUserIDByOldIDRequest) (*pb.GetUserIDByOldIDResponse, error) {
	userID, err := p.profileService.GetUserIDByOldID(ctx, req.OldId)
	if err != nil {
		log.Error(ctx, "%v. old id: %d", err, req.OldId)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	return &pb.GetUserIDByOldIDResponse{
		UserID: userID.String(),
	}, nil
}

func (p *Profile) GetByUsernames(ctx context.Context, req *pb.GetByUsernamesRequest) (*pb.GetByUsernamesResponse, error) {
	pr, err := p.profileService.GetByUsernames(ctx, req.Usernames)
	if err != nil {
		log.Error(ctx, "%v. usernames: %s", err, req.Usernames)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	if pr == nil {
		return &pb.GetByUsernamesResponse{}, nil
	}

	profiles := make([]*pb.Profile, 0)
	for _, profile := range pr {
		profiles = append(profiles, &pb.Profile{
			User: p.convertToProtoUser(profile.User),
		})
	}
	return &pb.GetByUsernamesResponse{
		Profile: profiles,
	}, nil
}

func (p *Profile) GetOldByEmailAndWl(
	ctx context.Context,
	req *pb.GetOldByEmailAndWlRequest,
) (*pb.GetOldByEmailAndWlResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "profile: get old by email and wl: parse wl id %s: %v", req.WhiteLabelID, err)
		return nil, fmt.Errorf("profile: get old by email and wl: parse wl id: %w", err)
	}
	pr, err := p.profileService.GetOldByEmailAndWl(ctx, req.Email, wlID)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	if pr == nil {
		return &pb.GetOldByEmailAndWlResponse{}, nil
	}
	return &pb.GetOldByEmailAndWlResponse{
		Profile: &pb.Profile{
			User: p.convertToProtoUser(pr.User),
		},
	}, nil
}

func (p *Profile) GetOldByEmailAndWlV2(
	ctx context.Context,
	req *pb.GetOldByEmailAndWlV2Request,
) (*pb.GetOldByEmailAndWlV2Response, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "profile: get old by email and wl: parse wl id %s: %v", req.WhiteLabelID, err)
		return nil, fmt.Errorf("profile: get old by email and wl: parse wl id: %w", err)
	}
	pr, err := p.profileService.GetOldByEmailAndWl(ctx, req.Email, wlID)
	if err != nil {
		err = fmt.Errorf("GetOldByEmailAndWlV2: %w", err)
		log.Error(ctx, err.Error())
		return nil, err
	}
	if pr == nil {
		return &pb.GetOldByEmailAndWlV2Response{}, nil
	}
	return &pb.GetOldByEmailAndWlV2Response{
		Profile: &pb.ProfileV3{
			User: p.convertToProtoUser(pr.User),
		},
	}, nil
}

func (p *Profile) GetSuspended(ctx context.Context, req *pb.GetSuspendedRequest) (*pb.GetSuspendedResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "getSuspended: parse: %s", err.Error())
		return nil, fmt.Errorf("getSuspended: parse: %w", err)
	}

	suspended, err := p.profileService.GetSuspended(ctx, userID)
	if err != nil {
		log.Error(ctx, "profile: %s", err.Error())
		return nil, fmt.Errorf("profile: %w", err)
	}
	return &pb.GetSuspendedResponse{
		Suspended: suspended,
	}, nil
}

func (p *Profile) SetSuspended(ctx context.Context, req *pb.SetSuspendedRequest) (*pb.SetSuspendedResponse, error) {
	ids := make([]uuid.UUID, len(req.UserIds))
	for i, id := range req.UserIds {
		uid, err := uuid.Parse(id)
		if err != nil {
			log.Error(ctx, "setSuspended: parse: %s", err.Error())
			return nil, fmt.Errorf("setSuspended: parse: %w", err)
		}
		ids[i] = uid
	}

	err := p.profileService.SetSuspended(ctx, ids, req.Suspended)
	if err != nil {
		log.Error(ctx, "profile: %s", err.Error())
		return nil, fmt.Errorf("profile: %w", err)
	}
	return &pb.SetSuspendedResponse{}, nil
}

func (p *Profile) ChangeWalletAddress(
	ctx context.Context,
	req *pb.ChangeWalletAddressRequest,
) (*pb.ChangeWalletAddressResponse, error) {
	err := p.profileService.SendEmailToChangeAddress(
		ctx,
		int(req.UserId),
		req.Username,
		req.Ip,
		req.Coin,
		req.Address,
		req.Domain,
	)
	if err != nil {
		log.Error(ctx, "changeAddress: %s", err.Error())
		switch {
		case errors.Is(err, service.ErrAddressChangeNotAllowed):
			return nil, errAddressChangeNotAllowed
		case errors.Is(err, service.ErrBothAddressIsTheSame):
			return nil, errBothAddressIsTheSame
		}
		return nil, fmt.Errorf("changeAddress: %w", err)
	}

	return &pb.ChangeWalletAddressResponse{}, nil
}

func (p *Profile) UpdateMinPay(ctx context.Context, req *pb.UpdateMinPayRequest) (*pb.UpdateMinPayResponse, error) {
	err := p.profileService.UpdateMinPay(ctx, int(req.UserId), req.Coin, req.Value)
	if err != nil {
		log.Error(ctx, "updateMinPay: %s", err.Error())
		switch {
		case errors.Is(err, service.ErrMinPayNotValid):
			return nil, errMinPayNotValid
		}
		return nil, fmt.Errorf("updateMinPay: %w", err)
	}
	return &pb.UpdateMinPayResponse{}, nil
}

func (p *Profile) ChangeWalletAddressConfirm(
	ctx context.Context,
	req *pb.ChangeWalletAddressConfirmRequest,
) (*pb.ChangeWalletAddressConfirmResponse, error) {
	res, err := p.profileService.ChangeWalletAddressConfirm(ctx, req.UserId, req.Token)
	if err != nil {
		log.Error(ctx, err.Error())
		switch {
		case errors.Is(err, service.ErrCoinNotFound):
			return nil, errCoinNotFound
		case errors.Is(err, service.ErrTokenNotFound):
			return nil, errTokenNotFound
		}
		return nil, err
	}
	return &pb.ChangeWalletAddressConfirmResponse{
		Address: res.Address,
		UserId:  res.UserID,
		CoinId:  res.CoinID,
	}, nil
}

func (p *Profile) UpdateUserIsActive(
	ctx context.Context,
	req *pb.UpdateUserIsActiveRequest,
) (*pb.UpdateUserIsActiveResponse, error) {
	err := p.profileService.UpdateUserIsActive(ctx, req.Email, req.Active)
	if err != nil {
		log.Error(ctx, "updateUserIsActive: %s", err.Error())
		return nil, err
	}
	return &pb.UpdateUserIsActiveResponse{}, nil
}

func (p *Profile) GetUserIsActive(ctx context.Context, req *pb.GetUserIsActiveRequest) (*pb.GetUserIsActiveResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "getUserIsActive: parse: %s", err.Error())
		return nil, fmt.Errorf("getUserIsActive: parse: %w", err)
	}
	isActive, err := p.profileService.GetUserIsActive(ctx, userID)
	if err != nil {
		log.Error(ctx, "getUserIsActive: %s", err.Error())
		return nil, fmt.Errorf("getUserIsActive: %w", err)
	}
	return &pb.GetUserIsActiveResponse{IsActive: isActive}, nil
}

func (p *Profile) GetKycStatus(ctx context.Context, req *pb.GetKycStatusRequest) (*pb.GetKycStatusResponse, error) {
	kyc, err := p.profileService.GetKycStatus(ctx, int(req.UserId))
	if err != nil {
		log.Error(ctx, "getKycStatus: %s", err.Error())
		return nil, fmt.Errorf("getKycStatus: %w", err)
	}
	return &pb.GetKycStatusResponse{
		RetryAfter: &timestamppb.Timestamp{
			Seconds: kyc.RetryAfter.Unix(),
		},
		DelayMinutes: int32(kyc.DelayMinutes),
		Status:       pb.KycStatus(kyc.Status),
		IsAllowed:    kyc.IsAllowed,
		Overall:      kyc.Overall,
		DocCheck:     kyc.DocCheck,
		FaceCheck:    kyc.FaceCheck,
	}, nil
}

func (p *Profile) SetKycStatus(ctx context.Context, req *pb.SetKycStatusRequest) (*pb.SetKycStatusResponse, error) {
	err := p.profileService.SetKycStatus(ctx, int(req.UserId), int(req.Status))
	if err != nil {
		log.Error(ctx, "setKycStatus: %s", err.Error())
		return nil, fmt.Errorf("setKycStatus: %w", err)
	}
	return &pb.SetKycStatusResponse{}, nil
}

func (p *Profile) InsertKycHistory(ctx context.Context, req *pb.InsertKycHistoryRequest) (*pb.InsertKycHistoryResponse, error) {
	err := p.profileService.InsertKycHistory(ctx, int(req.UserId), req.Data)
	if err != nil {
		log.Error(ctx, "insertKycHistory: %s", err.Error())
		return nil, fmt.Errorf("insertKycHistory: %w", err)
	}
	return &pb.InsertKycHistoryResponse{}, nil
}

func (p *Profile) CheckAppleAccount(
	ctx context.Context,
	req *pb.CheckAppleAccountRequest,
) (*pb.CheckAppleAccountResponse, error) {
	registrationRequired, email, err := p.profileService.CheckAppleAccount(ctx, req.AppleId, req.Email)
	if err != nil {
		log.Error(ctx, "checkAppleAccount: %s", err.Error())
		return nil, err
	}

	return &pb.CheckAppleAccountResponse{
		RegistrationRequired: registrationRequired,
		Email:                email,
	}, nil
}

func (p *Profile) GetNotificationSettings(
	ctx context.Context,
	req *pb.GetNotificationSettingsRequest,
) (*pb.GetNotificationSettingsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "getNotificationSettings: parse user id: %s: %v", req.UserId, err)
		return nil, err
	}
	settings, err := p.profileService.GetNotificationSettings(ctx, userID)
	if err != nil {
		log.Error(ctx, "getNotificationSettings: %v", err)
		return nil, err
	}
	return &pb.GetNotificationSettingsResponse{
		Settings: toProtoNotificationsSettings(settings),
	}, nil
}

func toProtoNotificationsSettings(settings *model.NotificationSettings) *pb.NotificationSettings {
	if settings == nil {
		return nil
	}
	return &pb.NotificationSettings{
		Email:                  settings.Email,
		Language:               settings.Language,
		IsTgNotificationsOn:    settings.IsEmailNotificationsOn,
		TgId:                   int64(settings.TgID),
		WhitelabelId:           settings.WhiteLabelID.String(),
		IsEmailNotificationsOn: settings.IsEmailNotificationsOn,
		IsPushNotificationsOn:  settings.IsPushNotificationsOn,
	}
}

func (p *Profile) SaveNotificationSettings(
	ctx context.Context,
	req *pb.SaveNotificationSettingsRequest,
) (*pb.SaveNotificationSettingsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "saveNotificationSettings: parse user id: %s: %v", req.UserId, err)
		return nil, err
	}
	err = p.profileService.SaveNotificationSettings(ctx, &model.ChangeableNotificationSettings{
		UserID:                 userID,
		IsTgNotificationsOn:    req.IsTgNotificationsOn,
		IsEmailNotificationsOn: req.IsEmailNotificationsOn,
		TgID:                   req.TgId,
		IsPushNotificationsOn:  req.IsPushNotificationsOn,
	})
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.SaveNotificationSettingsResponse{}, nil
}

func (p *Profile) RelatedUsers(ctx context.Context, req *pb.RelatedUsersRequest) (*pb.RelatedUsersResponse, error) {
	firstID, err := uuid.Parse(req.FirstId)
	if err != nil {
		log.Error(ctx, "relatedAccounts: parse first_id: %s", err.Error())
		return &pb.RelatedUsersResponse{}, fmt.Errorf("parse first_id: %w", err)
	}
	secondID, err := uuid.Parse(req.SecondId)
	if err != nil {
		log.Error(ctx, "relatedAccounts: parse second_id: %s", err.Error())
		return &pb.RelatedUsersResponse{}, fmt.Errorf("parse second_id: %w", err)
	}
	related, err := p.profileService.RelatedUsers(ctx, firstID, secondID)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, fmt.Errorf("relatedAccounts: %s", err.Error())
	}
	return &pb.RelatedUsersResponse{
		Related: related,
	}, nil
}

func (p *Profile) GetAllSubUsers(ctx context.Context, req *pb.GetAllSubUsersRequest) (*pb.GetAllSubUsersResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "getAllSubUsers: parse: %s", err.Error())
		return nil, fmt.Errorf("getAllSubUsers: parse: %w", err)
	}
	users, err := p.profileService.GetAllSubUsers(ctx, userID)
	if err != nil {
		log.Error(ctx, "getAllSubUsers: %s", err.Error())
		return nil, err
	}
	subs := make([]*pb.GetAllSubUsersResponse_Subs, len(users))
	for i, u := range users {
		subs[i] = &pb.GetAllSubUsersResponse_Subs{
			UserId:    u.ID.String(),
			Username:  u.Username,
			UserOldId: u.OldID,
		}
	}
	return &pb.GetAllSubUsersResponse{
		Subs: subs,
	}, nil
}

func (p *Profile) GetAllUserIDsByUsername(
	ctx context.Context,
	req *pb.GetAllUserIDsByUsernameRequest,
) (*pb.GetAllUserIDsByUsernameResponse, error) {
	ids, err := p.profileService.GetAllUserIDsByUsername(ctx)
	if err != nil {
		log.Error(ctx, "getAllUserIDsByUsername: %s", err.Error())
		return nil, fmt.Errorf("getAllUserIDsByUsername: %w", err)
	}
	m := make(map[string]string, len(ids))
	for k, v := range ids {
		m[k] = v.String()
	}
	return &pb.GetAllUserIDsByUsernameResponse{
		Ids: m,
	}, nil
}

func (p *Profile) GetReferrals(ctx context.Context, req *pb.GetReferralsRequest) (*pb.GetReferralsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "getReferrals: parse: %s", err.Error())
		return nil, fmt.Errorf("GetReferrals: parse: %w", err)
	}

	referrals, count, err := p.profileService.GetReferrals(
		ctx,
		userID,
		int(req.Skip),
		int(req.Take),
		req.Sort.Field,
		req.Sort.Asc,
	)
	if err != nil {
		log.Error(ctx, "getReferrals: %s", err.Error())
		return nil, fmt.Errorf("getReferrals: %w", err)
	}

	refs := make([]*pb.Referral, len(referrals))
	for i, r := range referrals {
		refs[i] = &pb.Referral{
			Id:        r.ID.String(),
			Username:  r.Username,
			Email:     r.Email,
			CreatedAt: timestamppb.New(r.CreatedAt),
		}
	}

	return &pb.GetReferralsResponse{
		Referrals:  refs,
		TotalCount: int32(count),
	}, nil
}

func (p *Profile) GetUsernamesByIDs(
	ctx context.Context,
	req *pb.GetUsernamesByIDsRequest,
) (*pb.GetUsernamesByIDsResponse, error) {
	ids := make([]uuid.UUID, len(req.UserIds))
	for i := range ids {
		id, err := uuid.Parse(req.UserIds[i])
		if err != nil {
			log.Error(ctx, "getUsernamesByIds: parse user id: %s: %v", req.UserIds[i], err)
			return nil, err
		}
		ids = append(ids, id)
	}
	usernames, err := p.profileService.GetUsernamesByIDs(ctx, ids)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	usernamesProto := make(map[string]string)
	for k, v := range usernames {
		usernamesProto[k.String()] = v
	}
	return &pb.GetUsernamesByIDsResponse{
		Usernames: usernamesProto,
	}, nil
}

func (p *Profile) GetEmailsByIDs(ctx context.Context, req *pb.GetEmailsByIDsRequest) (*pb.GetEmailsByIDsResponse, error) {
	ids := make([]uuid.UUID, len(req.UserIds))
	for i := range ids {
		id, err := uuid.Parse(req.UserIds[i])
		if err != nil {
			log.Error(ctx, "getEmailByIds: parse user id: %s: %v", req.UserIds[i], err)
			return nil, err
		}
		ids = append(ids, id)
	}
	emails, err := p.profileService.GetEmailsByIDs(ctx, ids)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	emailsProto := make(map[string]string)
	for k, v := range emails {
		emailsProto[k.String()] = v
	}
	return &pb.GetEmailsByIDsResponse{
		Emails: emailsProto,
	}, nil
}

func (p *Profile) SafeDeleteByID(ctx context.Context, req *pb.SafeDeleteByIDRequest) (*pb.SafeDeleteByIDResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "safeDeleteByID: parse: %s", err.Error())
		return nil, fmt.Errorf("SafeDeleteByID: parse: %w", err)
	}

	err = p.profileService.SafeDeleteByID(ctx, userID)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.SafeDeleteByIDResponse{}, nil
}

func (p *Profile) SaveV4(ctx context.Context, req *pb.SaveV4Request) (*pb.SaveV4Response, error) {
	user, err := p.parseProtoUserV2(req.User)
	if err != nil {
		log.Error(ctx, "SaveV4: %v", err)
		return nil, err
	}
	id, err := p.profileService.SaveV4(ctx, user)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.SaveV4Response{
		UserId: id,
	}, nil
}

func (p *Profile) UpdateRefID(ctx context.Context, req *pb.UpdateRefIDRequest) (*pb.UpdateRefIDResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("update ref_id: parse user id %v: %w", req.UserID, err)
	}
	newRefID, err := p.profileService.GetUserIDByOldID(ctx, req.RefID)
	if err != nil {
		return nil, fmt.Errorf("update ref_id: get new ref id %v: %w", req.RefID, err)
	}
	if err := p.profileService.UpdateRefID(ctx, req.OldUserID, userID, req.RefID, newRefID); err != nil {
		return nil, fmt.Errorf("update ref_id: %w", err)
	}
	return &pb.UpdateRefIDResponse{}, nil
}

func (p *Profile) parseProtoUserV2(u *pb.UserV2) (*model.User, error) {
	userID, err := uuid.Parse(u.ID)
	if err != nil {
		return nil, fmt.Errorf("parse user_id %v: %w", u.ID, err)
	}
	wlID, err := uuid.Parse(u.WhiteLabelID)
	if err != nil {
		return nil, fmt.Errorf("parse whiteLabelID %v: %w", u.WhiteLabelID, err)
	}
	var parentID uuid.UUID
	if u.ParentId != "" {
		parentID, err = uuid.Parse(u.ParentId)
		if err != nil {
			return nil, fmt.Errorf("parse parent_id %v: %w", u.ParentId, err)
		}
	}
	var newRefID uuid.UUID
	if u.NewRefId != "" {
		newRefID, err = uuid.Parse(u.NewRefId)
		if err != nil {
			return nil, fmt.Errorf("parse new_ref_id: %w %s", err, u.NewRefId)
		}
	}

	poolType := defaultPoolType
	if u.PoolType != nil {
		poolType = *u.PoolType
		// ? looks strange, but *bug*
		// if client using old library version it pass empty string instead of nil pointer
		if len(poolType) == 0 {
			poolType = defaultPoolType
		}
	}
	return &model.User{
		ID:           userID,
		Username:     u.Username,
		Vip:          u.Vip,
		SegmentID:    int(u.SegmentID),
		RefID:        int(u.RefID),
		ParentID:     parentID,
		Email:        u.Email,
		Password:     u.Password,
		CreatedAt:    u.CreatedAt.AsTime(),
		WhiteLabelID: wlID,
		ApiKey:       u.ApiKey,
		IsActive:     u.IsActive,
		AppleID:      u.AppleId,
		NewRefID:     newRefID,
		PoolType:     poolType,
		Language:     u.Language,
		TgID:         u.TgId,
		TgUsername:   u.TgUsername,
	}, nil
}

func (p *Profile) GetFlagReferralLinkGenerated(
	ctx context.Context,
	req *pb.GetFlagReferralLinkGeneratedRequest,
) (*pb.GetFlagReferralLinkGeneratedResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error(ctx, "profile: get flag referral link generated %v: %v", req.UserID, err)
		return nil, fmt.Errorf("profile: get flag referral link generated: %w", err)
	}
	flg, err := p.profileService.WasReferralLinkGenerated(ctx, userID)
	if err != nil {
		log.Error(ctx, "profile: get flag referral link generated %s", err.Error())
		return nil, fmt.Errorf("profile: get flag referral link generated%w", err)
	}
	return &pb.GetFlagReferralLinkGeneratedResponse{Value: flg}, nil
}

func (p *Profile) SetFlagReferralLinkGenerated(
	ctx context.Context,
	req *pb.SetFlagReferralLinkGeneratedRequest,
) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error(ctx, "profile: set flag referral link generated %v: %v", req.UserID, err)
		return nil, fmt.Errorf("profile: set flag referral link generated: %w", err)
	}
	err = p.profileService.SetFlagReferralLinkGenerated(ctx, userID, true)
	if err != nil {
		log.Error(ctx, "profile: get flag referral link generated %s", err.Error())
		return nil, fmt.Errorf("profile: get flag referral link generated%w", err)
	}
	return &emptypb.Empty{}, nil
}

func (p *Profile) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	user, err := p.parseProtoUserV2(req.User)
	if err != nil {
		log.Error(ctx, "UpdateProfile: parse user: %v", err)
		return nil, err
	}
	err = p.profileService.UpdateProfile(ctx, user)
	if err != nil {
		log.Error(ctx, "UpdateProfile: update: %v", err)
		return nil, err
	}
	return nil, nil
}

func (p *Profile) SetTimezone(ctx context.Context, req *pb.SetTimezoneRequest) (*pb.SetTimezoneResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "SetTimezone: parse user id: %s. %v", req.UserId, err)
		return nil, fmt.Errorf("SetTimezone: parse user id: %w", err)
	}

	err = p.profileService.SetTimezone(ctx, userID, req.Timezone)
	if err != nil {
		log.Error(ctx, fmt.Errorf("SetTimezone: %w", err).Error())
		return nil, fmt.Errorf("SetTimezone: %w", err)
	}

	return &pb.SetTimezoneResponse{}, nil
}

func (p *Profile) SetLanguage(ctx context.Context, req *pb.SetLanguageRequest) (*pb.SetLanguageResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "SetLanguage: parse user id: %s. %v", req.UserId, err)
		return nil, fmt.Errorf("SetLanguage: parse user id: %w", err)
	}

	err = p.profileService.SetLanguage(ctx, userID, req.Language)
	if err != nil {
		log.Error(ctx, fmt.Errorf("SetLanguage: %w", err).Error())
		return nil, fmt.Errorf("SetLanguage: %w", err)
	}

	return &pb.SetLanguageResponse{}, nil
}

func (p *Profile) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "GetAddresses: parse user id: %s. %v", req.UserId, err)
		return nil, fmt.Errorf("GetAddresses: parse user id: %w", err)
	}

	addresses, err := p.profileService.GetAddresses(ctx, userID)
	if err != nil {
		log.Error(ctx, fmt.Errorf("GetAddresses: %w", err).Error())
		return nil, fmt.Errorf("GetAddresses: %w", err)
	}
	addrs := make([]*pb.Address, len(addresses))
	for i := range addrs {
		addrs[i] = &pb.Address{
			Coin:          addresses[i].Coin,
			Minpay:        addresses[i].MinPay.String(),
			WalletAddress: addresses[i].WalletAddress,
			MiningAddress: addresses[i].MiningAddress,
		}
	}
	return &pb.GetAddressesResponse{
		Addresses: addrs,
	}, nil
}

func (p *Profile) GetUsersWithWL(ctx context.Context, req *pb.GetUsersWithWLRequest) (*pb.GetUsersWithWLResponse, error) {
	wlID, err := uuid.Parse(req.GetWlUuid())
	if err != nil {
		log.Error(ctx, "profile: GetUsersWithWL %v: %v", req.GetWlUuid(), err)
		return nil, fmt.Errorf("GetUsersWithWL: %w", err)
	}
	if req.GetLimit() == 0 {
		log.Error(ctx, "profile: GetUsersWithWL %v: limit is zero", req.GetWlUuid())
		return nil, fmt.Errorf("profile: GetUsersWithWL %v: limit is zero", req.GetWlUuid())
	}
	userNames, count, err := p.profileService.GetCountUserWithWL(ctx, wlID, req.GetOffset(), req.GetLimit())
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.GetUsersWithWLResponse{
		Users: toProtoUsers(userNames),
		Total: count,
	}, nil
}

func (p *Profile) GetUsersByUUIDs(ctx context.Context, req *pb.GetUsersByUUIDsRequest) (*pb.GetUsersByUUIDsResponse, error) {
	rawUUIDs := req.GetUsersUuids()
	uuids := make([]uuid.UUID, len(rawUUIDs))
	for i := range uuids {
		id, err := uuid.Parse(rawUUIDs[i])
		if err != nil {
			log.Error(ctx, "GetUserByUUIDs: parse user id: %s: %v", rawUUIDs[i], err)
			return nil, err
		}
		uuids = append(uuids, id)
	}
	users, err := p.profileService.GetUsersByUUIDs(ctx, uuids)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &pb.GetUsersByUUIDsResponse{
		Users: toProtoUsers(users),
	}, nil
}

func toProtoUsers(users []model.UserShortInfo) []*pb.ShortUserInfo {
	userInfos := make([]*pb.ShortUserInfo, len(users))
	for i := range users {
		userInfos[i] = &pb.ShortUserInfo{
			UserUuid:  users[i].UserUUID,
			UserName:  users[i].UserName,
			Email:     users[i].Email,
			CreatedAt: timestamppb.New(users[i].CreatedAt),
		}
	}

	return userInfos
}

func (p *Profile) GetByUsernamesForReferrals(
	req *pb.GetByUsernamesForReferralsRequest,
	srv pb.ProfileService_GetByUsernamesForReferralsServer,
) error {
	pr, err := p.profileService.GetByUsernamesForReferrals(srv.Context(), req.Usernames)
	if err != nil {
		log.Error(srv.Context(), "%v. usernames: %s", err, req.Usernames)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return errInconsistentUserInfo
		}
		return err
	}
	if pr == nil {
		return nil
	}

	batchSize := 1000
	for i := 0; i < len(pr); i += batchSize {
		end := i + batchSize
		if end > len(pr) {
			end = len(pr)
		}

		profiles := make([]*pb.Profile, 0, len(pr[i:end]))
		for _, profile := range pr[i:end] {
			profiles = append(profiles, &pb.Profile{
				User: p.convertToProtoUser(profile.User),
			})
		}

		msg := &pb.GetByUsernamesForReferralsResponse{
			Profile: profiles,
		}

		if err := srv.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func (p *Profile) CreateProfile(ctx context.Context, in *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateProfile: email is required")
	}

	if in.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateProfile: username is required")
	}

	wlUUID := uuid.Nil
	var err error
	if in.GetWhiteLabelUuid() != "" {
		wlUUID, err = uuid.Parse(in.GetWhiteLabelUuid())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "CreateProfile: white label uuid is invalid")
		}
	}

	resp, err := p.profileService.CreateProfile(ctx, service.CreateProfileRequest{
		Username:         in.GetUsername(),
		Email:            in.GetEmail(),
		Password:         in.GetPassword(),
		RefId:            in.GetRefId(),
		RefUuid:          in.GetRefUuid(),
		WhiteLabelUuid:   wlUUID,
		AppleId:          in.GetAppleId(),
		Language:         in.GetLanguage(),
		TelegramId:       in.GetTelegramId(),
		TelegramUserName: in.GetTelegramUserName(),
	})
	if err != nil {
		log.Error(ctx, "CreateProfile: %v", err)
		return nil, err
	}

	return &pb.CreateProfileResponse{
		UserUuid: resp.UserUUID.String(),
		UserId:   resp.UserID,
	}, nil
}

func (p *Profile) CreateSubUser(ctx context.Context, in *pb.CreateSubUserRequest) (*pb.CreateSubUserResponse, error) {
	parentUUID, err := uuid.Parse(in.GetParentUserUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userUUID, id, err := p.profileService.CreateSubUser(ctx,
		parentUUID,
		in.GetUsername(),
		p.toCoinAndAddresses(in.GetAddresses()))
	if err != nil {
		return nil, err
	}

	return &pb.CreateSubUserResponse{
		NewUserUUID: userUUID.String(),
		OldID:       id,
	}, nil
}

func (p *Profile) toCoinAndAddresses(in []*pb.CoinAndAddress) []*model.CoinAndAddress {
	res := make([]*model.CoinAndAddress, 0)
	for _, addr := range in {
		res = append(res, &model.CoinAndAddress{
			Coin:    addr.Coin,
			Address: addr.Address,
		})
	}
	return res
}

func (p *Profile) GetUserByTg(
	ctx context.Context,
	req *pb.GetUserByTgRequest,
) (*pb.GetUserByTgResponse, error) {
	u, err := p.profileService.GetUserByTg(ctx, req.TgID)
	if err != nil {
		log.Error(ctx, "GetUserByTg %v. tg: %s", err, req.TgID)
		switch {
		case errors.Is(err, service.ErrInconsistentUserInfo):
			return nil, errInconsistentUserInfo
		}
		return nil, err
	}
	return &pb.GetUserByTgResponse{
		User: p.convertToProtoUser(u),
	}, nil
}

func (p *Profile) SetUserAttributes(ctx context.Context, r *profile.SetUserAttributesRequest) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(r.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var whiteLabelID uuid.UUID
	if r.GetWhiteLabelID() != "" {
		whiteLabelID, err = uuid.Parse(r.GetWhiteLabelID())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	var parentID uuid.UUID
	if r.GetParentId() != "" {
		parentID, err = uuid.Parse(r.GetParentId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	attrs := model.UserAttributes{
		Username: r.Username,
		ParentID: func() *uuid.UUID {
			if parentID == uuid.Nil {
				return nil
			}
			return &parentID
		}(),
		Language: r.Language,
		WhiteLabelID: func() *uuid.UUID {
			if whiteLabelID == uuid.Nil {
				return nil
			}
			return &whiteLabelID
		}(),
		PoolType:                 r.PoolType,
		WasReferralLinkGenerated: r.WasReferralLinkGenerated,
		IsAmbassador:             r.IsAmbassador,
	}

	err = p.profileService.SetUserAttributes(ctx, userID, attrs)
	if err != nil {
		return nil, fmt.Errorf("SetUserAttributes: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (p *Profile) GetOldIDByID(ctx context.Context, req *profile.GetOldIDByIDRequest) (*profile.GetOldIDByIDResponse, error) {
	userUUID, err := uuid.Parse(req.Id)
	if err != nil {
		err = fmt.Errorf("GetOldIDByID: parse new id %w", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	oldNewIDs, err := p.oldIDResolver.GetIDs(ctx, []uuid.UUID{userUUID})
	if err != nil {
		return nil, fmt.Errorf("GetIDByOldID: %w", err)
	}
	if len(oldNewIDs) == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &profile.GetOldIDByIDResponse{
		NewId: userUUID.String(),
		OldId: oldNewIDs[0].Old,
	}, nil
}
