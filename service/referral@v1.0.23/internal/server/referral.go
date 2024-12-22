// Package server implements a simple web server.
package server

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/service"
	pb "code.emcdtech.com/emcd/service/referral/protocol/referral"
)

// ReferralService is an interface that defines the methods for managing referrals.
// It provides functionality to create, update, delete, get, list, and retrieve the history of referrals.
type ReferralService interface {
	Create(ctx context.Context, in *model.Referral) error
	Update(ctx context.Context, in *model.Referral) error
	Delete(ctx context.Context, userID uuid.UUID, product string, coin string) error
	Get(ctx context.Context, userID uuid.UUID, product string, coin string) (*model.Referral, error)
	List(ctx context.Context, userID uuid.UUID, skip int32, take int32) ([]*model.Referral, int, error)
	History(ctx context.Context, userID uuid.UUID, product string, coin string) ([]*model.Referral, error)
	CreateMultiple(ctx context.Context, rs []*model.Referral) error
	GetUserReferrals(ctx context.Context, userID uuid.UUID, skip, limit int32) ([]*model.UserReferral, int64, error)
	UpdateFeeWithMultiplier(
		ctx context.Context,
		userID uuid.UUID,
		product string,
		coins []string,
		multiplier decimal.Decimal,
	) error
	UpdateFee(ctx context.Context, userID uuid.UUID, product string, fees map[string]decimal.Decimal) error
	UpdateWithPromoCode(ctx context.Context, cms *model.CoinsMultipliers) error

	// UpdateFeeByCoinAndProduct обновляет комиссии у пользователя по монетам и продукту.
	UpdateFeeByCoinAndProduct(ctx context.Context, userUUID string, fees []model.SettingForCoinAndProduct) error
	// UpdateReferralUUIDByUserUUID обновление реферальной связи у пользователя.
	UpdateReferralUUIDByUserUUID(ctx context.Context, userUUID string, referralUUID string) error
}

// ReferralServer is a type that represents a server for managing referrals.
// It implements the pb.UnimplementedReferralServiceServer interface.
// It also contains a service of type service.Referral, which provides the implementation for managing referrals.
type ReferralServer struct {
	pb.UnimplementedReferralServiceServer
	service                   ReferralService
	defaultSettings           service.DefaultSettings
	defaultWhiteLabelSettings service.DefaultWhitelabelSettings
}

// NewReferralServer creates a new instance of ReferralServer with the provided service.
func NewReferralServer(
	s ReferralService,
	defaultSettings service.DefaultSettings,
	defaultWlSettings service.DefaultWhitelabelSettings,
) *ReferralServer {
	return &ReferralServer{
		service:                   s,
		defaultSettings:           defaultSettings,
		defaultWhiteLabelSettings: defaultWlSettings,
	}
}

// Create creates a referral in the ReferralServer.
// It accepts a context and a pb.CreateRequest as input,
// which contains the referral data to be created.
// The referral data is extracted from the input request and stored in a model.Referral struct.
// The Create method of the r.service is called passing the extracted referral data.
// If any common occurs during the creation process, it is logged and returned.
// Finally, it returns a pb.CreateResponse and a nil common.
func (r *ReferralServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	req := new(model.Referral)

	err := req.FromProto(in.Referral)
	if err != nil {
		log.Error(ctx, "ReferralServer.Create.FromProto: %v", err)
		return nil, err
	}

	err = r.service.Create(ctx, req)
	if err != nil {
		log.Error(ctx, "ReferralServer.Create: %v", err)
		return nil, err
	}

	return &pb.CreateResponse{}, nil
}

// Update updates a referral in the ReferralServer.
// It accepts a context and a pb.UpdateRequest as input,
// which contains the updated referral data.
// The referral data is extracted from the input request and stored in a model.Referral struct.
// The Update method of the r.service is called passing the extracted referral data.
// If any common occurs during the update process, it is logged and returned.
// Finally, it returns a pb.UpdateResponse and a nil common.
func (r *ReferralServer) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	req := new(model.Referral)

	err := req.FromProto(in.Referral)
	if err != nil {
		log.Error(ctx, "ReferralServer.Update.FromProto: %v", err)
		return nil, err
	}

	err = r.service.Update(ctx, req)
	if err != nil {
		log.Error(ctx, "ReferralServer.Update: %v", err)
		return nil, err
	}

	return &pb.UpdateResponse{}, nil
}

// Delete removes a referral record based on the provided information.
// It calls the Delete method of the ReferralService to delete the record.
// If an common occurs during the deletion, it logs the common and returns it.
// Otherwise, it returns an empty DeleteResponse.
func (r *ReferralServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		log.Error(ctx, "ReferralServer.Delete.ParseUserID: %v", err)
		return nil, err
	}

	err = r.service.Delete(ctx, userID, in.Product, in.Coin)
	if err != nil {
		log.Error(ctx, "ReferralServer.Delete: %v", err)
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}

// Get retrieves a referral from the ReferralServer.
// It accepts a context and a pb.GetRequest as input,
// which contains the user ID, product, and coin of the referral to be retrieved.
// The Get method of the r.service is called passing the extracted parameters.
// If any common occurs during the retrieval process, it is logged and returned.
// Finally, it creates a pb.GetResponse and sets its fields with the retrieved referral data.
// It returns the pb.GetResponse and a nil common.
func (r *ReferralServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		log.Error(ctx, "ReferralServer.Get.ParseUserID: %v", err)
		return nil, err
	}

	resp, err := r.service.Get(ctx, userID, in.Product, in.Coin)
	if err != nil {
		log.Error(ctx, "ReferralServer.Get: %v", err)
		return nil, err
	}

	out := &pb.GetResponse{
		Referral: &pb.Referral{
			UserId:        resp.UserID.String(),
			Product:       resp.Product,
			Coin:          resp.Coin,
			WhitelabelId:  resp.WhitelabelID.String(),
			Fee:           resp.Fee.String(),
			WhitelabelFee: resp.WhitelabelFee.String(),
			ReferralFee:   resp.ReferralFee.String(),
			ReferralId:    resp.ReferralID.String(),
			CreatedAt:     timestamppb.New(resp.CreatedAt),
		},
	}

	return out, nil
}

// List is a method of the ReferralServer struct that handles the request to list referrals.
// It takes a context.Context object and a *pb.ListRequest object as parameters.
// It calls the service.List method to retrieve the list of referrals.
// If there is an common, it logs the common message and returns the common.
// It then prepares the ListResponse object and populates it with the referral information.
// Finally, it sets the TotalCount field of the response and returns the response object.
func (r *ReferralServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		log.Error(ctx, "ReferralServer.List.ParseUserID: %v", err)
		return nil, err
	}

	list, count, err := r.service.List(ctx, userID, in.Skip, in.Take)
	if err != nil {
		log.Error(ctx, "ReferralServer.List: %v", err)
		return nil, err
	}

	out := &pb.ListResponse{}

	for i := range list {
		o := &pb.Referral{
			UserId:        list[i].UserID.String(),
			Product:       list[i].Product,
			Coin:          list[i].Coin,
			WhitelabelId:  list[i].WhitelabelID.String(),
			Fee:           list[i].Fee.String(),
			WhitelabelFee: list[i].WhitelabelFee.String(),
			ReferralFee:   list[i].ReferralFee.String(),
			ReferralId:    list[i].ReferralID.String(),
			CreatedAt:     timestamppb.New(list[i].CreatedAt),
		}

		out.List = append(out.List, o)
	}

	out.TotalCount = int32(count)

	return out, nil
}

// History retrieves the referral history for a user based on the provided parameters.
// It takes a context.Context object and a *pb.HistoryRequest object as input parameters. It returns a *pb.HistoryResponse object
// and an common as the return values.
// If the referral server fails to retrieve the history or encounters an common, it logs the common and returns nil and the common.
// The function first calls the service.History method with the context, user ID, product, and coin obtained from the HistoryRequest object.
// It then constructs a new pb.HistoryResponse object and populates it with the referral history returned by the service.
// Finally, it returns the constructed pb.HistoryResponse object and nil as the common.
func (r *ReferralServer) History(ctx context.Context, in *pb.HistoryRequest) (*pb.HistoryResponse, error) {
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		log.Error(ctx, "ReferralServer.History.ParseUserID: %v", err)
		return nil, err
	}

	history, err := r.service.History(ctx, userID, in.Product, in.Coin)
	if err != nil {
		log.Error(ctx, "ReferralServer.History: %v", err)
		return nil, err
	}

	out := &pb.HistoryResponse{}

	for i := range history {
		o := &pb.Referral{
			UserId:        history[i].UserID.String(),
			Product:       history[i].Product,
			Coin:          history[i].Coin,
			WhitelabelId:  history[i].WhitelabelID.String(),
			Fee:           history[i].Fee.String(),
			WhitelabelFee: history[i].WhitelabelFee.String(),
			ReferralFee:   history[i].ReferralFee.String(),
			ReferralId:    history[i].ReferralID.String(),
			CreatedAt:     timestamppb.New(history[i].CreatedAt),
		}

		out.History = append(out.History, o)
	}

	return out, nil
}

func (r *ReferralServer) GetUserReferrals(
	ctx context.Context,
	req *pb.GetUserReferralsRequest,
) (*pb.GetUserReferralsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "ReferralServer.GetUserReferrals.ParseUserID: %v", err)

		return nil, err
	}

	res, total, err := r.service.GetUserReferrals(ctx, userID, req.GetSkip(), req.GetLimit())
	if err != nil {
		log.Error(ctx, "ReferralServer.GetUserReferrals.service: %v", err)
		return nil, err
	}

	return r.convertGetUserReferralsToProto(res, total), nil
}

func (r *ReferralServer) convertGetUserReferralsToProto(referrals []*model.UserReferral, total int64) *pb.GetUserReferralsResponse {
	list := make([]*pb.UserReferral, 0, len(referrals))
	for _, ref := range referrals {
		list = append(list, r.convertUserReferralsToProto(ref))
	}

	return &pb.GetUserReferralsResponse{
		List:  list,
		Total: total,
	}
}

func (r *ReferralServer) convertUserReferralsToProto(referral *model.UserReferral) *pb.UserReferral {
	if referral == nil {
		return &pb.UserReferral{}
	}

	return &pb.UserReferral{UserId: referral.UserID.String()}
}

func (r *ReferralServer) CreateMultiple(ctx context.Context, req *pb.CreateMultipleRequest) (*pb.CreateMultipleResponse, error) {
	rs := make([]*model.Referral, len(req.Referrals))
	for i := range rs {
		m := model.Referral{}

		err := m.FromProto(req.Referrals[i])
		if err != nil {
			log.Error(ctx, "ReferralServer.CreateMultiple: %v", err)
			return nil, err
		}

		rs[i] = &m
	}
	err := r.service.CreateMultiple(ctx, rs)
	if err != nil {
		log.Error(ctx, "ReferralServer.CreateMultiple: %v", err)
		return nil, err
	}
	return &pb.CreateMultipleResponse{}, nil
}

func (r *ReferralServer) UpdateFeeWithMultiplier(
	ctx context.Context,
	req *pb.UpdateFeeWithMultiplierRequest,
) (*pb.UpdateFeeWithMultiplierResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "reward.UpdateWithMultiplier: parse user id: %s. %v", req.UserId, err)
		return nil, err
	}
	err = r.service.UpdateFeeWithMultiplier(ctx, userID, req.Product, req.Coins, decimal.NewFromFloat32(req.Multiplier))
	if err != nil {
		log.Error(ctx, "reward.UpdateWithMultiplier: %v", err)
		return nil, err
	}
	return &pb.UpdateFeeWithMultiplierResponse{}, nil
}

func (r *ReferralServer) UpdateFeeToDefault(
	ctx context.Context,
	req *pb.UpdateFeeToDefaultRequest,
) (*pb.UpdateFeeToDefaultResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "referral.UpdateFeeToDefault: parse user id: %s. %v", req.UserId, err)
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	fees := make(map[string]decimal.Decimal)
	for _, coin := range req.Coins {
		ref, err := r.service.Get(ctx, userID, req.Product, coin) //nolint:govet
		if err != nil {
			log.Error(
				ctx,
				"referral.UpdateFeeToDefault: user id: %s. product: %s. coin: %s. %v",
				req.UserId,
				req.Product,
				coin,
				err,
			)
			return nil, err
		}
		var fee decimal.Decimal
		if ref.WhitelabelID == uuid.Nil {
			ds, err := r.defaultSettings.Get(ctx, req.Product, coin)
			if err != nil {
				log.Error(ctx, "referral.UpdateFeeToDefault: product: %s. coin: %s. %v", req.Product, coin, err)
				continue
			}
			fee = ds.Fee
		} else {
			wds, err := r.defaultWhiteLabelSettings.GetV2ByCoin(ctx, req.Product, coin, ref.WhitelabelID)
			if err != nil {
				log.Error(ctx, "referral.UpdateFeeToDefault: product: %s. coin: %s. %v", req.Product, coin, err)
				continue
			}
			fee = wds.Fee
		}
		fees[coin] = fee
	}

	err = r.service.UpdateFee(ctx, userID, req.Product, fees)
	if err != nil {
		log.Error(ctx, "referral.UpdateFee: user id: %s. product: %s. %v", userID.String(), req.Product, err)
		return nil, err
	}
	return &pb.UpdateFeeToDefaultResponse{}, nil
}

func (r *ReferralServer) UpdateWithPromoCode(
	ctx context.Context,
	req *pb.UpdateWithPromoCodeRequest,
) (*pb.UpdateWithPromoCodeResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "referral.UpdateWithPromoCode: parse user id: %s. %v", req.UserId, err)
		return nil, err
	}
	actionID, err := uuid.Parse(req.ActionId)
	if err != nil {
		log.Error(ctx, "referral.UpdateWithPromoCode: parse action id: %s. %v", req.ActionId, err)
		return nil, err
	}
	feeMultipliers, err := toDecimalMap(req.FeeMultipliers)
	if err != nil {
		log.Error(ctx, "referral.UpdateWithPromoCode: toDecimalMap: feeMultipliers: %v", err)
		return nil, err
	}
	refFeeMultipliers, err := toDecimalMap(req.RefFeeMultipliers)
	if err != nil {
		log.Error(ctx, "referral.UpdateWithPromoCode: toDecimalMap: refFeeMultipliers: %v", err)
		return nil, err
	}
	err = r.service.UpdateWithPromoCode(ctx, &model.CoinsMultipliers{
		UserID:            userID,
		Product:           req.Product,
		FeeMultipliers:    feeMultipliers,
		RefFeeMultipliers: refFeeMultipliers,
		ActionID:          actionID,
		CreatedAt:         time.Now().UTC(),
	})
	if err != nil {
		log.Error(ctx, "referral.UpdateWithPromoCode: %v", err)
		return nil, err
	}
	return &pb.UpdateWithPromoCodeResponse{}, nil
}

func (r *ReferralServer) SetFee(ctx context.Context, req *pb.SetFeeRequest) (*pb.SetFeeResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		log.Error(ctx, "referral.SetFee: parse user uuid: %s. %v", req.GetUserUuid(), err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if len(req.GetFees()) == 0 {
		log.Error(ctx, "referral.SetFee: no fees to update")
		return nil, status.Errorf(codes.InvalidArgument, "referral.SetFee: no fees to update")
	}

	newFees := make([]model.SettingForCoinAndProduct, 0, len(req.GetFees()))
	for _, fee := range req.GetFees() {
		newFee, err := decimal.NewFromString(fee.GetFee())
		if err != nil {
			log.Error(ctx, "referral.SetFee: parse fee: %s. %v", fee.GetFee(), err)
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if fee.GetProduct() == "" {
			log.Error(ctx, "referral.SetFee: product is empty")
			return nil, status.Errorf(codes.InvalidArgument, "product is empty")
		}
		if fee.GetCoin() == "" {
			log.Error(ctx, "referral.SetFee: coin is empty")
			return nil, status.Errorf(codes.InvalidArgument, "coin is empty")
		}

		newFees = append(newFees, model.SettingForCoinAndProduct{
			Coin:    fee.GetCoin(),
			Product: fee.GetProduct(),
			Fee:     newFee,
		})
	}

	err = r.service.UpdateFeeByCoinAndProduct(context.WithoutCancel(ctx), userUUID.String(), newFees)
	if err != nil {
		log.Error(ctx, "SetFee.UpdateFeeByCoinAndProduct: %s, user_uuid: %s", err.Error(), userUUID.String())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.SetFeeResponse{}, nil
}

func toDecimalMap(m map[string]string) (map[string]decimal.Decimal, error) {
	res := make(map[string]decimal.Decimal, len(m))
	for k, v := range m {
		dec, err := decimal.NewFromString(v)
		if err != nil {
			return nil, fmt.Errorf("parse decimal: key: %s. val: %s. %w", k, v, err)
		}
		res[k] = dec
	}
	return res, nil
}

func (r *ReferralServer) SetReferralUUID(ctx context.Context, req *pb.SetReferralUUIDRequest) (*pb.SetReferralUUIDResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	referralUUID, err := uuid.Parse(req.GetReferralUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = r.service.UpdateReferralUUIDByUserUUID(context.WithoutCancel(ctx), userUUID.String(), referralUUID.String())
	if err != nil {
		log.Error(ctx, "SetReferralUUID.UpdateReferralUUIDByUserUUID: %s, user_uuid: %s, referral_uuid: %s",
			err.Error(), userUUID.String(), referralUUID.String())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.SetReferralUUIDResponse{}, nil
}
