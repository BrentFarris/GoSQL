package gosql

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	path string
	conn *sql.DB
}

func NewSQLite(path string) SQLite {
	return SQLite{
		path: path,
	}
}

func (db SQLite) IsOpen() bool {
	return db.conn != nil
}

func (db *SQLite) Open() error {
	if db.IsOpen() {
		return nil
	}
	conn, err := sql.Open("sqlite3", db.path)
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *SQLite) Close() error {
	err := db.conn.Close()
	db.conn = nil
	return err
}

func (db SQLite) Query(query Query) (*sql.Rows, error) {
	return db.conn.Query(query.Build())
}

func (db SQLite) Exec(query Query) (sql.Result, error) {
	return db.conn.Exec(query.Build())
}
