package mail_sender

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"net/mail"
	smtpLib "net/smtp"
	"path/filepath"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

type Smtp struct {
	setting model.SmtpSetting
}

func NewSmtp(setting model.SmtpSetting) *Smtp {
	return &Smtp{
		setting: setting,
	}
}

func (s *Smtp) SendHtml(ctx context.Context, m model.Mail) (err error) {
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

	auth := smtpLib.PlainAuth("", s.setting.Username, s.setting.Password, s.setting.ServerAddress)
	hostAndPort := fmt.Sprintf("%s:%d", s.setting.ServerAddress, s.setting.ServerPort)

	from := mail.Address{
		Name:    s.setting.FromAddressDisplayedAs,
		Address: s.setting.FromAddress,
	}

	rawMessage, err := formMessage(from.String(), m)
	if err != nil {
		return fmt.Errorf("make raw message: %w", err)
	}

	err = smtpLib.SendMail(
		hostAndPort,
		auth,
		s.setting.FromAddress,
		[]string{m.To},
		rawMessage,
	)
	if err != nil {
		return fmt.Errorf("smtp mail sendeer: send mail: %w", err)
	}
	return nil
}

func (s *Smtp) Name() string {
	return "smtp" // TODO (a.barsukovskij) maybe added address for identification???
}

func formMessage(from string, mail model.Mail) ([]byte, error) {
	header := make(map[string]string)
	header["From"] = from
	header["To"] = mail.To
	header["Subject"] = mail.Subject
	header["MIME-Version"] = "1.0"

	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)
	header["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()

	// Write headers to message buffer
	for k, v := range header {
		_, err := msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		if err != nil {
			return nil, fmt.Errorf("write header: %w", err)
		}
	}
	_, err := msg.WriteString("\r\n")
	if err != nil {
		return nil, fmt.Errorf("write message delimiter: %w", err)
	}

	// Write the email body
	bodyWriter, err := writer.CreatePart(map[string][]string{"Content-Type": {"text/html; charset=UTF-8"}})
	if err != nil {
		return nil, fmt.Errorf("write message header: %w", err)
	}
	_, err = bodyWriter.Write([]byte(mail.Message))
	if err != nil {
		return nil, fmt.Errorf("write message body: %w", err)
	}

	// Attach files
	for _, attachment := range mail.Attachments {
		if err := writeAttachment(writer, attachment); err != nil {
			return nil, fmt.Errorf("write attachment: %w", err)
		}
	}
	_ = writer.Close()

	return msg.Bytes(), nil
}

// writeAttachment write attachment to the email message.
func writeAttachment(writer *multipart.Writer, attachment model.Attachment) error {
	// Get the filename and MIME type
	filename := filepath.Base(attachment.Filename)
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	part, err := writer.CreatePart(map[string][]string{
		"Content-Disposition":       {"attachment; filename=" + filename},
		"Content-Type":              {mimeType},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return err
	}

	// Encode the file content as base64
	encoder := base64.NewEncoder(base64.StdEncoding, part)
	defer func() {
		_ = encoder.Close()
	}()
	_, err = encoder.Write(attachment.Data)
	if err != nil {
		return fmt.Errorf("write encoder: %w", err)
	}
	return nil
}
