package helper

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover() // overlaping
	if err != nil {
		errorRollback := tx.Rollback()
		PanicError(errorRollback)
		panic(err)
	} else {
		errorCommit := tx.Commit()
		PanicError(errorCommit)
	}
}

func SqlNullFromString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func SqlNullFromInt(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func SqlNullFromIntPointer(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
	return sql.NullInt64{
		Int64: *i,
		Valid: true,
	}
}

func SqlNullFromFloat(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func SqlNullFromTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func UuidNullFromUuid(u uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		UUID:  u,
		Valid: true,
	}
}

func SqlNullFromBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}
