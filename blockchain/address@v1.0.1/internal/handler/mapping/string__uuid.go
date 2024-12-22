package mapping

import (
	"fmt"

	"github.com/google/uuid"
)

func MapStringToUuid(uuidString string) (*uuid.UUID, error) {
	if uuidStringParsed, err := uuid.Parse(uuidString); err != nil {

		return nil, fmt.Errorf("failed parse string to uuid: %s, %w", uuidString, err)
	} else {

		return &uuidStringParsed, nil
	}
}
