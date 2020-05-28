package utils

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("mysql", "root:1234w5asd@tcp(localhost:3306)/chatsystem")
	if err != nil {
		panic(err.Error())
	}
}