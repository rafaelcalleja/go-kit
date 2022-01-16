package adapters

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func newMySQLConnection() (*sqlx.DB, error) {
	config := mysql.NewConfig()

	config.Net = "tcp"
	config.Addr = os.Getenv("MYSQL_ADDR")
	config.User = "user"
	config.Passwd = "password"
	config.DBName = os.Getenv("MYSQL_DATABASE")
	config.ParseTime = true // with that parameter, we can use time.Time in mysqlHour.Hour

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to MySQL")
	}

	return db, nil
}

func TestNewTransactionInitializerDb(t *testing.T) {
	connection, _ := newMySQLConnection()

	s := transaction.NewTransactionalSession(
		NewTransactionInitializerDb(connection),
	)

	called := false

	_ = s.ExecuteAtomically(func() error {
		called = true

		return nil
	})

	assert.True(t, called)
}
