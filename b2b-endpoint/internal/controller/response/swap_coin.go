package response

type Coin struct {
	Title    string     `json:"title"`
	IconUrl  string     `json:"icon_url"`
	Networks []*Network `json:"networks"`
}

type Network struct {
	Title             string `json:"title"`
	Accuracy          int    `json:"accuracy"`
	WithdrawSupported bool   `json:"withdraw_supported"`
}
