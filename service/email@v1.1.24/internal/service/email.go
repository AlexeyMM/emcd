package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
)

const (
	aliveStateEn   = "On"
	deadStateEn    = "Off"
	aliveStateRu   = "Включен"
	deadStateRu    = "Выключен"
	english        = "en"
	russian        = "ru"
	eventWlVersion = 2
)

type Email interface {
	ListMessages(
		ctx context.Context,
		email *string,
		eventType *string,
		skip, take int32) ([]*model.EmailMessageEvent, int, error)
	SendRegister(ctx context.Context, wlID uuid.UUID, email, domain, token, language string) error
	SendMobileTwoFaOff(ctx context.Context, userID uuid.UUID, domain string) error
	SendMobileTwoFaOn(ctx context.Context, userID uuid.UUID, token, domain string) error
	SendGoogleTwoFaOff(ctx context.Context, userID uuid.UUID, domain string) error
	SendGoogleTwoFaOn(ctx context.Context, userID uuid.UUID, token, domain string) error
	SendPasswordChange(ctx context.Context, userID uuid.UUID, token, domain string) error
	SendPhoneDelete(ctx context.Context, userID uuid.UUID, token, domain string) error
	SendWalletChangedAddress(
		ctx context.Context,
		userID uuid.UUID,
		domain, token, coinCode string,
	) error
	SendPasswordRestoration(
		ctx context.Context,
		userID uuid.UUID, token, domain string) error
	SendWorkerChangedState(
		ctx context.Context,
		worker *model.Worker,
		email, domain string,
		whiteLabelID uuid.UUID,
		language string) error
	SendUserHashrateDecreased(
		ctx context.Context,
		email, domain string,
		whiteLabelID uuid.UUID,
		language string,
		decreasedBy decimal.Decimal,
		coin string) error
	SendReferralRewardPayouts(
		ctx context.Context,
		userID uuid.UUID,
		domain string,
		from, to time.Time,
		attachments []model.Attachment,
	) error
	SendChangeEmail(ctx context.Context, whiteLabelID uuid.UUID, email, domain, token, language string) error
	SendStatisticsReport(
		ctx context.Context,
		email, reportLink, language string,
		reportType model.CodeTemplate,
	) error
	SendSwapSupportMessage(
		ctx context.Context,
		name, email, text string,
	) error
	SendInitialSwapMessage(
		ctx context.Context,
		email, language, link string,
	) error
	SendSuccessfulSwapMessage(
		ctx context.Context,
		swapID uuid.UUID,
		from, to, address, email, language, executionTime string,
	) error
}

type email struct {
	emailSender            *EmailSender
	emailMessages          repository.EmailMessages
	whitelabel             repository.Whitelabel
	whiteLabelEventClients repository.WhiteLabelEventClients
	profile                repository.Profile
	links                  map[model.CodeTemplate]string
	domains                map[string]string
}

func NewEmail(
	emailSender *EmailSender,
	emailMessages repository.EmailMessages,
	whitelabel repository.Whitelabel,
	whiteLabelEventClients repository.WhiteLabelEventClients,
	profile repository.Profile,
	links map[model.CodeTemplate]string,
	domains map[string]string,
) *email {
	return &email{
		emailSender:            emailSender,
		emailMessages:          emailMessages,
		whitelabel:             whitelabel,
		whiteLabelEventClients: whiteLabelEventClients,
		profile:                profile,
		links:                  links,
		domains:                domains,
	}
}

func (e *email) SendPasswordRestoration(
	ctx context.Context,
	userID uuid.UUID, token, domain string,
) error {
	const op = "service.email.SendPasswordRestoration"
	user, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	placeholder := struct {
		Token  string
		Domain string
	}{
		Token:  token,
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if user.WhiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
		wlDomain, err := e.getWlDomain(ctx, user.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, user.WhiteLabelID, user.Language, model.PasswordRestoration, user.Email, placeholder)
}

func (e *email) getDomain(ctx context.Context, host string) string {
	domain, ok := e.domains[host]
	if !ok || domain == "" {
		log.Warn(ctx, "param domain (`%s`) not found", host)
		domain = e.domains[""]
	}
	return domain
}

func (e *email) getWlDomain(ctx context.Context, wlID uuid.UUID) (string, error) {
	wl, err := e.whitelabel.GetWlByID(ctx, wlID)
	if err != nil {
		return "", err
	}
	origin, err := e.whitelabel.GetWlFullDomain(ctx, wl.UserId)
	if err != nil {
		log.Error(ctx, "failed to find wl: %s domain, err: %s", wlID.String(), err.Error())
		return "", err
	}
	return origin, nil
}

func (e *email) SendWorkerChangedState(
	ctx context.Context,
	worker *model.Worker,
	email, domain string,
	whiteLabelID uuid.UUID,
	language string,
) error {
	var state string
	switch language {
	case english:
		if worker.IsOn {
			state = aliveStateEn
		} else {
			state = deadStateEn
		}
	case russian:
		if worker.IsOn {
			state = aliveStateRu
		} else {
			state = deadStateRu
		}
	}
	placeholder := struct {
		Username       string
		WorkerName     string
		Coin           string
		State          string
		StateChangedAt time.Time
		Domain         string
	}{
		Username:       worker.Username,
		WorkerName:     worker.Name,
		Coin:           worker.Coin,
		State:          state,
		StateChangedAt: worker.StateChangedAt,
		Domain:         e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if whiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, false)
		wlDomain, err := e.getWlDomain(ctx, whiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, whiteLabelID, language, model.WorkerChangedState, email, placeholder)
}

func (e *email) SendWalletChangedAddress(
	ctx context.Context,
	userID uuid.UUID,
	domain, token, coinCode string,
) error {
	const op = "service.email.SendWalletChangedAddress"
	user, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service.email.SendWalletChangedAddress: request profile: %w", err)
	}
	var (
		strategy    SendEmailStrategy
		placeholder any
		userEmail   string
	)
	if user.WhiteLabelID != uuid.Nil {
		wlDomain, err := e.getWlDomain(ctx, user.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder = struct {
			Domain string
			Token  string
		}{
			Domain: wlDomain,
			Token:  token,
		}
		userEmail = user.Email
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
	} else {
		settings, err := e.profile.GetNotificationSettings(ctx, userID)
		if err != nil {
			return fmt.Errorf("%s: get notification settings: %w", op, err)
		}
		userEmail = settings.Email
		possibleDomain := e.getDomain(ctx, domain)
		url := strings.ReplaceAll(e.links[model.WalletChangedAddress], "{{.AutopayAddressToken}}", token)
		url = strings.ReplaceAll(url, "{{.Domain}}", possibleDomain)
		placeholder = struct {
			Coin       string
			Domain     string
			PoolWebUrl string
			Year       int
		}{
			Coin:       strings.ToUpper(coinCode),
			Domain:     possibleDomain,
			Year:       time.Now().UTC().Year(),
			PoolWebUrl: url,
		}
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, user.WhiteLabelID, user.Language, model.WalletChangedAddress, userEmail, placeholder)
}

func (e *email) SendRegister(ctx context.Context, whiteLabelID uuid.UUID, email, domain, token, language string) error {
	placeholder := struct {
		Token  string
		Domain string
	}{
		Token:  token,
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if whiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, false)
		wlDomain, err := e.getWlDomain(ctx, whiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, whiteLabelID, language, model.Register, email, placeholder)
}

func (e *email) SendGoogleTwoFaOff(ctx context.Context, userID uuid.UUID, domain string) error {
	const op = "service.email.SendGoogleTwoFaOff"
	pr, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: get profile: %w", op, err)
	}
	placeholder := struct {
		Year   int
		Domain string
	}{
		Year:   time.Now().UTC().Year(),
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if pr.WhiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
		wlDomain, err := e.getWlDomain(ctx, pr.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(
		ctx,
		pr.WhiteLabelID,
		pr.Language,
		model.GoggleTwoFaOff,
		pr.Email,
		placeholder,
	)
}

func (e *email) SendGoogleTwoFaOn(ctx context.Context, userID uuid.UUID, token, domain string) error {
	pr, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service.email.SendGoogleTwoFaOn: request profile: %w", err)
	}
	placeholder := struct {
		Token  string
		Year   int
		Domain string
	}{
		Token:  token,
		Year:   time.Now().UTC().Year(),
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if pr.WhiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
		wlDomain, err := e.getWlDomain(ctx, pr.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, pr.WhiteLabelID, pr.Language, model.GoggleTwoFaOn, pr.Email, placeholder)
}

func (e *email) SendMobileTwoFaOff(ctx context.Context, userID uuid.UUID, domain string) error {
	pr, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service.email.SendMobileTwoFaOff: request profile: %w", err)
	}
	placeholder := struct {
		Year   int
		Domain string
	}{
		Year:   time.Now().UTC().Year(),
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if pr.WhiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
		wlDomain, err := e.getWlDomain(ctx, pr.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, pr.WhiteLabelID, pr.Language, model.MobileTwoFaOff, pr.Email, placeholder)
}

func (e *email) SendMobileTwoFaOn(ctx context.Context, userID uuid.UUID, token, domain string) error {
	pr, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service.email.SendMobileTwoFaOn: request profile: %w", err)
	}
	placeholder := struct {
		Token  string
		Year   int
		Domain string
	}{
		Token:  token,
		Year:   time.Now().UTC().Year(),
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if pr.WhiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
		wlDomain, err := e.getWlDomain(ctx, pr.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, pr.WhiteLabelID, pr.Language, model.MobileTwoFaOn, pr.Email, placeholder)
}

func (e *email) SendUserHashrateDecreased(
	ctx context.Context,
	email, domain string,
	whiteLabelID uuid.UUID,
	language string,
	decreasedBy decimal.Decimal, coin string,
) error {
	placeholder := struct {
		DecreasedByPercent string
		Coin               string
		Domain             string
	}{
		DecreasedByPercent: toPercent(decreasedBy).String(),
		Coin:               coin,
		Domain:             e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if whiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, false)
		wlDomain, err := e.getWlDomain(ctx, whiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, whiteLabelID, language, model.UserHashrateDecreased, email, placeholder)
}

func toPercent(d decimal.Decimal) decimal.Decimal {
	return d.Mul(decimal.NewFromInt(100)).Round(2)
}

func (e *email) ListMessages(ctx context.Context,
	email, eventType *string,
	skip, take int32,
) ([]*model.EmailMessageEvent, int, error) {
	if take <= 0 {
		return nil, 0, fmt.Errorf("getReferrals: Must be greater than 0")
	}

	return e.emailMessages.ListMessages(ctx, email, eventType, skip, take)
}

func (e *email) SendPasswordChange(ctx context.Context, userID uuid.UUID, token, domain string) error {
	pr, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service.email.SendPasswordChange: %w", err)
	}
	placeholder := struct {
		Token  string
		Domain string
	}{
		Token:  token,
		Domain: e.getDomain(ctx, domain),
	}
	var strategy SendEmailStrategy
	if pr.WhiteLabelID != uuid.Nil {
		strategy = newWLSendEmailStrategy(e.emailSender, e.whitelabel, true)
		wlDomain, err := e.getWlDomain(ctx, pr.WhiteLabelID)
		if err != nil {
			return err
		}
		placeholder.Domain = wlDomain
	} else {
		strategy = newPoolSendEmailStrategy(e.emailSender)
	}
	return strategy.SendEmail(ctx, pr.WhiteLabelID, pr.Language, model.PasswordChange, pr.Email, placeholder)
}

func (e *email) SendPhoneDelete(ctx context.Context, userID uuid.UUID, token, domain string) error {
	pr, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service.email.SendPhoneDelete: %w", err)
	}
	placeholder := struct {
		Token  string
		Year   int
		Domain string
	}{
		Token:  token,
		Year:   time.Now().Year(),
		Domain: e.getDomain(ctx, domain),
	}
	return e.emailSender.Send(ctx, uuid.Nil, pr.Language, model.PhoneDelete, pr.Email, placeholder)
}

func (e *email) SendReferralRewardPayouts(
	ctx context.Context,
	userID uuid.UUID,
	domain string,
	from, to time.Time,
	attachments []model.Attachment,
) error {
	const op = "service.email.SendReferralRewardPayouts"
	user, err := e.profile.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: get profile: %w", op, err)
	}
	placeholder := struct {
		From   time.Time
		To     time.Time
		Domain string
	}{
		From:   from,
		To:     to,
		Domain: e.getDomain(ctx, domain),
	}

	settings, err := e.profile.GetNotificationSettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: get notification settings: %w", op, err)
	}
	return e.emailSender.Send(
		ctx,
		uuid.Nil, // use Nil for backward compatibility
		user.Language,
		model.ReferralRewardPayouts,
		settings.Email,
		placeholder,
		attachments...,
	)
}

func (e *email) SendChangeEmail(
	ctx context.Context,
	whiteLabelID uuid.UUID,
	email, domain, token, language string,
) error {
	placeholder := struct {
		Token  string
		Domain string
	}{
		Token:  token,
		Domain: e.getDomain(ctx, domain),
	}
	return e.emailSender.Send(ctx, whiteLabelID, language, model.ChangeEmail, email, placeholder)
}

func (e *email) SendStatisticsReport(
	ctx context.Context,
	email, reportLink, language string,
	reportType model.CodeTemplate,
) error {
	placeholder := struct {
		ReportLink string
	}{
		reportLink,
	}
	return e.emailSender.Send(ctx, uuid.Nil, language, reportType, email, placeholder)
}

func (e *email) SendSwapSupportMessage(
	ctx context.Context,
	name, email, text string,
) error {
	placeholder := struct {
		Name  string
		Email string
		Text  string
	}{
		Name:  name,
		Email: email,
		Text:  text,
	}
	return e.emailSender.Send(ctx, uuid.Nil, "ru", model.SwapSupport, "support@emcd.io", placeholder)
}

func (e *email) SendInitialSwapMessage(
	ctx context.Context,
	email, language, link string,
) error {
	placeholder := struct {
		Email    string
		Language string
		Link     string
	}{
		Email:    email,
		Language: language,
		Link:     link,
	}

	return e.emailSender.Send(ctx, uuid.Nil, language, model.SwapMessage, email, placeholder)
}

func (e *email) SendSuccessfulSwapMessage(
	ctx context.Context,
	swapID uuid.UUID,
	from, to, address, email, language, executionTime string,
) error {
	placeholder := struct {
		SwapID        uuid.UUID
		From          string
		To            string
		Address       string
		ExecutionTime string
	}{
		SwapID:        swapID,
		From:          from,
		To:            to,
		Address:       address,
		ExecutionTime: executionTime,
	}
	return e.emailSender.Send(ctx, uuid.Nil, language, model.SwapSuccessful, email, placeholder)
}

type SendEmailStrategy interface {
	SendEmail(ctx context.Context,
		whiteLabelID uuid.UUID,
		language string,
		_type model.CodeTemplate,
		to string,
		placeholder any,
		attachments ...model.Attachment,
	) error
}

type poolSendEmailStrategy struct {
	emailSender *EmailSender
}

func newPoolSendEmailStrategy(sender *EmailSender) SendEmailStrategy {
	return &poolSendEmailStrategy{
		emailSender: sender,
	}
}

func (p *poolSendEmailStrategy) SendEmail(ctx context.Context,
	_ uuid.UUID,
	language string,
	_type model.CodeTemplate,
	to string,
	placeholder any,
	attachments ...model.Attachment,
) error {
	return p.emailSender.Send(ctx, uuid.Nil, language, _type, to, placeholder, attachments...)
}

type wlSendEmailStrategy struct {
	emailSender      *EmailSender
	wlClient         repository.Whitelabel
	checkEmailPrefix bool
}

func newWLSendEmailStrategy(
	sender *EmailSender,
	wlClient repository.Whitelabel,
	checkEmailPrefix bool,
) SendEmailStrategy {
	return &wlSendEmailStrategy{
		emailSender:      sender,
		wlClient:         wlClient,
		checkEmailPrefix: checkEmailPrefix,
	}
}

func (w *wlSendEmailStrategy) SendEmail(ctx context.Context,
	whiteLabelID uuid.UUID,
	language string,
	_type model.CodeTemplate,
	to string,
	placeholder any,
	attachments ...model.Attachment,
) error {
	emailTo := to
	wl, err := w.wlClient.GetWlByID(ctx, whiteLabelID)
	if err != nil {
		return err
	}
	if w.checkEmailPrefix && wl.Prefix != "" {
		emailTo = to[1:]
	}
	return w.emailSender.Send(ctx, whiteLabelID, language, _type, emailTo, placeholder, attachments...)
}
