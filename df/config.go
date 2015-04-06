package df

import (
	"database/sql"
	"github.com/mikeshimura/dbflute/log"
)
var Bci *BehaviorCommandInvoker 
var DBCurrent_I *DBCurrent

var LogStop bool

func DFLog(logstr string) {
	if !LogStop {
		log.DebugConv("dfLog", logstr)
	}
}
func DFErrorLog(logstr string) {
	log.ErrorConv("dfLog", logstr)
}

func DbBegin(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		DFErrorLog("Begin Transaction Error:" + err.Error())
		return tx, err
	}
	DFLog("Begin Transaction")
	return tx, err
}
func TxCommit(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		DFErrorLog("Transaction Commit Error:" + err.Error())
		return err
	}
	DFLog("Transaction Commit")
	return err
}
func TxRollback(tx *sql.Tx) error {
	err := tx.Rollback()
	if err != nil {
		DFErrorLog("Transaction Rollback Error:" + err.Error())
		return err
	}
	DFLog("Transaction Rollback")
	return err
}
