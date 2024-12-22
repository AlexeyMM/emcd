package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/email/internal/berrors"
	"code.emcdtech.com/emcd/service/email/internal/mail_sender"
	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
)

type ProviderFactory func(p model.Provider) (EmailProvider, error)

type EmailProvider interface {
	SendHtml(ctx context.Context, mail model.Mail) error
}

type EmailSender struct {
	templateRepo        repository.Template
	provideSettingsRepo repository.ProvideSettings
	emailMessageRepo    repository.EmailMessages
	providers           map[model.ProviderName]ProviderFactory
}

func (s *EmailSender) Send(
	ctx context.Context,
	whiteLabelID uuid.UUID,
	language string,
	_type model.CodeTemplate,
	to string,
	placeholder any,
	attachments ...model.Attachment,
) error {
	tpl, err := s.templateRepo.Get(ctx, whiteLabelID, language, _type)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("get template: %w", berrors.ErrTemplateNotFound)
		}
		return fmt.Errorf("get template: %w", err)
	}
	subject, err := tpl.Subject.Render(placeholder)
	if err != nil {
		return fmt.Errorf("render subject: %w", err)
	}
	body, err := tpl.Template.Render(placeholder)
	if err != nil {
		return fmt.Errorf("render body: %w", err)
	}

	mail := model.Mail{
		To:          to,
		Subject:     subject,
		Message:     body,
		Attachments: attachments,
	}

	var wasSend bool
	ps, err := s.provideSettingsRepo.Get(ctx, whiteLabelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("get provider settings: %w", berrors.ErrProviderSettingNotFound)
		}
		return fmt.Errorf("get provider settings: %w", err)
	}
	for _, p := range ps.Providers {
		newProvider, ok := s.providers[p.Name]
		if !ok {
			return fmt.Errorf("provider %s not found", p.Name)
		}
		sender, err := newProvider(p) //nolint:govet
		if err != nil {
			log.Error(ctx, "get provider %s: %s", p.Name, err.Error())
			continue
		}
		err = sender.SendHtml(ctx, mail)
		if err != nil {
			log.Warn(ctx, "sending %s: %s", p.Name, err.Error())
			continue
		}
		wasSend = true
		break
	}

	if wasSend {
		err = s.emailMessageRepo.Create(ctx, &model.EmailMessageEvent{
			Email:     mail.To,
			Type:      _type,
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Error(ctx, "save sent email (to:  %s, template: %s): %s", mail.To, _type, err.Error())
		}
		return nil
	}
	return fmt.Errorf("not send email to %s: %w", mail.To, berrors.ErrNoMailSenderAvailable)
}

func NewEmailSender(
	templateRepo repository.Template,
	provideSettingsRepo repository.ProvideSettings,
	emailMessageRepo repository.EmailMessages,
) *EmailSender {
	return &EmailSender{
		templateRepo:        templateRepo,
		provideSettingsRepo: provideSettingsRepo,
		emailMessageRepo:    emailMessageRepo,
		providers: map[model.ProviderName]ProviderFactory{
			model.SMTPProviderName: func(p model.Provider) (EmailProvider, error) {
				s, err := model.GetSetting[model.SmtpSetting](p)
				if err != nil {
					return nil, fmt.Errorf("get setting: %w", err)
				}
				return mail_sender.NewSmtp(s), nil
			},
			model.MailgunProviderName: func(p model.Provider) (EmailProvider, error) {
				s, err := model.GetSetting[model.MailgunSetting](p)
				if err != nil {
					return nil, fmt.Errorf("get setting: %w", err)
				}
				return mail_sender.NewMailgun(s), nil
			},
		},
	}
}
