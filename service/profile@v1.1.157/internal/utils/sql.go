package utils

import (
	"database/sql"
)

func StringToStringNull(s string) sql.NullString {

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
