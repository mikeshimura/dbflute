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
	TableDbName string
	Purpose     Purpose
	DBMetaProvider *DBMetaProvider
	SqlClause *SqlClause
	Name string
}
func (b *BaseConditionBean) AllowEmptyStringQuery(){
	(*b.SqlClause).AllowEmptyStringQuery()
}
func (b *BaseConditionBean) GetName() string{
	return b.Name
}
func (b *BaseConditionBean) AsTableDbName() string {
	return b.TableDbName
}
func (b *BaseConditionBean) GetSqlClause() *SqlClause{
	return b.SqlClause
}

