package kafka

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"context"
	"errors"
	"github.com/segmentio/kafka-go/sasl/plain"
	"time"

	"github.com/cenkalti/backoff"
	kafkago "github.com/segmentio/kafka-go"
)

type RawWriter interface {
	WriteMessages(ctx context.Context, msgs ...Message) error
}

type Writer struct {
	writer *kafkago.Writer
	policy backoff.BackOff
}

func NewWriter(
	ctx context.Context,
	brokers []string,
	saslUser string,
	saslPassword string,
	saslEnabled bool,
	batchSize int,
) *Writer {
	funcLogError := func(format string, msg ...any) {
		sdkLog.Error(ctx, format, msg...)
	}

	if saslEnabled {
		transport := &kafkago.Transport{
			Dial:           nil,
			DialTimeout:    0,
			IdleTimeout:    0,
			MetadataTTL:    0,
			MetadataTopics: nil,
			ClientID:       "",
			TLS:            nil,
			SASL: plain.Mechanism{
				Username: saslUser,
				Password: saslPassword,
			},
			Resolver: nil,
			Context:  nil,
		}

		write := &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Topic:                  "",
			Balancer:               nil,
			MaxAttempts:            0,
			WriteBackoffMin:        0,
			WriteBackoffMax:        0,
			BatchSize:              batchSize,
			BatchBytes:             0,
			BatchTimeout:           0,
			ReadTimeout:            0,
			WriteTimeout:           0,
			RequiredAcks:           kafkago.RequireAll,
			Async:                  false,
			Completion:             nil,
			Compression:            0,
			Logger:                 nil,
			ErrorLogger:            kafkago.LoggerFunc(funcLogError),
			Transport:              transport,
			AllowAutoTopicCreation: true,
		}

		return WrapWriter(write)
	} else {
		// transport := &kafkago.Transport{
		// 	Dial: (&net.Dialer{
		// 		Timeout:   3 * time.Second,
		// 		DualStack: true,
		// 	}).DialContext,
		// 	DialTimeout:    0,
		// 	IdleTimeout:    0,
		// 	MetadataTTL:    0,
		// 	MetadataTopics: nil,
		// 	ClientID:       "",
		// 	TLS:            nil,
		// 	SASL: plain.Mechanism{
		// 		Username: "",
		// 		Password: "",
		// 	},
		// 	Resolver: nil,
		// 	Context:  nil,
		// }

		write := &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Topic:                  "",
			Balancer:               nil,
			MaxAttempts:            0,
			WriteBackoffMin:        0,
			WriteBackoffMax:        0,
			BatchSize:              batchSize,
			BatchBytes:             0,
			BatchTimeout:           0,
			ReadTimeout:            0,
			WriteTimeout:           0,
			RequiredAcks:           kafkago.RequireAll,
			Async:                  false,
			Completion:             nil,
			Compression:            0,
			Logger:                 nil,
			ErrorLogger:            kafkago.LoggerFunc(funcLogError),
			Transport:              nil,
			AllowAutoTopicCreation: true,
		}

		return WrapWriter(write)
	}
}

func WrapWriter(writer *kafkago.Writer) *Writer {
	const (
		retries  = 3                      // magic number
		duration = 250 * time.Millisecond // magic number
	)

	return &Writer{
		writer: writer,
		policy: backoff.WithMaxRetries(backoff.NewConstantBackOff(duration), retries),
	}
}

func (w *Writer) WriteMessages(ctx context.Context, msgs ...Message) error {
	if len(msgs) == 0 {
		return nil
	}
	rawMsgs := make([]kafkago.Message, 0, len(msgs))
	for i := range msgs {
		rawMsgs = append(rawMsgs, *msgs[i].toKafkaMessage())
	}
	err := backoff.RetryNotify(
		func() error {
			err := w.writer.WriteMessages(ctx, rawMsgs...)
			if err != nil && !errors.Is(err, kafkago.LeaderNotAvailable) {
				return backoff.Permanent(err)
			}

			return err
		},
		backoff.WithContext(w.policy, ctx),
		func(err error, _ time.Duration) {
			sdkLog.Error(ctx, "write message: %v", err)
		},
	)
	var e kafkago.WriteErrors
	if errors.As(err, &e) {
		return e[0]
	}

	return err
}
