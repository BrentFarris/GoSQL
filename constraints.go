package db

const (
	ConditionEquals              = "="
	ConditionNotEquals           = "<>"
	ConditionLessThan            = "<"
	ConditionGreaterThan         = ">"
	ConditionLessThanOrEquals    = "<="
	ConditionGreaterThanOrEquals = ">="
	ConditionLike                = "LIKE"
	ConditionNotLike             = "NOT LIKE"
	ConditionIn                  = "IN"
	ConditionNotIn               = "NOT IN"
	ConditionBetween             = "BETWEEN"
	ConditionNotBetween          = "NOT BETWEEN"
	ConditionIsNull              = "IS NULL"
	ConditionIsNotNull           = "IS NOT NULL"

	ConjunctionAnd = "AND"
	ConjunctionOr  = "OR"
)

type Constraint struct {
	field       string
	condition   string
	conjunction string
	value       any
}

type Constraints struct {
	constraints []Constraint
}

func (c *Constraints) conjunction(field, condition, conjoin string, value any) {
	c.constraints = append(c.constraints, Constraint{
		field:       field,
		condition:   condition,
		conjunction: conjoin,
		value:       value,
	})
}

func (c *Constraints) And(field, condition string, value any) *Constraints {
	c.conjunction(field, condition, ConjunctionAnd, value)
	return c
}

func (c *Constraints) Or(field, condition string, value any) *Constraints {
	c.conjunction(field, condition, ConjunctionOr, value)
	return c
}
