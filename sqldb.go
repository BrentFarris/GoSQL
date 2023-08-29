package gosql

import "database/sql"

type SQLDB interface {
	IsOpen() bool
	Open() error
	Close() error
	Query(Query) (*sql.Rows, error)
	Exec(Query) (sql.Result, error)
}
