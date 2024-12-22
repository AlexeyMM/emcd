package utils

import (
	"github.com/google/uuid"
)

func UuidToUuidNull(u uuid.UUID) uuid.NullUUID {

	return uuid.NullUUID{
		UUID:  u,
		Valid: true,
	}
}
