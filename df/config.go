/*
 * Copyright 2014-2015 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 */
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

func TxBegin(db *sql.DB) (*sql.Tx, error) {
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
