package common

import "github.com/mattn/go-sqlite3"

func IsSqliteErr(err error) bool {
	if _, ok := err.(sqlite3.Error); ok {
		return true
	}
	return false
}

func IsSqliteErrConstraint(err error) bool {
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return true
		}
	}
	return false
}
