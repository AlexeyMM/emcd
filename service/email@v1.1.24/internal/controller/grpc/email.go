package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	businessErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/email/internal/berrors"
	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/service"
	pb "code.emcdtech.com/emcd/service/email/protocol/email"
)

var (
	errTemplateNotFound        = businessErr.NewError("email-0001", berrors.ErrTemplateNotFound.Error())
	errNoMailSenderAvailable   = businessErr.NewError("email-0002", berrors.ErrNoMailSenderAvailable.Error())
	errProviderSettingNotFound = businessErr.NewError("email-0003", berrors.ErrProviderSettingNotFound.Error())
	errInvalidArgument         = businessErr.NewError("email-0004", "invalid argument")
)

type EmailServiceServer struct {
	pb.UnsafeEmailServiceServer
	emailService service.Email
}

func NewEmail(emailService service.Email) *EmailServiceServer {
	return &EmailServiceServer{
		emailService: emailService,
	}
}

func (e *EmailServiceServer) SendPasswordRestoration(
	ctx context.Context,
	req *pb.SendPasswordRestorationRequest,
) (*pb.SendPasswordRestorationResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	if err := e.emailService.SendPasswordRestoration(ctx, userId, req.Token, req.GetDomain()); err != nil {
		log.Error(ctx, "SendPasswordRestoration: user: %s: %w", req.UserId, err)
		return nil, convertError(err)
	}
	return &pb.SendPasswordRestorationResponse{}, nil
}

func (e *EmailServiceServer) SendWorkerChangedState(
	ctx context.Context,
	req *pb.SendWorkerChangedStateRequest,
) (*pb.SendWorkerChangedStateResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "send worker changed state: parse wl id: %w. wl id: %v", err, req.WhiteLabelID)
		return nil, fmt.Errorf("send worker changed state: parse wl id: %w", err)
	}
	err = e.emailService.SendWorkerChangedState(ctx, &model.Worker{
		Name:           req.WorkerName,
		Username:       req.Username,
		Coin:           req.Coin,
		IsOn:           req.IsOn,
		StateChangedAt: req.StateChangedAt.AsTime(),
	}, req.Email, req.GetDomain(), wlID, req.Language)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, convertError(err)
	}
	return &pb.SendWorkerChangedStateResponse{}, nil
}

func (e *EmailServiceServer) SendWalletChangedAddress(
	ctx context.Context,
	req *pb.SendWalletChangedAddressRequest,
) (*pb.SendWalletChangedAddressResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		err = fmt.Errorf("send wallet changed state: parsing user_id=%s: %w", req.GetUserId(), err)
		log.Error(ctx, err.Error())
		return nil, err
	}

	err = e.emailService.SendWalletChangedAddress(ctx, userID, req.Domain, req.Token, req.CoinCode)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, convertError(err)
	}

	return &pb.SendWalletChangedAddressResponse{}, nil
}

func (e *EmailServiceServer) SendRegister(
	ctx context.Context,
	req *pb.SendRegisterRequest,
) (*pb.SendRegisterResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "send register: parsing wl id: %v. %v", err, req.WhiteLabelID)
		return nil, fmt.Errorf("send register: parsing wl id: %w", err)
	}
	log.Info(ctx, "send register request: %v", req)
	err = e.emailService.SendRegister(ctx, wlID, req.Email, req.GetDomain(), req.Token, req.Language)
	if err != nil {
		log.Error(ctx, "send register: %v", err)
		return nil, convertError(err)
	}
	return &pb.SendRegisterResponse{}, nil
}

func (e *EmailServiceServer) SendMobileTwoFaOff(
	ctx context.Context,
	req *pb.SendMobileTwoFaOffRequest,
) (*pb.SendMobileTwoFaOffResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	err = e.emailService.SendMobileTwoFaOff(ctx, userId, req.GetDomain())
	if err != nil {
		log.Error(ctx, "sendMobileTwoFaOff: %s", err.Error())
		return nil, convertError(err)
	}
	return &pb.SendMobileTwoFaOffResponse{}, nil
}

func (e *EmailServiceServer) SendMobileTwoFaOn(
	ctx context.Context,
	req *pb.SendMobileTwoFaOnRequest,
) (*pb.SendMobileTwoFaOnResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	err = e.emailService.SendMobileTwoFaOn(ctx, userId, req.Token, req.GetDomain())
	if err != nil {
		log.Error(ctx, "sendMobileTwoFaOn: %s", err.Error())
		return nil, convertError(err)
	}
	return &pb.SendMobileTwoFaOnResponse{}, nil
}

func (e *EmailServiceServer) SendGoogleTwoFaOff(
	ctx context.Context,
	req *pb.SendGoogleTwoFaOffRequest,
) (*pb.SendGoogleTwoFaOffResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	err = e.emailService.SendGoogleTwoFaOff(ctx, userId, req.GetDomain())
	if err != nil {
		log.Error(ctx, "sendGoogleTwoFaOff: %s", err.Error())
		return nil, convertError(err)
	}

	return &pb.SendGoogleTwoFaOffResponse{}, nil
}

func (e *EmailServiceServer) SendGoogleTwoFaOn(
	ctx context.Context,
	req *pb.SendGoogleTwoFaOnRequest,
) (*pb.SendGoogleTwoFaOnResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	err = e.emailService.SendGoogleTwoFaOn(ctx, userId, req.Token, req.GetDomain())
	if err != nil {
		log.Error(ctx, "sendGoogleTwoFaOn: %s", err.Error())
		return nil, convertError(err)
	}

	return &pb.SendGoogleTwoFaOnResponse{}, nil
}

func (e *EmailServiceServer) SendUserHashrateDecreased(
	ctx context.Context,
	req *pb.SendUserHashrateDecreasedRequest,
) (*pb.SendUserHashrateDecreasedResponse, error) {
	decreasedBy, err := decimal.NewFromString(req.DecreasedBy)
	if err != nil {
		log.Error(ctx, "sendUserHashrateDecreased: parse decreased by: %s: email: %s: %v", req.DecreasedBy, req.Email, err)
		return nil, err
	}
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "sendUserHashrateDecreased: parsing wl id: %v. email: %s: %v", err, req.Email, req.WhiteLabelID)
		return nil, err
	}
	err = e.emailService.SendUserHashrateDecreased(
		ctx,
		req.Email,
		req.GetDomain(),
		wlID,
		req.Language,
		decreasedBy,
		req.Coin,
	)
	if err != nil {
		log.Error(ctx, "sendUserHashrateDecreased: email: %s: %v", req.Email, err)
		return nil, convertError(err)
	}
	return &pb.SendUserHashrateDecreasedResponse{}, nil
}

func (e *EmailServiceServer) ListMessages(
	ctx context.Context,
	in *pb.GetSentEmailMessagesByEmailRequest,
) (*pb.GetSentEmailMessagesByEmailResponse, error) {
	var err error

	list, count, err := e.emailService.ListMessages(ctx, in.Email, in.Type, in.Skip, in.Take)
	if err != nil {
		log.Error(ctx, "EmailServer.ListMessages: %v", err)
		return nil, err
	}

	out := &pb.GetSentEmailMessagesByEmailResponse{}

	for i := range list {
		o := &pb.Email{
			Id:        list[i].ID.String(),
			Email:     list[i].Email,
			Type:      list[i].Type.String(),
			CreatedAt: timestamppb.New(list[i].CreatedAt),
		}

		out.List = append(out.List, o)
	}

	out.TotalCount = int32(count)

	return out, nil
}

func (e *EmailServiceServer) SendPasswordChange(
	ctx context.Context,
	req *pb.SendPasswordChangeRequest,
) (*pb.SendPasswordChangeResponse, error) {
	const op = "grpc.email.SendPasswordChange"
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	err = e.emailService.SendPasswordChange(ctx, userId, req.Token, req.GetDomain())
	if err != nil {
		log.Error(ctx, "%s: user: %s: %v", op, req.UserId, err)
		return nil, convertError(err)
	}
	return &pb.SendPasswordChangeResponse{}, nil
}

func (e *EmailServiceServer) SendPhoneDelete(
	ctx context.Context,
	req *pb.SendPhoneDeleteRequest,
) (*pb.SendPhoneDeleteResponse, error) {
	const op = "grpc.email.SendPhoneDelete"
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Warn(ctx, "parse user id: %s", err.Error())
		return nil, errInvalidArgument
	}
	err = e.emailService.SendPhoneDelete(ctx, userId, req.Token, req.GetDomain())
	if err != nil {
		log.Error(ctx, "%s: user: %s: %v", op, req.UserId, err)
		return nil, convertError(err)
	}
	return &pb.SendPhoneDeleteResponse{}, nil
}

func (e *EmailServiceServer) SendReferralRewardPayouts(
	ctx context.Context,
	req *pb.SendReferralRewardPayoutsRequest,
) (*pb.SendReferralRewardPayoutsResponse, error) {
	attachments := make([]model.Attachment, 0, len(req.Attachments))
	for _, a := range req.Attachments {
		attachment := model.Attachment{
			Filename: a.Name,
			Data:     a.Body,
		}
		attachments = append(attachments, attachment)
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, errInvalidArgument
	}
	err = e.emailService.SendReferralRewardPayouts(
		ctx,
		userID,
		req.GetDomain(),
		req.GetFrom().AsTime(),
		req.GetTo().AsTime(),
		attachments,
	)
	if err != nil {
		log.Error(ctx, "send referral reward payouts (%s): %w", userID, err)
		return nil, convertError(err)
	}
	return &pb.SendReferralRewardPayoutsResponse{}, nil
}

func (e *EmailServiceServer) SendChangeEmail(
	ctx context.Context,
	req *pb.SendChangeEmailRequest,
) (*pb.SendChangeEmailResponse, error) {
	wlID, err := uuid.Parse(req.WhiteLabelID)
	if err != nil {
		log.Error(ctx, "send change email: parsing wl id: %v. %v", err, req.WhiteLabelID)
		return nil, fmt.Errorf("send change email: parsing wl id: %w", err)
	}
	log.Info(ctx, "send change email request: %v", req)
	if req.Language == "" {
		req.Language = "en"
	}
	err = e.emailService.SendChangeEmail(ctx, wlID, req.Email, req.GetDomain(), req.Token, req.Language)
	if err != nil {
		log.Error(ctx, "send change email: %v", err)
		return nil, convertError(err)
	}
	return &pb.SendChangeEmailResponse{}, nil
}

func (e *EmailServiceServer) SendStatisticsReport(
	ctx context.Context,
	req *pb.SendStatisticsReportRequest,
) (*pb.SendStatisticsReportResponse, error) {
	const op = "grpc.email.SendStatisticsReport"
	reportType, err := convertReportTypeToModel(req.ReportType)
	if err != nil {
		log.Error(ctx, "%s: email: %s: %v", op, req.Email, err)
		return nil, convertError(err)
	}
	if err := e.emailService.SendStatisticsReport(ctx, req.Email, req.ReportLink, req.Language, reportType); err != nil {
		log.Error(ctx, "%s: email: %s: %v", op, req.Email, err)
		return nil, convertError(err)
	}
	return &pb.SendStatisticsReportResponse{}, nil
}

func (e *EmailServiceServer) SendSwapSupportMessage(
	ctx context.Context,
	req *pb.SendSwapSupportMessageRequest,
) (*pb.SendSwapSupportMessageResponse, error) {
	err := e.emailService.SendSwapSupportMessage(ctx, req.Name, req.UserEmail, req.Text)
	if err != nil {
		log.Error(ctx, "sendSwapSupportMessage: %v", err)
		return nil, convertError(err)
	}
	return &pb.SendSwapSupportMessageResponse{}, nil
}

func (e *EmailServiceServer) SendInitialSwapMessage(
	ctx context.Context,
	req *pb.SendInitialSwapMessageRequest,
) (
	*pb.SendInitialSwapMessageResponse, error,
) {
	err := e.emailService.SendInitialSwapMessage(ctx, req.Email, req.Language, req.Link)
	if err != nil {
		log.Error(ctx, "sendSwapMessage: %v", err)
		return nil, convertError(err)
	}
	return &pb.SendInitialSwapMessageResponse{}, nil
}

func (e *EmailServiceServer) SendSuccessfulSwapMessage(
	ctx context.Context,
	req *pb.SendSuccessfulSwapMessageRequest,
) (*pb.SendSuccessfulSwapMessageResponse, error) {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		log.Error(ctx, "sendSuccessfulSwapMessage: parse swap id: %v", err)
		return nil, convertError(err)
	}

	err = e.emailService.SendSuccessfulSwapMessage(
		ctx,
		swapID,
		req.From,
		req.To,
		req.Address,
		req.Email,
		req.Language,
		req.ExecutionTime,
	)
	if err != nil {
		log.Error(ctx, "sendSuccessfulSwapMessage: %v", err)
		return nil, convertError(err)
	}
	return &pb.SendSuccessfulSwapMessageResponse{}, nil
}
