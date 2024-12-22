package berrors

import "errors"

var (
	ErrProviderSettingNotFound = errors.New("whitelabel smtp settings not found")
	ErrTemplateNotFound        = errors.New("template not found")
	ErrNoMailSenderAvailable   = errors.New("no mail sender available or its settings doesn't exist for this whitelabel")
)
