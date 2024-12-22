package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProviderName string

const (
	SMTPProviderName    ProviderName = "smtp"
	MailgunProviderName ProviderName = "mailgun"
)

func (p ProviderName) String() string {
	return string(p)
}

// Setting настройка по отправке email для white label.
type Setting struct {
	WhiteLabelID uuid.UUID `db:"white_label_id"`
	Providers    Providers `db:"providers"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Provider сдержит настройки провайдера email сервиса.
type Provider struct {
	Name       ProviderName    `json:"name"`
	RawSetting json.RawMessage `json:"setting"`
}

// SmtpSetting конфигурация SMTP сервера.
type SmtpSetting struct {
	markerProviderSetting
	Username               string `json:"username"`
	Password               string `json:"password"`
	ServerAddress          string `json:"server_address"`
	ServerPort             int    `json:"server_port"`
	FromAddress            string `json:"from_address"`
	FromAddressDisplayedAs string `json:"from_address_displayed_as"`
}

func (s SmtpSetting) ProviderName() ProviderName {
	return SMTPProviderName
}

// MailgunSetting конфигурация Mailgun сервера.
type MailgunSetting struct {
	markerProviderSetting
	Domain                 string `json:"domain"`
	ApiKey                 string `json:"api_key"`
	FromAddress            string `json:"from_address"`
	ApiBase                string `json:"api_base"`
	FromAddressDisplayedAs string `json:"from_address_displayed_as"`
}

func (s MailgunSetting) ProviderName() ProviderName {
	return MailgunProviderName
}

type Providers []Provider

// Scan implements the Scanner interface.
func (x *Providers) Scan(value interface{}) (err error) {
	var b []byte
	switch v := value.(type) {
	case nil:
		b = nil
	case string:
		b = ([]byte)(v)
	case []byte:
		b = v
	}
	// по идее, данное условие не должно выполнятся, т.к. из БД должен всегда считываться
	if len(b) == 0 || strings.ToLower(string(b)) == "null" {
		*x = nil
		return nil
	}
	return json.Unmarshal(b, x)
}

// Value implements the driver Valuer interface.
func (x Providers) Value() (driver.Value, error) {
	if len(x) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(x)
}

// ProviderSetting маркер, что тип является dto описывающий настройки для email провайдера
type ProviderSetting interface {
	mustEmbedUnimplemented()
	ProviderName() ProviderName
}

type markerProviderSetting struct{}

func (markerProviderSetting) mustEmbedUnimplemented() {}

// GetSetting получение DTO настроек провайдера email, выдаст ошибку в случае если в ProviderSetting один,
// а запрошена DTO от другого.
func GetSetting[T any, PT interface {
	ProviderSetting
	*T
}](p Provider) (T, error) {
	r := new(T)
	if p.Name != (interface{}(r).(ProviderSetting)).ProviderName() {
		return *interface{}(r).(*T), fmt.Errorf("invalid type %T", *interface{}(r).(*T))
	}
	err := json.Unmarshal(p.RawSetting, r)
	return *interface{}(r).(*T), err
}

func GetProvider[T any, PT interface {
	ProviderSetting
	*T
}](v T) (Provider, error) {
	p := Provider{
		Name: (interface{}(v).(ProviderSetting)).ProviderName(),
	}
	b, err := json.Marshal(v)
	if err != nil {
		return p, fmt.Errorf("marshal %T: %w", v, err)
	}
	p.RawSetting = b
	return p, nil
}
