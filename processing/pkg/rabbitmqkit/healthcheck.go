package rabbitmqkit

import (
	"context"

	"code.emcdtech.com/emcd/sdk/app"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ChannelHealthChecker struct {
	ch   *amqp.Channel
	name string
}

func NewChannelHealthChecker(ch *amqp.Channel, name string) *ChannelHealthChecker {
	return &ChannelHealthChecker{ch: ch, name: name}
}

func (c *ChannelHealthChecker) Check(ctx context.Context) app.HealthCheckServingStatus {
	if c.ch.IsClosed() {
		return app.NotServingHealthCheckServingStatus
	}

	return app.ServingHealthCheckServingStatus
}

func (c *ChannelHealthChecker) ServiceName() string {
	return c.name
}
