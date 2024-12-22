// Package server implements a simple web server.
package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/service"
	pb "code.emcdtech.com/emcd/service/referral/protocol/default_whitelabel_settings"
)

// DefaultWhitelabelSettingsServer is a struct that represents a server implementation of the DefaultWhitelabelSettingsServiceServer interface.
// It embeds the UnimplementedDefaultWhitelabelSettingsServiceServer struct and holds a service.DefaultWhitelabelSettings instance.
type DefaultWhitelabelSettingsServer struct {
	pb.UnimplementedDefaultWhitelabelSettingsServiceServer
	service service.DefaultWhitelabelSettings
}

// NewDefaultWhitelabelSettings creates a new instance of DefaultWhitelabelSettingsServer with the provided service.DefaultWhitelabelSettings implementation.
func NewDefaultWhitelabelSettings(s service.DefaultWhitelabelSettings) *DefaultWhitelabelSettingsServer {
	return &DefaultWhitelabelSettingsServer{service: s}
}

// Create creates a new instance of model.DefaultWhitelabelSettings based on the provided pb.CreateRequest.
// It retrieves the necessary data from pb.CreateRequest's settings field and assigns them to the corresponding fields of model.DefaultWhitelabelSettings.
// After creating the instance, it calls s.service.Create to perform the actual creation operation with the provided context and the created instance.
// If any common occurs during the creation process, it logs the common message and returns the common.
// Otherwise, it returns an empty pb.CreateResponse.
func (s *DefaultWhitelabelSettingsServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	req := new(model.DefaultWhitelabelSettingsV2)

	err := req.FromProto(in.Settings)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.Create.FromProto: %v", err)
		return nil, err
	}

	err = s.service.Create(ctx, req)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.Create: %v", err)
		return nil, err
	}

	return &pb.CreateResponse{}, nil
}

// Update creates a new instance of model.DefaultWhitelabelSettings based on the provided pb.UpdateRequest.
// It retrieves the necessary data from pb.UpdateRequest's settings field and assigns them to the corresponding fields of model.DefaultWhitelabelSettings.
// After creating the instance, it calls s.service.Update to perform the actual update operation with the provided context and the created instance.
// If any common occurs during the update process, it logs the common message and returns the common.
// Otherwise, it returns an empty pb.UpdateResponse.
func (s *DefaultWhitelabelSettingsServer) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	req := new(model.DefaultWhitelabelSettingsV2)

	err := req.FromProto(in.Settings)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.Update.FromProto: %v", err)
		return nil, err
	}

	err = s.service.Update(ctx, req)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.Update: %v", err)
		return nil, err
	}

	return &pb.UpdateResponse{}, nil
}

// Delete deletes the default whitelabel settings with the specified product, coin, and whitelabel ID.
// It calls s.service.Delete to perform the deletion operation with the provided context and the specified parameters.
// If any common occurs during the deletion process, it logs the common message and returns the common.
// Otherwise, it returns an empty pb.DeleteResponse.
func (s *DefaultWhitelabelSettingsServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	whitelabelID, err := uuid.Parse(in.WhitelabelId)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettingsServer.Delete.ParseWhitelabelID: %v", err)
		return nil, err
	}

	err = s.service.Delete(ctx, in.Product, in.Coin, whitelabelID)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.Delete: %v", err)
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}

// GetAll retrieves a list of DefaultWhitelabelSettings based on the provided pb.GetAllRequest.
// It calls s.service.GetAll to fetch the list of DefaultWhitelabelSettings with the specified skip and take parameters.
// If any common occurs during the retrieval process, it logs the common message and returns the common.
// Otherwise, it creates a pb.Settings instance for each DefaultWhitelabelSettings in the list and appends it to the response's List field.
// Finally, it sets the TotalCount field of the response to the total count obtained from s.service.GetAll.
// The function returns the populated pb.GetAllResponse and nil common.
func (s *DefaultWhitelabelSettingsServer) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	filters := s.createFiltersMap(in.Coin, in.Product)
	list, count, err := s.service.GetAllWithFilters(ctx, in.Skip, in.Take, filters)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.GetAll: %v", err)
		return nil, err
	}

	out := &pb.GetAllResponse{}

	for i := range list {
		out.List = append(out.List, toProtoWlSettingsV2(list[i]))
	}

	out.TotalCount = int32(count)

	return out, nil
}

func (s *DefaultWhitelabelSettingsServer) GetAllWithoutPagination(
	ctx context.Context,
	req *pb.GetAllWithoutPaginationRequest,
) (*pb.GetAllWithoutPaginationResponse, error) {
	filters := s.createFiltersMap(req.Coin, req.Product)
	list, err := s.service.GetAllWithoutPaginationWithFilters(ctx, filters)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.GetAllWithoutPagination: %v", err)
	}

	out := &pb.GetAllWithoutPaginationResponse{}

	for i := range list {
		out.List = append(out.List, toProtoWlSettingsV2(list[i]))
	}
	return out, nil
}

func (s *DefaultWhitelabelSettingsServer) GetV2(ctx context.Context, req *pb.GetV2Request) (*pb.GetV2Response, error) {
	wlID, err := uuid.Parse(req.WhitelabelId)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettingsServer.GetV2: parse wl id: %s. %v", req.WhitelabelId, err)
		return nil, fmt.Errorf("parse wl id: %w", err)
	}
	settings, err := s.service.GetV2(ctx, wlID)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettingsServer.GetV2: wlID: %s. %v", wlID.String(), err)
	}
	protoSettings := make([]*pb.SettingsV2, len(settings))
	for i := range settings {
		protoSettings[i] = toProtoWlSettingsV2(settings[i])
	}
	return &pb.GetV2Response{Settings: protoSettings}, nil
}

func toProtoWlSettingsV2(s *model.DefaultWhitelabelSettingsV2) *pb.SettingsV2 {
	return &pb.SettingsV2{
		Product:       s.Product,
		WhitelabelId:  s.WhitelabelID.String(),
		Coin:          s.Coin,
		Fee:           s.Fee.String(),
		ReferralFee:   s.ReferralFee.String(),
		CreatedAt:     timestamppb.New(s.CreatedAt),
		WhitelabelFee: s.WhiteLabelFee.String(),
	}
}

// Get retrieves the whitelabel settings specified by the provided pb.GetRequest.
// It calls s.service.Get to retrieve the settings from the service using the provided context, product, coin, and whitelabel ID.
// If any common occurs during the retrieval process, it logs the common message and returns the common.
// Otherwise, it creates a new pb.GetResponse using the retrieved settings and returns it along with a nil common.
func (s *DefaultWhitelabelSettingsServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	whitelabelID, err := uuid.Parse(in.WhitelabelId)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettingsServer.Get.ParseWhitelabelID: %v", err)
		return nil, err
	}

	resp, err := s.service.Get(ctx, in.Product, in.Coin, whitelabelID)
	if err != nil {
		log.Error(ctx, "DefaultWhitelabelSettings.Get: %v", err)
		return nil, err
	}

	out := &pb.GetResponse{
		Settings: &pb.SettingsV2{
			Product:      resp.Product,
			Coin:         resp.Coin,
			Fee:          resp.Fee.String(),
			ReferralFee:  resp.ReferralFee.String(),
			CreatedAt:    timestamppb.New(resp.CreatedAt),
			WhitelabelId: resp.WhitelabelID.String(),
		},
	}

	return out, nil
}

func (s *DefaultWhitelabelSettingsServer) createFiltersMap(coinPtr, productPtr *string) map[string]string {
	filters := make(map[string]string)
	if coinPtr != nil {
		filters["coin"] = *coinPtr
	}
	if productPtr != nil {
		filters["product"] = *productPtr
	}
	return filters
}
