package accessory

import (
"errors"
	//"fmt"
)

type EntityModifiedProperties struct {
}

var propertyNameMap = make(map[string]bool)

func (e *EntityModifiedProperties) AddPropertyName(property string)error{
	if len(property)==0{
		return errors.New("df005:Property length 0")
	}
	propertyNameMap[property]=true
	return nil
}
func (e *EntityModifiedProperties) GetPropertyNamesArray() []string {
    keys := make([]string, 0, len(propertyNameMap))
    for k := range propertyNameMap {
        keys = append(keys, k)
    }
    return keys
}
func (e *EntityModifiedProperties) IsModifiedProperty(property string)bool{
	_, ok:=propertyNameMap[property]
	return ok
}
func (e *EntityModifiedProperties) PropertyNameMapClear(){
	propertyNameMap = make(map[string]bool)
}