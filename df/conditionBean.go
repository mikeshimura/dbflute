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
//"github.com/mikeshimura/dbflute/log"
)

type Purpose int

const (
	P_Normal_Use = iota
	P_SCALAR_SELECT
	P_OR_SCOPE_QUERY
)

type ConditionBean interface {
	//AsTableDbName() string
	//	GetName() string
	GetBaseConditionBean() *BaseConditionBean
}
type BaseConditionBean struct {
	TableDbName    string
	Purpose        Purpose
	DBMetaProvider *DBMetaProvider
	SqlClause      *SqlClause
	Name           string
}

func (b *BaseConditionBean) AllowEmptyStringQuery() {
	(*b.SqlClause).AllowEmptyStringQuery()
}
func (b *BaseConditionBean) GetName() string {
	return b.Name
}
func (b *BaseConditionBean) AsTableDbName() string {
	return b.TableDbName
}
func (b *BaseConditionBean) GetSqlClause() *SqlClause {
	return b.SqlClause
}
func (b *BaseConditionBean) SelectHint() string{
	return (*b.GetSqlClause()).CreateSelectHint()
}
func (b *BaseConditionBean) DoSetupSelect(queryLocal *BaseConditionQuery,
	queryRemote *BaseConditionQuery) {
	foreignPropertyName := queryRemote.ForeignPropertyName
	foreignTableAliasName := queryRemote.AliasName
	localRelationPath := queryLocal.RelationPath
	foreignRelationPath := queryRemote.RelationPath
	localPropertyName:=queryLocal.TableDbName
	(*b.SqlClause).RegisterSelectedRelation(foreignTableAliasName, localPropertyName, foreignPropertyName,
		localRelationPath, foreignRelationPath,queryRemote.TableDbName)
}
