package df

import (
	"errors"
	//"fmt"
	"reflect"
)

//Entityは必ず以下のTypeの内どれか
//int64
//float64
//bool
//string
//Time.time
//df.Date
//df.Timestamp
//Numeric
//NullBool
//NullFloat64
//NullInt64
//NullString
//NullNumeric
//NullTime
//df.NullDate
//df.NullTimestamp

type Entity interface {
	AsTableDbName() string
	HasPrimaryKeyValue() bool
	GetAsInterfaceArray() []interface{}
	GetModifiedPropertyNamesArray() []string
	IsModifiedProperty(property string) bool
	SetUp()
	GetDBMeta() *DBMeta
}

type BaseEntity struct {
	EntityModifiedProperties
}

func GetEntityValue(entity *Entity, property string) interface{} {
	cno, ok := (*(*entity).GetDBMeta()).GetColumnInfoMap()[property]
	if !ok {
		return nil
	}
	return (*entity).GetAsInterfaceArray()[cno]
}
func SetEntityValue(entity *Entity, property string, value interface{}) bool {
	var entityi interface{} = *entity
	v := reflect.ValueOf(entityi)
	m := v.MethodByName("Set" + InitCap(property))
	if !m.IsValid() {
		return false
	}
	m.Call([]reflect.Value{reflect.ValueOf(value)})
	return true
}

type EntityModifiedProperties struct {
	propertyNameMap map[string]bool
}

func (e *EntityModifiedProperties) AddPropertyName(property string) error {
	if e.propertyNameMap == nil {
		e.propertyNameMap = make(map[string]bool)
	}
	if len(property) == 0 {
		return errors.New("df005:Property length 0")
	}
	e.propertyNameMap[property] = true
	return nil
}
func (e *EntityModifiedProperties) GetModifiedPropertyNamesArray() []string {
	keys := make([]string, 0, len(e.propertyNameMap))
	for k := range e.propertyNameMap {
		keys = append(keys, k)
	}
	return keys
}
func (e *EntityModifiedProperties) IsModifiedProperty(property string) bool {
	_, ok := e.propertyNameMap[property]
	return ok
}
func (e *EntityModifiedProperties) PropertyNameMapClear() {
	e.propertyNameMap = make(map[string]bool)
}

