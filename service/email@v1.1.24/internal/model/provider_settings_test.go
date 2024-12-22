package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

func TestGetSetting(t *testing.T) {
	p := model.Provider{
		Name:       model.SMTPProviderName,
		RawSetting: []byte("{}"),
	}
	_, err := model.GetSetting[model.SmtpSetting](p)
	require.NoError(t, err)
	_, err = model.GetSetting[model.MailgunSetting](p)
	require.Error(t, err)
}
