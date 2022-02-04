package transaction

import (
	"database/sql"
)

type TxRepository interface {
	WithTrx(tx *sql.Tx) interface{}
}
