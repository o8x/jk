package puresqlite

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func Get() *sql.DB {
	return db
}

func Default() *sql.DB {
	_ = Init("file:defmemsqlite?mode=memory")
	return db
}

func Init(file string) error {
	ins, err := sql.Open("sqlite", file)
	if err != nil {
		return err
	}

	db = ins
	return err
}
