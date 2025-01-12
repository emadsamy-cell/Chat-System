package utils

import (
	"chat_with_go/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", config.MySQLDSN)
}
