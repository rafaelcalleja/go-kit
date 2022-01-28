package mysql_tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	db, _ := NewMySQLConnection()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(50000 * time.Microsecond)

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			tx, _ := db.Begin()
			mSql := "select * from db.products"
			rows, _ := tx.Query(mSql)
			err := rows.Close() // here, if you do not release the connection to the pool, other concurrencies will block after five runs
			require.NoError(t, err)
			err = tx.Rollback()
			require.NoError(t, err)
			wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}
