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

const defaultSettingsListSize = 100

type SettingController struct {
	pb.UnsafeEmailProviderSettingsServiceServer
	repo repository.ProvideSettings
}

func (c *SettingController) CreateEmailProviderSetting(
	ctx context.Context,
	request *pb.CreateEmailProviderSettingRequest,
) (*pb.CreateEmailProviderSettingResponse, error) {
	validate := func() error {
		err := validation.ValidateStructWithContext(ctx, request,
			validation.Field(&request.Setting.WhiteLabelId, validation.Required),
			validation.Field(&request.Setting, validation.Required),
		)
		if err != nil {
			return err
		}
		err = validation.ValidateStructWithContext(ctx, &request.Setting,
			validation.Field(&request.Setting.Providers, validation.Required, validation.Each(
				validation.By(func(value interface{}) error {
					provider, ok := value.(*pb.EmailProviderSetting_Provider)
					if !ok {
						return fmt.Errorf("not support type %T", value)
					}
					switch v := provider.Value.(type) {
					case *pb.EmailProviderSetting_Provider_Smtp:
						return validation.ValidateStructWithContext(ctx, v.Smtp,
							validation.Field(&v.Smtp.Username, validation.Required),
							validation.Field(&v.Smtp.Password, validation.Required),
							validation.Field(&v.Smtp.Host, validation.Required),
							validation.Field(&v.Smtp.Port, validation.Required),
							validation.Field(&v.Smtp.FromAddress, validation.Required),
							validation.Field(&v.Smtp.FromAddressDisplayedAs, validation.Required),
							validation.Field(&v.Smtp.Username, validation.Required),
						)
					case *pb.EmailProviderSetting_Provider_Mailgun:
						return validation.ValidateStructWithContext(ctx, v.Mailgun,
							validation.Field(&v.Mailgun.Domain, validation.Required),
							validation.Field(&v.Mailgun.ApiKey, validation.Required),
							validation.Field(&v.Mailgun.ApiBase, validation.Required),
							validation.Field(&v.Mailgun.FromAddress, validation.Required),
							validation.Field(&v.Mailgun.FromAddressDisplayedAs, validation.Required),
						)
					default:
						return fmt.Errorf("unexpected validation for type %T", provider)
					}
				}),
			)),
		)
		if err != nil {
			return err
		}
		return nil
	}
	if err := validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	setting, err := c.pbSettingToSetting(request.Setting)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	err = c.repo.Create(ctx, setting)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.CreateEmailProviderSettingResponse{}, nil
}

func (c *SettingController) GetEmailProviderSetting(
	ctx context.Context,
	request *pb.GetEmailProviderSettingRequest,
) (*pb.GetEmailProviderSettingResponse, error) {
	err := validation.ValidateStructWithContext(ctx, request,
		validation.Field(&request.WhiteLabelId, validation.Required, is.UUID),
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	whiteLabelId, err := uuid.Parse(request.WhiteLabelId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "parse WhiteLableId: %v", err)
	}
	setting, err := c.repo.Get(ctx, whiteLabelId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	s, err := c.settingToPbSetting(setting)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.GetEmailProviderSettingResponse{
		Setting: s,
	}, nil
}

func (c *SettingController) UpdateEmailProviderSetting(
	ctx context.Context,
	request *pb.UpdateEmailProviderSettingRequest,
) (*pb.UpdateEmailProviderSettingResponse, error) {
	validate := func() error {
		err := validation.ValidateStructWithContext(ctx, request,
			validation.Field(&request.Setting.WhiteLabelId, validation.Required, is.UUID),
			validation.Field(&request.Setting, validation.Required),
		)
		if err != nil {
			return err
		}
		err = validation.ValidateStructWithContext(ctx, &request.Setting,
			validation.Field(&request.Setting.Providers, validation.Required, validation.Each(
				validation.By(
					func(value interface{}) error {
						provider, ok := value.(*pb.EmailProviderSetting_Provider)
						if !ok {
							return fmt.Errorf("not support type %T", value)
						}
						switch v := provider.Value.(type) {
						case *pb.EmailProviderSetting_Provider_Smtp:
							return validation.ValidateStructWithContext(ctx, v.Smtp,
								validation.Field(&v.Smtp.Username, validation.Required),
								validation.Field(&v.Smtp.Password, validation.Required),
								validation.Field(&v.Smtp.Host, validation.Required),
								validation.Field(&v.Smtp.Port, validation.Required),
								validation.Field(&v.Smtp.FromAddress, validation.Required),
								validation.Field(&v.Smtp.FromAddressDisplayedAs, validation.Required),
								validation.Field(&v.Smtp.Username, validation.Required),
							)
						case *pb.EmailProviderSetting_Provider_Mailgun:
							return validation.ValidateStructWithContext(ctx, v.Mailgun,
								validation.Field(&v.Mailgun.Domain, validation.Required),
								validation.Field(&v.Mailgun.ApiKey, validation.Required),
								validation.Field(&v.Mailgun.ApiBase, validation.Required),
								validation.Field(&v.Mailgun.FromAddress, validation.Required),
								validation.Field(&v.Mailgun.FromAddressDisplayedAs, validation.Required),
							)
						default:
							return fmt.Errorf("unexpected validation for type %T", provider)
						}
					},
				),
			),
			),
		)
		return err
	}
	if err := validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	setting, err := c.pbSettingToSetting(request.Setting)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	err = c.repo.Update(ctx, setting)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.UpdateEmailProviderSettingResponse{}, nil
}

func (c *SettingController) DeleteEmailProviderSetting(
	ctx context.Context,
	request *pb.DeleteEmailProviderSettingRequest,
) (*pb.DeleteEmailProviderSettingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSetting not implemented")
}

func (c *SettingController) ListEmailProviderSettings(
	ctx context.Context,
	request *pb.ListEmailProviderSettingsRequest,
) (*pb.ListEmailProviderSettingsResponse, error) {
	var err error
	pg := repository.NewPagination(defaultSettingsListSize)
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
	settings, err := c.settingsToPbSettings(l)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.ListEmailProviderSettingsResponse{
		NextPageToken: nextPageToken,
		Settings:      settings,
	}, nil
}

func (c *SettingController) pbSettingToSetting(s *pb.EmailProviderSetting) (model.Setting, error) {
	whiteLabelID, err := uuid.Parse(s.WhiteLabelId)
	if err != nil {
		return model.Setting{}, fmt.Errorf("parse WhiteLabelId: %w", err)
	}
	r := model.Setting{
		WhiteLabelID: whiteLabelID,
	}

	return r, nil
}

func (c *SettingController) settingToPbSetting(s model.Setting) (*pb.EmailProviderSetting, error) {
	r := &pb.EmailProviderSetting{
		WhiteLabelId: s.WhiteLabelID.String(),
	}
	for _, provider := range s.Providers {
		switch provider.Name {
		case model.SMTPProviderName:
			s, err := model.GetSetting[model.SmtpSetting](provider)
			if err != nil {
				return nil, fmt.Errorf("get setting: %w", err)
			}
			r.Providers = append(r.Providers,
				&pb.EmailProviderSetting_Provider{
					Value: &pb.EmailProviderSetting_Provider_Smtp{
						Smtp: &pb.EmailProviderSetting_SMTP{
							Username:               s.Username,
							Password:               s.Password,
							Host:                   s.ServerAddress,
							Port:                   int32(s.ServerPort),
							FromAddress:            s.FromAddress,
							FromAddressDisplayedAs: s.FromAddressDisplayedAs,
						},
					},
				},
			)
		case model.MailgunProviderName:
			s, err := model.GetSetting[model.MailgunSetting](provider)
			if err != nil {
				return nil, fmt.Errorf("get setting: %w", err)
			}
			r.Providers = append(r.Providers,
				&pb.EmailProviderSetting_Provider{
					Value: &pb.EmailProviderSetting_Provider_Mailgun{
						Mailgun: &pb.EmailProviderSetting_Mailgun{
							Domain:                 s.Domain,
							ApiKey:                 s.ApiKey,
							ApiBase:                s.ApiBase,
							FromAddress:            s.FromAddress,
							FromAddressDisplayedAs: s.FromAddressDisplayedAs,
						},
					},
				},
			)
		default:
			return nil, fmt.Errorf("unknown provider: %s", provider.Name)
		}
	}
	return r, nil
}

func (c *SettingController) settingsToPbSettings(l []model.Setting) ([]*pb.EmailProviderSetting, error) {
	r := make([]*pb.EmailProviderSetting, 0, len(l))
	for i := range l {
		s, err := c.settingToPbSetting(l[i])
		if err != nil {
			return nil, fmt.Errorf("get setting: %w", err)
		}
		r = append(r, s)
	}
	return r, nil
}

func NewSettingsController(
	repo repository.ProvideSettings,
) *SettingController {
	return &SettingController{
		repo: repo,
	}
}
