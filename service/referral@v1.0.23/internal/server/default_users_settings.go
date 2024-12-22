// Package server implements a simple web server.
package server

import (
	"context"
	"errors"
	"strings"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/service"
	pb "code.emcdtech.com/emcd/service/referral/protocol/default_users_settings"
)

const (
	oneHundred float64 = 100
)

type DefaultUsersSettingsServer struct {
	pb.UnimplementedDefaultUsersSettingsServiceServer
	service *service.DefaultUsersSettings
}

func NewDefaultUsersSettings(s *service.DefaultUsersSettings) *DefaultUsersSettingsServer {
	return &DefaultUsersSettingsServer{
		service: s,
	}
}

func (s *DefaultUsersSettingsServer) CreateUsersSettings(ctx context.Context, req *pb.CreateUsersSettingsRequest) (*pb.CreateUsersSettingsResponse, error) {
	users := req.GetUsers()
	if len(users) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user list is required")
	}

	createRequest, err := toReferralSettings(users)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.CreateUsersSettings(ctx, createRequest)
	if err != nil {
		log.Error(ctx, "CreateUsersSettings: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUsersSettingsResponse{}, nil
}
func (s *DefaultUsersSettingsServer) UpdateUsersSettings(ctx context.Context, req *pb.UpdateUsersSettingsRequest) (*pb.UpdateUsersSettingsResponse, error) {
	users := req.GetUsers()
	if len(users) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user list is required")
	}

	updateRequest, err := toReferralSettings(users)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updateMode, err := model.NewUpdateMode(req.GetUpdateMode())

	err = s.service.UpdateUsersSettings(ctx, updateRequest, updateMode)
	if err != nil {
		log.Error(ctx, "UpdateUsersSettings: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateUsersSettingsResponse{}, nil
}
func (s *DefaultUsersSettingsServer) GetUsersSettings(
	req *pb.GetUsersSettingsRequest,
	srv grpc.ServerStreamingServer[pb.GetUsersSettingsResponse]) error {
	userRaws := req.GetUserUuids()
	if len(userRaws) == 0 || len(userRaws) > 100 {
		return status.Error(codes.InvalidArgument, "user list is required and cannot be more than 100")
	}

	uuids := make([]uuid.UUID, 0, len(userRaws))
	for _, rawUUID := range userRaws {
		userUUID, err := uuid.Parse(rawUUID)
		if err != nil {
			return err
		}
		uuids = append(uuids, userUUID)
	}

	ctx := srv.Context()

	settings, err := s.service.GetUsersSettings(ctx, uuids, req.GetProducts(), req.GetCoins())
	if err != nil {
		log.Error(ctx, "GetUsersSettings: %v", err)
		return status.Error(codes.Internal, err.Error())
	}

	batchSize := 1000
	usersSettings := make([]*pb.UserSettings, 0, batchSize)
	for _, referralSettings := range settings {
		settingsPB := make([]*pb.UserPreference, 0, len(referralSettings.Preferences))
		for _, preference := range referralSettings.Preferences {
			settingsPB = append(settingsPB, &pb.UserPreference{
				Product:     preference.Product,
				Coin:        preference.Coin,
				Fee:         preference.Fee,
				ReferralFee: preference.ReferralFee,
				CreatedAt:   timestamppb.New(preference.CreatedAt),
				UpdatedAt:   timestamppb.New(preference.UpdatedAt),
			})
		}

		usersSettings = append(usersSettings, &pb.UserSettings{
			UserUuid: referralSettings.ReferralUUID.String(),
			Settings: settingsPB,
		})

		if len(usersSettings) >= batchSize {
			if err := srv.Send(&pb.GetUsersSettingsResponse{
				Users: usersSettings,
			}); err != nil {
				return err
			}
			usersSettings = make([]*pb.UserSettings, 0, batchSize)
		}
	}

	if len(usersSettings) > 0 {
		if err := srv.Send(&pb.GetUsersSettingsResponse{
			Users: usersSettings,
		}); err != nil {
			return err
		}
	}

	return nil
}

func toReferralSettings(pbIn []*pb.UserSettings) ([]model.ReferralSettings, error) {
	list := make([]model.ReferralSettings, 0, len(pbIn))
	for _, user := range pbIn {
		userUUID, err := uuid.Parse(user.GetUserUuid())
		if err != nil {
			return nil, err
		}

		preferences, err := toPreferences(user.GetSettings())
		if err != nil {
			return nil, err
		}

		list = append(list, model.ReferralSettings{
			ReferralUUID: userUUID,
			Preferences:  preferences,
		})
	}
	return list, nil
}

func toPreferences(pbIn []*pb.UserPreference) ([]model.ReferralPreference, error) {
	list := make([]model.ReferralPreference, 0, len(pbIn))

	for _, preference := range pbIn {
		if preference.GetProduct() == "" {
			return nil, errors.New("product field is required")
		}
		if preference.GetCoin() == "" {
			return nil, errors.New("coin field is required")
		}

		fee := preference.GetFee()
		if fee > 100 || fee < 0 {
			return nil, errors.New("fee field is required and cannot be more than 100 or less than 0")
		}

		referralFee := preference.GetReferralFee()
		if referralFee > oneHundred {
			return nil, errors.New("referral fee is greater than 100")
		}

		list = append(list, model.ReferralPreference{
			Product:     strings.ToLower(preference.GetProduct()),
			Coin:        strings.ToUpper(preference.GetCoin()),
			Fee:         fee,
			ReferralFee: referralFee,
		})
	}
	return list, nil
}
