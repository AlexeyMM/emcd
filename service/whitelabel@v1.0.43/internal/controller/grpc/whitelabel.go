package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"code.emcdtech.com/emcd/service/whitelabel/internal/model"
	"code.emcdtech.com/emcd/service/whitelabel/internal/service"
	pb "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
)

// TODO: move to (templates, etc...)
const (
	defaultFirmwareInstructionEn = `
		Скачайте <a href="https://emcd.io/fw/" target='_blank' class='link'>прошивку</a> для увеличения хешрейта.
	`
	defaultFirmwareInstructionRu = `
		Please download <a href="https://emcd.io/fw/" target='_blank' class='link'>patch</a> for increase hash_rate.
	`
)

type WhiteLabel struct {
	whiteLabelService service.WhiteLabel
	pb.UnimplementedWhitelabelServiceServer

	langs []string
}

func NewWhiteLabel(whiteLabelService service.WhiteLabel, langs []string) *WhiteLabel {
	return &WhiteLabel{
		whiteLabelService: whiteLabelService,
		langs:             langs,
	}
}

func (w *WhiteLabel) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	wl, err := w.parseProto(req.WhiteLabel)
	if err != nil {
		log.Error().Msgf("create: %v", err)
		return nil, fmt.Errorf("whitelabel create: %w", err)
	}
	id, err := w.whiteLabelService.Create(ctx, wl)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.CreateResponse{
		Id: id.String(),
	}, nil
}

func (w *WhiteLabel) parseProto(wl *pb.WhiteLabel) (*model.WhiteLabel, error) {
	var (
		id  uuid.UUID
		err error
	)
	if wl.Id != "" {
		id, err = uuid.Parse(wl.Id)
		if err != nil {
			return nil, fmt.Errorf("parse proto: parse wl id %s: %w", wl.Id, err)
		}
	}
	return &model.WhiteLabel{
		ID:                    id,
		Domain:                wl.Domain,
		SegmentID:             wl.SegmentId,
		Origin:                wl.Origin,
		Prefix:                wl.Prefix,
		SenderEmail:           wl.SenderEmail,
		APIKey:                wl.ApiKey,
		URL:                   wl.Url,
		Version:               int(wl.Version),
		MasterSlave:           wl.MasterSlave,
		MasterFee:             wl.MasterFee,
		IsTwoFAEnabled:        wl.IsTwoFaEnabled,
		IsCaptchaEnabled:      wl.IsCaptchaEnabled,
		IsEmailConfirmEnabled: wl.IsEmailConfirmEnabled,
	}, nil
}

func (w *WhiteLabel) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	wls, count, err := w.whiteLabelService.GetAll(ctx, int(req.Skip), int(req.Take), req.Sort.Field, req.Sort.Asc)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	protoWls := make([]*pb.WhiteLabel, len(wls))
	for i := range wls {
		protoWls[i] = w.convertToProto(wls[i])
	}
	return &pb.GetAllResponse{
		WhiteLabels: protoWls,
		TotalCount:  int32(count),
	}, nil
}

func (w *WhiteLabel) convertToProto(wl *model.WhiteLabel) *pb.WhiteLabel {
	if wl == nil {
		return nil
	}
	return &pb.WhiteLabel{
		Id:                    wl.ID.String(),
		UserId:                wl.UserID,
		SegmentId:             wl.SegmentID,
		Origin:                wl.Origin,
		Prefix:                wl.Prefix,
		SenderEmail:           wl.SenderEmail,
		Domain:                wl.Domain,
		ApiKey:                wl.APIKey,
		Url:                   wl.URL,
		Version:               int32(wl.Version),
		MasterSlave:           wl.MasterSlave,
		MasterFee:             wl.MasterFee,
		IsTwoFaEnabled:        wl.IsTwoFAEnabled,
		IsCaptchaEnabled:      wl.IsCaptchaEnabled,
		IsEmailConfirmEnabled: wl.IsEmailConfirmEnabled,
	}
}

func (w *WhiteLabel) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	wl, err := w.parseProto(req.WhiteLabel)
	if err != nil {
		log.Error().Msgf("update: %v", err)
		return nil, fmt.Errorf("whitelabel update: %w", err)
	}
	err = w.whiteLabelService.Update(ctx, wl)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.UpdateResponse{}, nil
}

func (w *WhiteLabel) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Error().Msgf("delete: parse id %s: %v", req.Id, err)
		return nil, fmt.Errorf("whitelabel: delete: parse id %s: %w", req.Id, err)
	}
	err = w.whiteLabelService.Delete(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.DeleteResponse{}, nil
}

func (w *WhiteLabel) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Error().Msgf("get by id: parse id %s: %v", req.Id, err)
		return nil, fmt.Errorf("whitelabel: get by id: parse id %s: %w", req.Id, err)
	}
	wl, err := w.whiteLabelService.GetByID(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.GetByIDResponse{
		WhiteLabel: w.convertToProto(wl),
	}, nil
}

func (w *WhiteLabel) GetBySegmentID(ctx context.Context, req *pb.GetBySegmentIDRequest) (*pb.GetWLResponse, error) {
	wl, err := w.whiteLabelService.GetBySegmentID(ctx, int(req.GetSegmentId()))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return &pb.GetWLResponse{
		WhiteLabel: w.convertToProto(wl),
	}, nil
}

func (w *WhiteLabel) GetByUserID(ctx context.Context, req *pb.GetByUserIDRequest) (*pb.GetWLResponse, error) {
	wl, err := w.whiteLabelService.GetByUserID(ctx, int(req.GetUserId()))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return &pb.GetWLResponse{
		WhiteLabel: w.convertToProto(wl),
	}, nil
}

func (w *WhiteLabel) GetByOrigin(ctx context.Context, req *pb.GetByOriginRequest) (*pb.GetWLResponse, error) {
	wl, err := w.whiteLabelService.GetByOrigin(ctx, req.GetOrigin())
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return &pb.GetWLResponse{
		WhiteLabel: w.convertToProto(wl),
	}, nil
}

func (w *WhiteLabel) CheckByUserID(ctx context.Context, req *pb.CheckByUserIDRequest) (*pb.CheckWLResponse, error) {
	success, err := w.whiteLabelService.CheckByUserID(ctx, int(req.GetUserId()))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return &pb.CheckWLResponse{
		Success: success,
	}, nil
}

func (w *WhiteLabel) CheckByUserIDAndOrigin(ctx context.Context, req *pb.CheckByUserIDAndOriginRequest) (*pb.CheckWLResponse, error) {
	success, err := w.whiteLabelService.CheckByUserIDAndOrigin(ctx, int(req.GetUserId()), req.GetOrigin())
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return &pb.CheckWLResponse{
		Success: success,
	}, nil
}

func (w *WhiteLabel) GetV2WLs(ctx context.Context, _ *pb.GetV2WLsRequest) (*pb.GetV2WLsResponse, error) {
	resp, err := w.whiteLabelService.GetV2WLs(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	response := pb.GetV2WLsResponse{
		List: make([]*pb.WhiteLabel, len(resp)),
	}

	for i := range resp {
		response.List[i] = w.convertToProto(resp[i])
	}
	return &response, nil
}

func (w *WhiteLabel) GetConfigByOrigin(ctx context.Context, req *pb.GetConfigByOriginRequest) (*pb.WLConfigResponse, error) {
	resp, err := w.whiteLabelService.GetConfigByOrigin(ctx, req.GetOrigin())
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	wl, err := w.whiteLabelService.GetByID(ctx, uuid.MustParse(resp.WhitelabelID))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	response := pb.WLConfigResponse{
		Config: &pb.Config{
			RefId:                 resp.RefID,
			Origin:                resp.Origin,
			MediaId:               resp.MediaID,
			Title:                 resp.Title,
			Logo:                  resp.Logo,
			Favicon:               resp.Favicon,
			Commission:            resp.Commission,
			Colors:                w.structColorsToMap(resp.Colors),
			FirmwareInstruction:   resp.FirmwareInstruction,
			PossibleLangs:         resp.PossibleLang,
			Lang:                  resp.Lang,
			WhitelabelId:          resp.WhitelabelID,
			IsEmailConfirmEnabled: wl.IsEmailConfirmEnabled,
			IsTwoFaEnabled:        wl.IsTwoFAEnabled,
			IsCaptchaEnabled:      wl.IsCaptchaEnabled,
			Prefix:                wl.Prefix,
		},
	}
	for _, s := range resp.StratumLists {
		response.Config.StratumList = append(response.Config.StratumList, &pb.Stratum{
			Coin:   s.Coin,
			Region: s.Region,
			Number: s.Number,
			Url:    s.Url,
		})
	}
	return &response, nil
}

func (w *WhiteLabel) SetConfigByRefID(ctx context.Context, req *pb.SetConfigByRefIDRequest) (*pb.SetConfigByRefIDResponse, error) {
	lang := strings.ToLower(req.Config.Lang)
	if ok := w.isLanguageValid(lang); !ok {
		return nil, status.Error(codes.InvalidArgument, "current language not support")
	}

	cfg := req.GetConfig()
	config := model.WlConfig{
		RefID:               cfg.GetRefId(),
		Origin:              cfg.GetOrigin(),
		MediaID:             cfg.GetMediaId(),
		Title:               cfg.GetTitle(),
		Logo:                cfg.GetLogo(),
		Favicon:             cfg.GetFavicon(),
		Commission:          cfg.GetCommission(),
		Colors:              w.mapColorsToStruct(cfg.GetColors()),
		FirmwareInstruction: cfg.FirmwareInstruction,
		Lang:                lang,
	}
	if config.FirmwareInstruction == "" {
		config.FirmwareInstruction = getDefaultFirmwareInstruction(lang)
	}
	if err := w.whiteLabelService.SetConfigByRefID(ctx, &config); err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	response := pb.SetConfigByRefIDResponse{
		Success: true,
	}

	return &response, nil
}

func (w *WhiteLabel) SetAllowOrigin(ctx context.Context, req *pb.AllowOrigin) (*pb.SuccessResponse, error) {
	opt := model.AllowOrigin{
		UserID: req.GetUserId(),
		Origin: req.GetOrigin(),
	}

	err := w.whiteLabelService.SetAllowOrigin(ctx, &opt)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	response := pb.SuccessResponse{
		Success: true,
	}

	return &response, nil
}

func (w *WhiteLabel) GetAllowOrigins(ctx context.Context, _ *pb.EmptyRequest) (*pb.GetAllowOriginsResponse, error) {
	resp, err := w.whiteLabelService.GetAllowOrigins(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	response := pb.GetAllowOriginsResponse{
		List: make([]*pb.AllowOrigin, len(resp)),
	}

	for i := range resp {
		response.List[i] = &pb.AllowOrigin{
			UserId: resp[i].UserID,
			Origin: resp[i].Origin,
		}
	}

	return &response, nil
}

func (w *WhiteLabel) SetStratum(ctx context.Context, req *pb.Stratum) (*pb.SuccessResponse, error) {
	opt := model.Stratum{
		RefID:  req.GetRefId(),
		Coin:   req.GetCoin(),
		Region: req.GetRegion(),
		Number: req.GetNumber(),
		Url:    req.GetUrl(),
	}

	err := w.whiteLabelService.SetStratum(ctx, &opt)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	response := pb.SuccessResponse{
		Success: true,
	}

	return &response, nil
}

func (w *WhiteLabel) mapColorsToStruct(mapColors map[string]string) model.Colors {
	var colors model.Colors

	res, err := json.Marshal(mapColors)
	if err != nil {
		log.Error().Msg(err.Error())
		return colors
	}

	if err = json.Unmarshal(res, &colors); err != nil {
		log.Error().Msg(err.Error())
	}

	return colors
}

func (w *WhiteLabel) structColorsToMap(colors model.Colors) map[string]string {
	mapColors := make(map[string]string)

	res, err := json.Marshal(colors)
	if err != nil {
		log.Error().Msg(err.Error())
		return mapColors
	}

	if err = json.Unmarshal(res, &mapColors); err != nil {
		log.Error().Msg(err.Error())
	}

	return mapColors
}

func (w *WhiteLabel) GetFullByUserID(ctx context.Context, req *pb.GetByUserIDRequest) (*pb.GetWLResponse, error) {
	wl, err := w.whiteLabelService.GetFullByUserID(ctx, int(req.UserId))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.GetWLResponse{
		WhiteLabel: w.convertToProto(wl),
	}, nil
}

func (w *WhiteLabel) GetCoins(ctx context.Context, in *pb.GetCoinsRequest) (*pb.GetCoinsResponse, error) {
	uuid, err := uuid.Parse(in.WlId)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	coins, err := w.whiteLabelService.GetCoins(ctx, uuid)
	var pbCoins []*pb.Coin
	for _, coin := range coins {
		pbCoins = append(pbCoins, w.toCoinProto(coin))
	}

	return &pb.GetCoinsResponse{
		Coins: pbCoins,
	}, err
}

func (w *WhiteLabel) toCoinProto(coin *model.WLCoins) *pb.Coin {
	return &pb.Coin{
		CoinId: coin.CoinID,
		WlId:   coin.WlID.String(),
	}
}

func (w *WhiteLabel) AddCoin(ctx context.Context, in *pb.AddCoinRequest) (*pb.AddCoinResponse, error) {
	wlID, err := uuid.Parse(in.Coin.WlId)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	err = w.whiteLabelService.AddCoin(ctx, wlID, in.Coin.CoinId)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.AddCoinResponse{}, nil
}

func (w *WhiteLabel) DeleteCoin(ctx context.Context, in *pb.DeleteCoinRequest) (*pb.DeleteCoinResponse, error) {
	wlID, err := uuid.Parse(in.Coin.WlId)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	err = w.whiteLabelService.DeleteCoin(ctx, wlID, in.Coin.CoinId)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &pb.DeleteCoinResponse{}, nil
}

func getDefaultFirmwareInstruction(lang string) string {
	if lang == "ru" {
		return defaultFirmwareInstructionRu
	}
	return defaultFirmwareInstructionEn
}

func (w *WhiteLabel) isLanguageValid(lang string) bool {
	for _, l := range w.langs {
		if lang == l {
			return true
		}
	}
	return false
}

func (w *WhiteLabel) GetStratumList(ctx context.Context, req *pb.GetStratumListRequest) (*pb.GetStratumListResponse, error) {
	wlID := req.GetWhitelabelId()

	id, err := uuid.Parse(wlID)
	if err != nil {
		log.Error().Msgf("get by id: parse id: %s, err: %v", wlID, err)
		return nil, fmt.Errorf("GetStratumList: get whitelabel by id: parse id :%s, err: %w", wlID, err)
	}
	wl, err := w.whiteLabelService.GetByID(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	log.Info().Msgf("segment_id: %d", wl.SegmentID)
	stratumList, err := w.whiteLabelService.GetWLStratumList(ctx, wl.SegmentID)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	resp := &pb.GetStratumListResponse{StratumList: make([]*pb.Stratum, 0, len(stratumList))}
	for _, s := range stratumList {
		resp.StratumList = append(resp.StratumList, &pb.Stratum{
			RefId:  s.RefID,
			Coin:   s.Coin,
			Region: s.Region,
			Number: s.Number,
			Url:    s.Url,
		})
	}
	return resp, nil
}
