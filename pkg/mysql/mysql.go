package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var GlobalMySQLDB *sql.DB

func init() {
	db, err := NewMySQLDB("dev_main:kr4e65a2xJ9@tcp(localhost:3312)/dev_main")

	if err != nil {
		panic(err)
	}

	GlobalMySQLDB = db
}

func NewMySQLDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")

	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}

