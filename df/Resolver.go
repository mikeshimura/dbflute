package df

import (

)
const(
	LOCAL_ALIAS_MARK = "$$localAlias$$"
	FOREIGN_ALIAS_MARK = "$$foreignAlias$$"
	SQ_BEGIN_MARK = "$$sqbegin$$"
	SQ_END_MARK = "$$sqend$$"
	INLINE_MARK = "$$inline$$"
	LOCATION_BASE_MARK = "$$locationBase$$"
	OPTIMIZED_MARK = "$$optimized$$"
)

type InlineViewResource struct{
	
}
type FixedConditionResolver interface {
	ResolveVariable(fixedCondition string, fixedInline bool) string
	ResolveFixedInlineView(foreignTable string, treatedAsInnerJoin bool) string
	HasOverRelation(fixedCondition string) bool
}

type HpFixedConditionQueryResolver struct{
	localCQ *ConditionQuery
	foreignCQ *ConditionQuery
	resolvedFixedCondition string
	inlineViewResourceMap map[string]*InlineViewResource
	inlineViewOptimizedCondition string
	inlineViewOptimizationWholeCondition bool
	inlineViewOptimizedLineNumberSet map[int]int
	Resolver *FixedConditionResolver
}
func (p *HpFixedConditionQueryResolver)ResolveVariable(fixedCondition string, fixedInline bool) string{
	
	panic("")
	return ""
}
func (p *HpFixedConditionQueryResolver)ResolveFixedInlineView(foreignTable string, treatedAsInnerJoin bool) string{
	
	panic("")
	return ""
}
func (p *HpFixedConditionQueryResolver)HasOverRelation(fixedCondition string) bool{
	panic("")	
	return false
}
