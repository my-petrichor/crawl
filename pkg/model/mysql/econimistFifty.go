package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlEconomistFiftyTableCreate = iota
	mysqlEconomistFiftyInsert
	mysqlEconomistFiftyGet
)

var economistFiftySQLString = []string{
	`CREATE TABLE IF NOT EXISTS data.economistFifty(
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(20) NOT NULL,
		identity VARCHAR(100) NOT NULL,
		academy VARCHAR(32) NOT NULL
		)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	`INSERT INTO data.economistFifty (name, identity, academy) VALUES(?, ?, ?)`,
	`SELECT name FROM data.economistFifty WHERE name = ?`,
}

func EconomistFiftyTableCreate(db *sql.DB) {
	_, err := db.Exec(economistFiftySQLString[mysqlEconomistFiftyTableCreate])
	if err != nil {
		log.Printf("economistFifty create error: %v\n", err)
	}
}

func EconomistFiftyInsert(db *sql.DB, data []string) {
	_, err := db.Exec(economistFiftySQLString[mysqlEconomistFiftyInsert], data[0], data[1], data[2])
	if err != nil {
		log.Printf("economistFifty insert error: %v\n", err)
	}
}

func EconomistFiftyGet(db *sql.DB, data []string) (name string) {
	err := db.QueryRow(economistFiftySQLString[mysqlEconomistFiftyGet], data[0]).Scan(&name)
	if err != nil {
		return ""
	}

	return name
}
