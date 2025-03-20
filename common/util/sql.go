package util

import (
	"database/sql"
	"time"
)

func StringToSql(str string) (res sql.NullString) {
	if str != "" {
		res.Valid = true
		res.String = str
	}
	return
}

func SqlToString(info sql.NullString) (res string) {
	if info.Valid {
		res = info.String
	}
	return
}

func Float64ToSql(r float64) (res sql.NullFloat64) {
	if r != 0 {
		res.Valid = true
		res.Float64 = r
	}
	return
}

func SqlToFloat64(info sql.NullFloat64) (res float64) {
	if info.Valid {
		res = info.Float64
	}
	return
}

func TimeToSql(now time.Time) (res sql.NullTime) {
	res.Valid = true
	res.Time = now
	return
}
