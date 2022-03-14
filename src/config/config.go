package config

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

func GetDB() (db *sql.DB, errorMessage error) {
	db, errorMessage = sql.Open("mssql", "server=TRUONGNBSE62373\\SQLEXPRESS;user id=sa;password=123456789;database=Assignment")
	return
}
func CloseDB(db *sql.DB) {
	db.Close()
}
