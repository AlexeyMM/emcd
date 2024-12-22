package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	sdkLog "code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitService interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
}

type rabbitServiceImpl struct {
	exchangeName string
	rabbitRepo   repository.RabbitRepository
	mu           *sync.Mutex
}

func NewRabbitService(exchangeName string, rabbitRepo repository.RabbitRepository) RabbitService {

	return &rabbitServiceImpl{
		exchangeName: exchangeName,
		rabbitRepo:   rabbitRepo,
		mu:           new(sync.Mutex),
	}
}

func (r *rabbitServiceImpl) Publish(ctx context.Context, routingKey string, body []byte) error {
	if err := r.pubMsg(ctx, r.exchangeName, routingKey, body); err != nil {
		if err := r.reconnect(); err != nil {

			sdkLog.Panic(ctx, err.Error())
		} else if err := r.pubMsg(ctx, r.exchangeName, routingKey, body); err != nil {

			sdkLog.Panic(ctx, err.Error())
		}
	}

	return nil
}

func (r *rabbitServiceImpl) pubMsg(ctx context.Context, exchange, routingKey string, body []byte) error {

	return r.rabbitRepo.GetChannel().PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         nil,
			ContentType:     "application/json",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            body,
		},
	)
}

func (r *rabbitServiceImpl) reconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.rabbitRepo.Reconnect(); err != nil {

		return fmt.Errorf("reconnect: %w", err)
	} else {

		return nil
	}
}
