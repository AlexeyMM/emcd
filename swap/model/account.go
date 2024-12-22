package model

type Account struct {
	ID      int64
	Keys    *Secrets
	IsValid bool
}

type AccountFilter struct {
	ID *int64
}

type Secrets struct {
	AccountID int64
	ApiKey    string
	ApiSecret string
}

type Accounts []*Account
