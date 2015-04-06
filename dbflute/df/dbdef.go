package df

import (
	"errors"
	"github.com/mikeshimura/dbflute/log"
)

var locked = true

type DBDef interface {
	Code() string
	CodeAlias() string
}

type PostgreSQL struct {

}

func (w *PostgreSQL) Code() string {
	return "postgresql"
}
func (w *PostgreSQL) CodeAlias() string {
	return "postgre"
}
type MySQL struct {

}

func (w *MySQL) Code() string {
	return "mysql"
}
func (w *MySQL) CodeAlias() string {
	return ""
}
type SQLServer struct {

}

func (w *SQLServer) Code() string {
	return "sqlserver"
}
func (w *SQLServer) CodeAlias() string {
	return "mssql"
}

func Lock() {
	if locked {
		return
	}
	if log.IsEnabled() {
		log.Info("df.DBDef", "...Locking the singleton world of the DB definition!")
	}
	locked = true
}

func UnLock() {
	if locked == false {
		return
	}
	if log.IsEnabled() {
		log.Info("df.DBDef", "...Unlocking the singleton world of the DB definition!")
	}
	locked = false
}
func assertUnlocked() error {
	if !locked {
		return nil
	}
	return errors.New("df004:The DB definition is locked.")
}

