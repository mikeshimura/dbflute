package df

import (
	"database/sql"
	"reflect"
	"strings"
)

type ResultSetHandler struct {
	Enitities          List
	Infos              List
	BaseTableDBName    string
	BasePointAliasName string
	outerJoinMap       map[string]*LeftOuterJoinInfo
	aliasMap           map[string]interface{}
	simpleQuery        bool
}

func (p *ResultSetHandler) GetListResultBean(rows *sql.Rows, entity string,
	sqlClause interface{}) *ListResultBean {
			list := new(ListResultBean)
	list.New()
	p.initialize(entity, sqlClause)

	for rows.Next() {
		array := p.getColumnArray()
		err := rows.Scan(array...)
		if err != nil {
			panic(err.Error())
		}
		list.List.Add((p.Enitities.Get(0)).(*ResultEntity).Entity)
		if p.simpleQuery == false {
			p.setParent()
		}
	}
	list.AllRecordCount = list.List.Size()
	tmap := (*DBMetaProvider_I).TableDbNameInstanceMap[entity]
	if tmap == nil {
		list.TableDbName = "__" + entity
	} else {
		list.TableDbName = (*tmap).GetTableDbName()
	}

	return list
}

func (p *ResultSetHandler) setParent() {
	for i, ent := range p.Enitities.data {
		if i == 0 {
			continue
		}
		entity := ent.(*ResultEntity)
		var parentResultEntity *ResultEntity
		for ii := i - 1; i >= 0; i-- {
			temp := p.Enitities.Get(ii)
			tempRE := temp.(*ResultEntity)
			if tempRE.TableDbName == entity.Parent {
				parentResultEntity = tempRE
				break
			}
		}
		v := reflect.ValueOf(parentResultEntity.Entity)
		m := v.MethodByName("Set" + InitCap(entity.TableDbName) + "_R")
		m.Call([]reflect.Value{reflect.ValueOf(entity.Entity)})
	}
}
func (p *ResultSetHandler) getColumnArray() []interface{} {
	size := p.Infos.Size()
	ints := make([]interface{}, size)
	var arrayBu []interface{}
	for _, ent := range p.Enitities.data {
		entity := ent.(*ResultEntity)
		table, array := BhvUtil_I.GetEntityAndInterfaceArray(entity.TableDbName)
		entity.Entity = table
		entity.Columns = array
		arrayBu = array
	}
	if p.simpleQuery {
		return arrayBu
	}
	for i, c := range p.Infos.data {
		cinfo := c.(*ResultSetInfo)
		ent := p.Enitities.Get(cinfo.EntityNo)
		resultEntity := ent.(*ResultEntity)
		ints[i] = resultEntity.Columns[cinfo.ColumnNo]
	}
	return ints
}
func (p *ResultSetHandler) initialize(entity string, sqlClause interface{}) {
	dbmeta := DBMetaProvider_I.TableDbNameInstanceMap[entity]
	if dbmeta == nil {
		panic("System Error dbmeta Not Found:" + entity)
	}
	p.BaseTableDBName = entity
	if sqlClause == nil || entity=="D_Int64"{
		p.simpleQuery = true
		resultEntityBase := new(ResultEntity)
		resultEntityBase.TableDbName = p.BaseTableDBName
		p.addColumnList(resultEntityBase, dbmeta)
		p.Enitities.Add(resultEntityBase)
		return
	}
	bsc := reflect.ValueOf(sqlClause).Elem().FieldByName("BaseSqlClause")
	baseSqlClause := (bsc.Interface()).(BaseSqlClause)
	if entity != baseSqlClause.TableDbName {
		panic("System Errir entity != baseSqlClause.TableDbName")
	}

	p.BasePointAliasName = baseSqlClause.BasePorintAliasName
	p.outerJoinMap = baseSqlClause.outerJoinMap
	resultEntityBase := new(ResultEntity)
	resultEntityBase.TableDbName = p.BaseTableDBName
	resultEntityBase.AliasName = p.BasePointAliasName
	//fmt.Printf("entity %s \n", entity)
	p.addColumnList(resultEntityBase, dbmeta)
	p.aliasMap = make(map[string]interface{})
	p.aliasMap[resultEntityBase.AliasName] = p.Enitities.Size()
	p.Enitities.Add(resultEntityBase)

	for _, realColumnName := range baseSqlClause.selectClauseInfo.data {
		p.setInfos(realColumnName)
	}
}
func (p *ResultSetHandler) addColumnList(resultEntity *ResultEntity, dbmeta *DBMeta) {
	resultEntity.ColumnList = make(map[string]int)
	for i, ci := range (*dbmeta).GetColumnInfoList().data {
		columnInfo := ci.(*ColumnInfo)
		resultEntity.ColumnList[columnInfo.ColumnSqlName.ColumnSqlName] = i
	}
}
func (p *ResultSetHandler) setInfos(realColumnName string) {
	column := strings.Split(realColumnName, ".")
	if len(column) != 2 {
		panic("System Error realColumnName:" + realColumnName)
	}
	alias := column[0]
	name := column[1]
	entityno := p.aliasMap[alias]
	if entityno == nil {
		p.addEntity(alias)
		entityno = p.aliasMap[alias]
	}
	if entityno == nil {
		panic("System Error Alias Not Found:" + alias)
	}
	entityNo := entityno.(int)
	entity := (p.Enitities.Get(entityNo)).(*ResultEntity)
	si := new(ResultSetInfo)
	si.EntityNo = entityNo
	si.ColumnNo = entity.ColumnList[name]
	p.Infos.Add(si)
}
func (p *ResultSetHandler) addEntity(alias string) {
	joinInfo := p.outerJoinMap[alias]
	resultEntity := new(ResultEntity)
	resultEntity.TableDbName = joinInfo.foreignTableDbName
	resultEntity.AliasName = joinInfo.foreignAliasName
	resultEntity.Parent = joinInfo.localTableDbName
	dbmeta := DBMetaProvider_I.TableDbNameInstanceMap[joinInfo.foreignTableDbName]
	p.addColumnList(resultEntity, dbmeta)
	p.aliasMap[resultEntity.AliasName] = p.Enitities.Size()
	p.Enitities.Add(resultEntity)
}

type ResultEntity struct {
	Entity      interface{}
	Columns     []interface{}
	TableDbName string
	Parent      string
	AliasName   string
	ColumnList  map[string]int
}
type ResultSetInfo struct {
	EntityNo int
	ColumnNo int
}
