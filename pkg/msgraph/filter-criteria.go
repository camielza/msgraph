package msgraph

import "fmt"

type FilterCriteria struct {
}

type Criteria interface {
	String() string
}

const (
	OR  = 1
	AND = 2
	NOT = 3
)

type BinaryLogicOperator struct {
	operator  int8
	criteria1 *Criteria
	criteria2 *Criteria
}

type StartWithCriteria struct {
	Field     string
	StartWith string
}

func (c *FilterCriteria) LogicOr(criteria1 *Criteria, criteria2 *Criteria) *Criteria {
	op := BinaryLogicOperator{OR, criteria1, criteria2}
	criteria := Criteria(op)
	return &criteria
}

func (c *FilterCriteria) LogicAnd(criteria1 *Criteria, criteria2 *Criteria) *Criteria {
	op := BinaryLogicOperator{AND, criteria1, criteria2}
	criteria := Criteria(op)
	return &criteria
}
func (c *FilterCriteria) LogicNot(criteria1 *Criteria) *Criteria {
	op := BinaryLogicOperator{NOT, criteria1, nil}
	criteria := Criteria(op)
	return &criteria
}

func (c *FilterCriteria) StartWith(field string, startWith string) *Criteria {
	start := StartWithCriteria{field, startWith}
	criteria := Criteria(start)
	return &criteria
}

func (c BinaryLogicOperator) String() string {
	format := "%s %s %s"
	op := "AND"
	switch c.operator {
	case AND:
		op = "AND"
	case OR:
		op = "OR"
	case NOT:
		return fmt.Sprintf("NOT %s", (*c.criteria1).String())
	}
	return fmt.Sprintf(format, (*c.criteria1).String(), op, (*c.criteria2).String())
}

func (c StartWithCriteria) String() string {
	return fmt.Sprintf("startswith(%s,'%s')", c.Field, c.StartWith)
}
