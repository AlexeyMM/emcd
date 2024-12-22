package coinwatch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"code.emcdtech.com/emcd/sdk/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"

	oteltrace "go.opentelemetry.io/otel/trace"

	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/rabbitmqkit"
)

const (
	depositRoutingKey           = "deposit" // new txs, not confirmed
	topUpRoutingKey             = "topup"   // confirmed txs, user account balance topped up
	processingQueueForCoinwatch = "processing-queue-for-coinwatch"
)

type Consumer struct {
	rmqCh              *amqp.Channel
	exchangeName       string
	transactionService service.Transaction
}

func NewConsumer(
	rmqCh *amqp.Channel,
	exchangeName string,
	transactionService service.Transaction,
) *Consumer {
	return &Consumer{
		rmqCh:              rmqCh,
		exchangeName:       exchangeName,
		transactionService: transactionService,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	// Declare the queue
	q, err := c.rmqCh.QueueDeclare(
		processingQueueForCoinwatch, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue (queue=%s): %w", processingQueueForCoinwatch, err)
	}

	// Bind queue to exchange with deposit routing key
	err = c.rmqCh.QueueBind(
		q.Name,            // queue name
		depositRoutingKey, // routing key
		c.exchangeName,    // exchange
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange with deposit routing key: %w", err)
	}

	// Bind queue to exchange with topup routing key
	err = c.rmqCh.QueueBind(
		q.Name,          // queue name
		topUpRoutingKey, // routing key
		c.exchangeName,  // exchange
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange with topup routing key: %w", err)
	}

	msgs, err := c.rmqCh.ConsumeWithContext(
		ctx,
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer (queue=%s): %w", q.Name, err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return ctx.Err() // channel was probably closed because ctx got cancelled
			}

			if err := c.processMessage(ctx, msg); err != nil {
				// skip message if it has invalid body, no point in retrying
				if errors.Is(err, errInvalidJSON) {
					log.SError(ctx, "message with invalid json received", map[string]any{
						"error": err,
						"msg":   string(msg.Body),
					})

					if err := msg.Ack(false); err != nil {
						return fmt.Errorf("ack invalid json mesg: %w", err)
					}

					continue
				}

				return fmt.Errorf("processMessage: %w", err)
			}

			if err := msg.Ack(false); err != nil {
				return fmt.Errorf("ack: %w", err)
			}
		}
	}
}

var errInvalidJSON = errors.New("invalid json")

// coinwatchTransactionMessage contains info about a transaction that has occurred on blockchain.
// We expect 2 of such messages for each transaction:
// 1) when transaction is detected but is not confirmed yet
// 2) when it gets confirmed by blockchain.
type coinwatchTransactionMessage struct {
	Address     string          `json:"address"` // receiver address, must be the deposit address for some invoice
	Amount      decimal.Decimal `json:"amount"`
	TxHash      string          `json:"tx_hash"`
	CoinCode    string          `json:"coin_code"` // coin id
	IsConfirmed bool            `json:"is_confirmed"`
}

func (c *Consumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, rabbitmqkit.RabbitmqHeadersCarrier(msg.Headers))
	opts := []oteltrace.SpanStartOption{
		oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
		oteltrace.WithAttributes(semconv.MessagingSystemRabbitmq),
		oteltrace.WithAttributes(semconv.MessagingRabbitmqDestinationRoutingKeyKey.String(msg.RoutingKey)),
	}

	ctx, span := otel.Tracer("coinwatch-consumer").Start(ctx, msg.RoutingKey, opts...)
	defer span.End()

	var coinwatchMsg coinwatchTransactionMessage
	if err := json.Unmarshal(msg.Body, &coinwatchMsg); err != nil {
		return fmt.Errorf("unmarshal message: %w: %w", errInvalidJSON, err)
	}

	tx := &model.Transaction{
		Address:     coinwatchMsg.Address,
		Amount:      coinwatchMsg.Amount,
		Hash:        coinwatchMsg.TxHash,
		CoinID:      coinwatchMsg.CoinCode,
		IsConfirmed: coinwatchMsg.IsConfirmed,
	}

	if err := c.transactionService.ProcessTransaction(ctx, tx); err != nil {
		return fmt.Errorf("processTransaction: %w", err)
	}

	return nil
}
