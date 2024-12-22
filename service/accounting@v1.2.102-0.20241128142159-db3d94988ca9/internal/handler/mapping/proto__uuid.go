package mapping

import (
	"fmt"
	"github.com/google/uuid"
)

func MapProtoUuidToUuid(uuidStr string) (uuid.UUID, error) {
	if uuidParsed, err := uuid.Parse(uuidStr); err != nil {

		return uuid.UUID{}, fmt.Errorf("failed parse uuid: %s, %w", uuidStr, err)
	} else {

		return uuidParsed, nil
	}
}
