package grpc

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
	pb "code.emcdtech.com/emcd/service/email/protocol/email"
)

const defaultTemplateListSize = 100

type TemplateController struct {
	pb.UnsafeEmailTemplateServiceServer
	repo repository.Template
}

func (c *TemplateController) CreateEmailTemplate(
	ctx context.Context,
	request *pb.CreateEmailTemplateRequest,
) (*pb.CreateEmailTemplateResponse, error) {
	validate := func() error {
		err := validation.ValidateStructWithContext(ctx, request,
			validation.Field(&request.Template, validation.Required),
		)
		if err != nil {
			return err
		}
		return validation.ValidateStructWithContext(ctx, request.Template,
			validation.Field(&request.Template.WhiteLabelId, validation.Required, is.UUID),
			validation.Field(&request.Template.Language, validation.Required),
			validation.Field(&request.Template.Type, validation.Required),
			validation.Field(&request.Template.Template, validation.Required),
			validation.Field(&request.Template.Subject, validation.Required),
			validation.Field(&request.Template.Footer, validation.Required),
		)
	}
	if err := validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	template, err := c.pbEmailTemplateToEmailTemplate(request.Template)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	if err = validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	err = c.repo.Create(ctx, template)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.CreateEmailTemplateResponse{}, nil
}

func (c *TemplateController) GetEmailTemplate(
	ctx context.Context,
	request *pb.GetEmailTemplateRequest,
) (*pb.GetEmailTemplateResponse, error) {
	err := validation.ValidateStructWithContext(ctx, request,
		validation.Field(&request.WhiteLabelId, validation.Required, is.UUID),
		validation.Field(&request.Language, validation.Required),
		validation.Field(&request.Type, validation.Required),
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	whiteLabelId, err := uuid.Parse(request.WhiteLabelId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	_type := c.pbTemplateTypeToCodeTemplate(request.Type)
	template, err := c.repo.Get(ctx, whiteLabelId, request.Language, _type)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.GetEmailTemplateResponse{
		Template: c.emailTemplateToPbEmailTemplate(template),
	}, nil
}

func (c *TemplateController) UpdateEmailTemplate(
	ctx context.Context,
	request *pb.UpdateEmailTemplateRequest,
) (*pb.UpdateEmailTemplateResponse, error) {
	validate := func() error {
		err := validation.ValidateStructWithContext(ctx, request,
			validation.Field(&request.Template, validation.Required),
		)
		if err != nil {
			return err
		}
		return validation.ValidateStructWithContext(ctx, request.Template,
			validation.Field(&request.Template.WhiteLabelId, validation.Required, is.UUID),
			validation.Field(&request.Template.Language, validation.Required),
			validation.Field(&request.Template.Type, validation.Required),
			validation.Field(&request.Template.Template, validation.Required),
			validation.Field(&request.Template.Subject, validation.Required),
			validation.Field(&request.Template.Footer, validation.Required),
		)
	}
	if err := validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	template, err := c.pbEmailTemplateToEmailTemplate(request.Template)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	err = c.repo.Update(ctx, template)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.UpdateEmailTemplateResponse{}, nil
}

func (c *TemplateController) DeleteEmailTemplate(
	ctx context.Context,
	request *pb.DeleteEmailTemplateRequest,
) (*pb.DeleteEmailTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTemplate not implemented")
}

func (c *TemplateController) ListEmailTemplate(
	ctx context.Context,
	request *pb.ListEmailTemplateRequest,
) (*pb.ListEmailTemplateResponse, error) {
	var err error
	pg := repository.NewPagination(defaultTemplateListSize)
	if request.PageToken != "" {
		pg, err = paginationFromBase64(request.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
		}
	}

	l, err := c.repo.List(ctx, pg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	var nextPageToken string
	if len(l) != 0 {
		nextPageToken, err = paginationToBase64(pg.Next())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}

	return &pb.ListEmailTemplateResponse{
		NextPageToken: nextPageToken,
		Templates:     c.emailTemplatesToPbEmailTemplates(l),
	}, nil
}

func (c *TemplateController) pbTemplateTypeToCodeTemplate(d pb.TemplateType) model.CodeTemplate {
	switch d {
	case pb.TemplateType_WorkerChangedState:
		return model.WorkerChangedState
	case pb.TemplateType_PasswordRestoration:
		return model.PasswordRestoration
	case pb.TemplateType_WalletChangedAddress:
		return model.WalletChangedAddress
	case pb.TemplateType_Register:
		return model.Register
	case pb.TemplateType_UserHashrateDecreased:
		return model.UserHashrateDecreased
	case pb.TemplateType_Mobile2faOn:
		return model.MobileTwoFaOn
	case pb.TemplateType_Mobile2faOff:
		return model.MobileTwoFaOff
	case pb.TemplateType_ChangePassword:
		return model.PasswordChange
	case pb.TemplateType_DeletePhone:
		return model.PhoneDelete
	case pb.TemplateType_ReferralRewardPayouts:
		return model.ReferralRewardPayouts
	case pb.TemplateType_WorkerReport:
		return model.WorkerReport
	case pb.TemplateType_SwapSupportMessage:
		return model.SwapSupport
	case pb.TemplateType_SwapMessage:
		return model.SwapMessage
	case pb.TemplateType_SwapSuccessful:
		return model.SwapSuccessful
	case pb.TemplateType_IncomeReport:
		return model.IncomeReport
	case pb.TemplateType_GoggleTwoFaOn:
		return model.GoggleTwoFaOn
	case pb.TemplateType_GoggleTwoFaOff:
		return model.GoggleTwoFaOff
	case pb.TemplateType_ChangeEmail:
		return model.ChangeEmail
	case pb.TemplateType_PayoutReport:
		return model.PayoutReport
	case pb.TemplateType_Unknown:
		return ""
	default:
		return ""
	}
}

func (c *TemplateController) codeTemplateToPBTemplateType(d model.CodeTemplate) pb.TemplateType {
	switch d {
	case model.WorkerChangedState:
		return pb.TemplateType_WorkerChangedState
	case model.PasswordRestoration:
		return pb.TemplateType_PasswordRestoration
	case model.WalletChangedAddress:
		return pb.TemplateType_WalletChangedAddress
	case model.Register:
		return pb.TemplateType_Register
	case model.UserHashrateDecreased:
		return pb.TemplateType_UserHashrateDecreased
	case model.MobileTwoFaOn:
		return pb.TemplateType_Mobile2faOn
	case model.MobileTwoFaOff:
		return pb.TemplateType_Mobile2faOff
	case model.PasswordChange:
		return pb.TemplateType_ChangePassword
	case model.PhoneDelete:
		return pb.TemplateType_DeletePhone
	case model.ReferralRewardPayouts:
		return pb.TemplateType_ReferralRewardPayouts
	case model.WorkerReport:
		return pb.TemplateType_WorkerReport
	case model.SwapSupport:
		return pb.TemplateType_SwapSupportMessage
	case model.SwapMessage:
		return pb.TemplateType_SwapMessage
	case model.SwapSuccessful:
		return pb.TemplateType_SwapSuccessful
	case model.IncomeReport:
		return pb.TemplateType_IncomeReport
	case model.GoggleTwoFaOn:
		return pb.TemplateType_GoggleTwoFaOn
	case model.GoggleTwoFaOff:
		return pb.TemplateType_GoggleTwoFaOff
	case model.ChangeEmail:
		return pb.TemplateType_ChangeEmail
	case model.PayoutReport:
		return pb.TemplateType_PayoutReport
	default:
		return pb.TemplateType_Unknown
	}
}

func (c *TemplateController) pbEmailTemplateToEmailTemplate(t *pb.EmailTemplate) (model.Template, error) {
	whiteLabelID, err := uuid.Parse(t.WhiteLabelId)
	if err != nil {
		return model.Template{}, fmt.Errorf("parse whitelable_id: %w", err)
	}
	return model.Template{
		WhiteLabelID: whiteLabelID,
		Language:     t.Language,
		Template:     model.NewTextTemplate(t.Template),
		Type:         c.pbTemplateTypeToCodeTemplate(t.Type),
		Subject:      model.NewTextTemplate(t.Subject),
		Footer:       t.Footer,
	}, nil
}

func (c *TemplateController) emailTemplateToPbEmailTemplate(t model.Template) *pb.EmailTemplate {
	return &pb.EmailTemplate{
		WhiteLabelId: t.WhiteLabelID.String(),
		Language:     t.Language,
		Template:     t.Template.String(),
		Type:         c.codeTemplateToPBTemplateType(t.Type),
		Subject:      t.Subject.String(),
		Footer:       t.Footer,
	}
}

func (c *TemplateController) emailTemplatesToPbEmailTemplates(l []model.Template) []*pb.EmailTemplate {
	r := make([]*pb.EmailTemplate, 0, len(l))
	for i := range l {
		r = append(r, c.emailTemplateToPbEmailTemplate(l[i]))
	}
	return r
}

func NewTemplateController(repo repository.Template) *TemplateController {
	return &TemplateController{
		repo: repo,
	}
}
