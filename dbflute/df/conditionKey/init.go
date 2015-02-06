package conditionKey

import ()

func init() {
	CK_EQ = new(CK_EQ_T)
	CK_EQ.ConditionKeyS = C_EQ
	CK_EQ.Operand = "="
	CK_GT = new(CK_GT_T)
	CK_GT.ConditionKeyS = C_GT
	CK_GT.Operand = ">"
}
