package kafka

import (
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"time"
)

type Message struct {
	ID      uuid.UUID
	Topic   string
	Key     []byte
	Headers map[string][]byte
	Value   []byte
}

func (m Message) toKafkaMessage() *kafka.Message {
	getHeader := func(headers map[string][]byte) []kafka.Header {
		r := make([]kafka.Header, 0, len(headers))
		for k, v := range headers {
			r = append(r, kafka.Header{Key: k, Value: v})
		}
		return r
	}

	return &kafka.Message{
		Topic:         m.Topic,
		Partition:     0,
		Offset:        0,
		HighWaterMark: 0,
		Key:           m.Key,
		Value:         m.Value,
		Headers:       getHeader(m.Headers),
		WriterData:    nil,
		Time:          time.Time{},
	}

}
