package df

import ()

var CommonColumnAutoSetupper_I *CommonColumnAutoSetupper

type CommonColumnAutoSetupper interface {
	HandleCommonColumnOfInsertIfNeeds(entity *Entity, ctx *Context)
	DoHandleCommonColumnOfInsertIfNeeds(entity *Entity, ctx *Context)
	HandleCommonColumnOfUpdateIfNeeds(entity *Entity, ctx *Context)
	DoHandleCommonColumnOfUpdateIfNeeds(entity *Entity, ctx *Context)
}

type BaseCommonColumnAutoSetupper struct {
	CommonColumnAutoSetupper *CommonColumnAutoSetupper
}

func (p *BaseCommonColumnAutoSetupper) HandleCommonColumnOfInsertIfNeeds(
	entity *Entity, ctx *Context) {
	if entity == nil || ctx == nil {
		return
	}
	if p.checkHasCommonColumn(entity) == false{
		return
	}
	(*p.CommonColumnAutoSetupper).DoHandleCommonColumnOfInsertIfNeeds(entity,ctx)
}
func (p *BaseCommonColumnAutoSetupper) HandleCommonColumnOfUpdateIfNeeds(
	entity *Entity, ctx *Context) {
	if entity == nil || ctx == nil {
		return
	}
	if p.checkHasCommonColumn(entity) == false{
		return
	}
	(*p.CommonColumnAutoSetupper).DoHandleCommonColumnOfUpdateIfNeeds(entity,ctx)
}
func (p *BaseCommonColumnAutoSetupper) checkHasCommonColumn(entity interface{}) bool{
	
	return true
}
