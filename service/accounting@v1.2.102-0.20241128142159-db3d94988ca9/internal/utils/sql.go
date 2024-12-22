package utils

import (
	"database/sql"
	"time"
)

func Float64ToFloat64Null(f float64) sql.NullFloat64 {

	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func Int64ToInt64Null(i int64) sql.NullInt64 {

	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func Int32ToInt32Null(i int32) sql.NullInt32 {

	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}

func BoolToBoolNull(b bool) sql.NullBool {

	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func StringToStringNull(s string) sql.NullString {

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func TimeToTimeNull(t time.Time) sql.NullTime {

	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
