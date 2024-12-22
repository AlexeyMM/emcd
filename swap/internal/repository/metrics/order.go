package metrics

import (
	"context"
	"strconv"
	"time"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/prometheus/client_golang/prometheus"
)

type Order struct {
	order     repository.Order
	histogram *prometheus.HistogramVec
}

func NewOrder(order repository.Order, histogram *prometheus.HistogramVec) *Order {
	return &Order{
		order:     order,
		histogram: histogram,
	}
}

func (o Order) Add(ctx context.Context, order *model.Order) error {
	start := time.Now()
	err := o.order.Add(ctx, order)
	duration := time.Since(start).Seconds()

	o.histogram.WithLabelValues("order.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (o Order) Find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error) {
	start := time.Now()
	ord, err := o.order.Find(ctx, filter)
	duration := time.Since(start).Seconds()

	o.histogram.WithLabelValues("order.Find", strconv.FormatBool(err == nil)).Observe(duration)
	return ord, err
}

func (o Order) FindOne(ctx context.Context, filter *model.OrderFilter) (*model.Order, error) {
	start := time.Now()
	ord, err := o.order.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	o.histogram.WithLabelValues("order.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return ord, err
}

func (o Order) Update(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial) error {
	start := time.Now()
	err := o.order.Update(ctx, order, filter, partial)
	duration := time.Since(start).Seconds()

	o.histogram.WithLabelValues("order.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}
