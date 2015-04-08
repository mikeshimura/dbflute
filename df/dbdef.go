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

