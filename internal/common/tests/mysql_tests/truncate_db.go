package mysql_tests

import (
	"database/sql"
	"fmt"
)

func OpenDB(driver, address string, maxIdleConns int) *sql.DB {
	db, _ := sql.Open(driver, address)
	db.SetMaxIdleConns(maxIdleConns)

	return db
}

func CloseDB(db *sql.DB) {
	_ = db.Close()
}

func TruncateTables(db *sql.DB, tables []string) {
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=0;")

	for _, v := range tables {
		_, _ = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", v))
	}

	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=1;")
}
