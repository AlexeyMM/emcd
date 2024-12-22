package mail_sender

import (
	"context"
	"fmt"
	"net/mail"

	mailgunLib "github.com/mailgun/mailgun-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

type Mailgun struct {
	setting model.MailgunSetting
}

func NewMailgun(setting model.MailgunSetting) *Mailgun {
	return &Mailgun{
		setting: setting,
	}
}

func (s *Mailgun) SendHtml(ctx context.Context, m model.Mail) (err error) {
	tracer := otel.Tracer("send-email")
	_, span := tracer.Start(ctx, s.Name(), trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	mg := mailgunLib.NewMailgun(s.setting.Domain, s.setting.ApiKey)
	mg.SetAPIBase(s.setting.ApiBase)
	from := mail.Address{
		Name:    s.setting.FromAddressDisplayedAs,
		Address: s.setting.FromAddress,
	}
	message := mg.NewMessage(from.String(), m.Subject, "", m.To)
	message.SetHtml(m.Message)
	for _, attachment := range m.Attachments {
		message.AddBufferAttachment(attachment.Filename, attachment.Data)
	}
	_, _, err = mg.Send(message)
	if err != nil {
		return fmt.Errorf("mailgun mail sender: send mail: %w", err)
	}
	return nil
}

func (s *Mailgun) Name() string {
	return "mailgun"
}
