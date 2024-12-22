package kafka

import (
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func GetKafka[T any, PT interface {
	proto.Message
	*T
}]() (*accountingPb.Kafka, error) {
	dto := new(T)
	if t, ok := (any(dto)).(proto.Message); !ok {

		return nil, fmt.Errorf("cannot cast %T to %T", dto, any(dto))
	} else {
		eKafkaDesc := accountingPb.E_Kafka.TypeDescriptor()
		msg := t.ProtoReflect().Descriptor().Options().ProtoReflect()
		if !msg.Has(eKafkaDesc) {

			return nil, fmt.Errorf("no option field for type %T", dto)
		} else if kafka, ok := accountingPb.E_Kafka.InterfaceOf(msg.Get(eKafkaDesc)).(*accountingPb.Kafka); !ok {

			return nil, fmt.Errorf("cannot cast %T to %T", msg.Get(eKafkaDesc), accountingPb.E_Kafka)
		} else {

			return kafka, nil

		}
	}
}
