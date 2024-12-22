package notification

import (
	"context"
	"strconv"

	"github.com/segmentio/analytics-go/v3"

	"code.emcdtech.com/emcd/sdk/log"
)

const (
	analyticsIDParam      = "id"
	analyticsCoinParam    = "coin"
	analyticsEmailParam   = "email_notify"
	analyticsAddressParam = "address"
)

type Analytics interface {
	WithdrawalAddressChanged(coin string, emailNotify string, address string, segmentID int)
}

type analyticsNotifier struct {
	segmentClient analytics.Client
}

func NewAnalytics(key string) *analyticsNotifier {
	return &analyticsNotifier{
		segmentClient: analytics.New(key),
	}
}

func (a *analyticsNotifier) WithdrawalAddressChanged(coin string, emailNotify string, address string, segmentID int) {
	segmentIDStr := strconv.Itoa(segmentID)
	err := a.segmentClient.Enqueue(analytics.Track{
		UserId: segmentIDStr,
		Event:  "Withdrawal Address Changed",
		Properties: analytics.NewProperties().
			Set(analyticsIDParam, segmentID).
			Set(analyticsCoinParam, coin).
			Set(analyticsEmailParam, emailNotify).
			Set(analyticsAddressParam, address),
	})
	if err != nil {
		log.Error(context.Background(), "analytics: withdrawal address changed: %v", err)
	}
}
