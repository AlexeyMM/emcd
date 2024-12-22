package model

type Coin struct {
	ID string

	Title       string
	Description string
	MediaURL    string
	IsActive    bool

	Networks []*Network
}

type Network struct {
	ID     string
	CoinID string
	Title  string
}
