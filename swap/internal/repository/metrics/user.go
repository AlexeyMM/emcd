package metrics

import (
	"context"
	"strconv"
	"time"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/prometheus/client_golang/prometheus"
)

type User struct {
	user      repository.User
	histogram *prometheus.HistogramVec
}

func NewUser(user repository.User, histogram *prometheus.HistogramVec) *User {
	return &User{
		user:      user,
		histogram: histogram,
	}
}

func (u *User) Add(ctx context.Context, user *model.User) error {
	start := time.Now()
	err := u.user.Add(ctx, user)
	duration := time.Since(start).Seconds()

	u.histogram.WithLabelValues("user.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (u *User) FindOne(ctx context.Context, filter *model.UserFilter) (*model.User, error) {
	start := time.Now()
	us, err := u.user.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	u.histogram.WithLabelValues("user.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return us, err
}
