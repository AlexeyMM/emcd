package metrics

import (
	"context"
	"strconv"
	"time"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/prometheus/client_golang/prometheus"
)

type Account struct {
	account   repository.Account
	histogram *prometheus.HistogramVec
}

func NewAccount(account repository.Account, histogram *prometheus.HistogramVec) *Account {
	return &Account{
		account:   account,
		histogram: histogram,
	}
}

func (a *Account) Add(ctx context.Context, account *model.Account) error {
	start := time.Now()
	err := a.account.Add(ctx, account)
	duration := time.Since(start).Seconds()

	a.histogram.WithLabelValues("account.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (a *Account) Find(ctx context.Context, filter *model.AccountFilter) (model.Accounts, error) {
	start := time.Now()
	acc, err := a.account.Find(ctx, filter)
	duration := time.Since(start).Seconds()

	a.histogram.WithLabelValues("account.Find", strconv.FormatBool(err == nil)).Observe(duration)
	return acc, err
}

func (a *Account) FindOne(ctx context.Context, filter *model.AccountFilter) (*model.Account, error) {
	start := time.Now()
	acc, err := a.account.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	a.histogram.WithLabelValues("account.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return acc, err
}
