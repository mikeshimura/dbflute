package df

import (
 //"container/list"
)

type ListResultBean struct {
	List *List
	TableDbName    string
	AllRecordCount int
}

func (l *ListResultBean) New(){
	l.List=new(List)
}