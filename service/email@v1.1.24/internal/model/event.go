package model

const (
	ChangeWalletAddressEvent = "change_address"
)

type WLEventRequest struct {
	UserID int32
	Token  string
	Type   string
}
