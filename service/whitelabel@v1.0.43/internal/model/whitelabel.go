package model

import "github.com/google/uuid"

type WhiteLabel struct {
	ID                    uuid.UUID
	UserID                int32   `db:"user_id"`
	SegmentID             int32   `db:"segment_id"`
	Origin                string  `db:"origin"`
	Prefix                string  `db:"prefix"`
	SenderEmail           string  `db:"sender_email"`
	Domain                string  `db:"domain"`
	APIKey                string  `db:"api_key"`
	URL                   string  `db:"url"`
	Version               int     `db:"version"`
	MasterSlave           bool    `db:"master_slave"`
	MasterFee             float64 `db:"master_fee"`
	IsTwoFAEnabled        bool    `db:"is_two_fa_enabled"`
	IsCaptchaEnabled      bool    `db:"is_captcha_enabled"`
	IsEmailConfirmEnabled bool    `db:"is_email_confirm_enabled"`
}

type WlRole struct {
	UserID       uuid.UUID
	WhiteLabelID uuid.UUID
	Role         string
}

type (
	WlConfig struct {
		RefID               string    `json:"ref_id"`
		MediaID             string    `json:"media_id"`
		Origin              string    `json:"origin"`
		Title               string    `json:"title"`
		Commission          float64   `json:"commission"`
		Logo                string    `json:"logo"`
		Favicon             string    `json:"favicon"`
		StratumLists        []Stratum `json:"stratum_list"`
		Colors              Colors    `json:"colors"`
		ColorsJB            []byte    `json:"-"`
		FirmwareInstruction string
		Lang                string   `json:"lang"`
		PossibleLang        []string `json:"possible_lang"`
		WhitelabelID        string   `json:"whitelabel_id"`
	}

	Colors struct {
		Emcd            string `json:"emcd"`
		EmcdHover       string `json:"emcd-hover"`
		EmcdDisabled    string `json:"emcd-disabled"`
		EmcdProgressBar string `json:"emcd-progress-bar"`
		Error           string `json:"error"`
		Error2          string `json:"error-2"`
		Success         string `json:"success"`
		Attention       string `json:"attention"`
		Bg1             string `json:"bg-1"`
		Bg2             string `json:"bg-2"`
		Bg3             string `json:"bg-3"`
		Bg4             string `json:"bg-4"`
		Text1           string `json:"text-1"`
		Text2           string `json:"text-2"`
		Text3           string `json:"text-3"`
		Text4           string `json:"text-4"`
		Text5           string `json:"text-5"`
		AttentionStatus string `json:"attention-status"`
		SuccessStatus   string `json:"success-status"`
		NegativeStatus  string `json:"negative-status"`
		ProcessStatus   string `json:"process-status"`
		NeutralStatus   string `json:"neutral-status"`
	}

	Stratum struct {
		RefID  string `json:"ref_id"`
		Coin   string `json:"coin"`
		Region string `json:"region"`
		Number string `json:"number"`
		Url    string `json:"url"`
	}

	AllowOrigin struct {
		UserID int32  `json:"user_id" db:"user_id"`
		Origin string `json:"origin" db:"origin"`
	}
)

type WLCoins struct {
	WlID   uuid.UUID `json:"wl_id" db:"wl_id"`
	CoinID string    `json:"coin_id" db:"coin_id"`
}
