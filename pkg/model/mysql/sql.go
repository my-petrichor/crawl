package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlDBCreate = iota
)

var databaseSQLString = []string{`
	CREATE DATABASE IF NOT EXISTS data
`}

type controller struct {
	DB *sql.DB
}

var C = new(controller)

func init() {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.0.252:3306)/mysql")
	if err != nil {
		log.Fatal(err)
	}

	C.DB = db

	databaseCreate(db)
}

func databaseCreate(db *sql.DB) {
	_, err := db.Exec(databaseSQLString[mysqlDBCreate])
	if err != nil {
		log.Fatal(err)
	}
}
