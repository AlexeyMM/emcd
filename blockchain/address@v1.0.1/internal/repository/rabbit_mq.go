package repository

import (
	"context"
	"fmt"

	sdkLog "code.emcdtech.com/emcd/sdk/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitRepository interface {
	Connect() error
	Reconnect() error
	Close() error

	GetConnection() *amqp.Connection
	GetChannel() *amqp.Channel
}

type rabbitRepository struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitRepository(url string) RabbitRepository {
	return &rabbitRepository{
		url:  url,
		conn: nil,
		ch:   nil,
	}
}

func (r *rabbitRepository) Connect() error {
	if rabbitConn, err := amqp.Dial(r.url); err != nil {

		return fmt.Errorf("failed dial to rabbit mq: %s, %w", r.url, err)
	} else {
		r.conn = rabbitConn

		if rabbitChan, err := rabbitConn.Channel(); err != nil {
			return fmt.Errorf("failed open channel of rabbit mq: %s, %w", r.url, err)

		} else {
			r.ch = rabbitChan

			return nil
		}
	}
}

func (r *rabbitRepository) Reconnect() error {
	if err := r.Close(); err != nil {

		return fmt.Errorf("re close: %w", err)
	} else if err = r.Connect(); err != nil {

		return fmt.Errorf("re connect: %w", err)
	} else {

		return nil
	}
}

func (r *rabbitRepository) Close() error {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {

			return fmt.Errorf("rabbit conn close: %w", err)
		} else {
			sdkLog.Info(context.Background(), "rabbit conn close")

		}
	}

	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			return fmt.Errorf("chan close: %w", err)

		} else {
			sdkLog.Info(context.Background(), "rabbit chan close")

		}
	}

	return nil
}

func (r *rabbitRepository) GetConnection() *amqp.Connection {

	return r.conn
}

func (r *rabbitRepository) GetChannel() *amqp.Channel {

	return r.ch
}
