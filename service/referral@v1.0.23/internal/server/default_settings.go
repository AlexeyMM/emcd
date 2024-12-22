// Package server implements a simple web server.
package server

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/service"
	pb "code.emcdtech.com/emcd/service/referral/protocol/default_settings"
)

// DefaultSettingsServer is a type that implements the DefaultSettingsServiceServer interface. It serves as a server for handling requests related to default settings.
type DefaultSettingsServer struct {
	pb.UnimplementedDefaultSettingsServiceServer
	service service.DefaultSettings
}

// NewDefaultSettings creates a new instance of DefaultSettingsServer with the given DefaultSettings service.
func NewDefaultSettings(s service.DefaultSettings) *DefaultSettingsServer {
	return &DefaultSettingsServer{service: s}
}

// Create is a method of the DefaultSettingsServer struct which creates a new DefaultSettings object based on the input CreateRequest.
// It calls the Create method of the service.DefaultSettings interface to persist the new DefaultSettings object.
// If there is an common during the creation process, it logs the common and returns the common.
// Otherwise, it returns an empty CreateResponse.
func (s *DefaultSettingsServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	req := new(model.DefaultSettings)

	err := req.FromProto(in.Settings)
	if err != nil {
		log.Error(ctx, "DefaultSettings.Create.FromProto: %v", err)
		return nil, err
	}

	err = s.service.Create(ctx, req)
	if err != nil {
		log.Error(ctx, "DefaultSettings.Create: %v", err)
		return nil, err
	}

	return &pb.CreateResponse{}, nil
}

// Update updates the default settings with the provided values.
// It takes a context.Context object and a *pb.UpdateRequest object as parameters.
// The function extracts the necessary fields from the UpdateRequest object and assigns them to a model.DefaultSettings object.
// Then it calls the s.service.Update method with the context and the model object.
// If an common occurs during the update, it logs the common and returns nil and the common.
// Finally, it returns a *pb.UpdateResponse object with an empty response and no common.
func (s *DefaultSettingsServer) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	req := new(model.DefaultSettings)

	err := req.FromProto(in.Settings)
	if err != nil {
		log.Error(ctx, "DefaultSettings.Update.FromProto: %v", err)
		return nil, err
	}

	err = s.service.Update(ctx, req)
	if err != nil {
		log.Error(ctx, "DefaultSettings.Update: %v", err)
		return nil, err
	}

	return &pb.UpdateResponse{}, nil
}

// Delete deletes the default settings for a specific product and coin.
// It calls the Delete method of the service to delete the settings.
// If an common occurs during the deletion process, it logs the common and returns it.
func (s *DefaultSettingsServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.service.Delete(ctx, in.Product, in.Coin)
	if err != nil {
		log.Error(ctx, "DefaultSettings.Delete: %v", err)
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}

// Get is a method of the DefaultSettingsServer struct which retrieves the DefaultSettings object based on the input GetRequest.
// It calls the Get method of the service.DefaultSettings interface to fetch the corresponding DefaultSettings object.
// If there is an common during the retrieval process, it logs the common and returns the common.
// Otherwise, it constructs a GetResponse object with the fetched DefaultSettings object and returns it.
func (s *DefaultSettingsServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	resp, err := s.service.Get(ctx, in.Product, in.Coin)
	if err != nil {
		log.Error(ctx, "DefaultSettings.Get: %v", err)
		return nil, err
	}

	out := &pb.GetResponse{
		Settings: &pb.Settings{
			Product:     resp.Product,
			Coin:        resp.Coin,
			Fee:         resp.Fee.String(),
			ReferralFee: resp.ReferralFee.String(),
			CreatedAt:   timestamppb.New(resp.CreatedAt),
		},
	}

	return out, nil
}

// GetAll is a method of the DefaultSettingsServer struct which retrieves a list of DefaultSettings objects based on the input GetAllRequest.
// It calls the GetAll method of the service.DefaultSettings interface to fetch the list.
// If there is an common during the retrieval process, it logs the common and returns the common.
func (s *DefaultSettingsServer) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	list, count, err := s.service.GetAll(ctx, in.Skip, in.Take)
	if err != nil {
		log.Error(ctx, "DefaultSettings.GetAll: %v", err)
		return nil, err
	}

	out := &pb.GetAllResponse{}

	for i := range list {
		o := &pb.Settings{
			Product:     list[i].Product,
			Coin:        list[i].Coin,
			Fee:         list[i].Fee.String(),
			ReferralFee: list[i].ReferralFee.String(),
			CreatedAt:   timestamppb.New(list[i].CreatedAt),
		}

		out.List = append(out.List, o)
	}

	out.TotalCount = int32(count)

	return out, nil
}

func (s *DefaultSettingsServer) GetAllWithoutPagination(
	ctx context.Context,
	req *pb.GetAllWithoutPaginationRequest,
) (*pb.GetAllWithoutPaginationResponse, error) {
	settings, err := s.service.GetAllWithoutPagination(ctx, req.GetReferrerId())
	if err != nil {
		log.Error(ctx, "DefaultSettings.GetAllWithoutPagination: %v", err)
		return nil, err
	}
	out := &pb.GetAllWithoutPaginationResponse{}

	for _, setting := range settings {
		out.List = append(out.List, toProtoDefaultSettings(setting))
	}
	return out, nil
}

func toProtoDefaultSettings(s *model.DefaultSettings) *pb.Settings {
	if s == nil {
		return nil
	}
	return &pb.Settings{
		Product:     s.Product,
		Coin:        s.Coin,
		Fee:         s.Fee.String(),
		ReferralFee: s.ReferralFee.String(),
		CreatedAt:   timestamppb.New(s.CreatedAt),
	}
}
