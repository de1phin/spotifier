package dbmaster

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	TimeStampFormat  string = "2006-01-02 15:04:05"
)

type Database interface {

	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) (*sql.Row)

}

func DBConnect(dbname, dblogin, dbpass string) (db *sql.DB, err error) {
	authString := fmt.Sprintf("%s:%s@/%s", dblogin, dbpass, dbname)
	db, err = sql.Open("mysql", authString)
	if err != nil {
		return nil, err
	}

	return db, nil
}