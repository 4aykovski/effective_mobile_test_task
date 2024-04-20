package filter

import "fmt"

const (
	DataTypeStr = "string"
	DataTypeInt = "int"

	OperatorEq            = "eq"
	OperatorNotEq         = "neq"
	OperatorLowerThan     = "lt"
	OperatorLowerThanEq   = "lte"
	OperatorGreaterThan   = "gt"
	OperatorGreaterThanEq = "gte"
)

type options struct {
	fields []Field
}

func NewOptions() Options {
	return &options{}
}

type Field struct {
	Name  string
	Op    string
	Value string
	Type  string
}

type Options interface {
	AddField(name, op, value, type_ string) error
	Fields() []Field
}

func (o *options) AddField(name, op, value, type_ string) error {

	if err := validateOperator(op); err != nil {
		return fmt.Errorf("can't add field: %w", err)
	}

	o.fields = append(o.fields, Field{
		Name:  name,
		Op:    op,
		Value: value,
		Type:  type_,
	})

	return nil
}
func (o *options) Fields() []Field {
	return o.fields
}

func ParseOperator(op string) (string, error) {
	switch op {
	case OperatorEq:
		return "=", nil
	case OperatorNotEq:
		return "!=", nil
	case OperatorLowerThan:
		return "<", nil
	case OperatorLowerThanEq:
		return "<=", nil
	case OperatorGreaterThan:
		return ">", nil
	case OperatorGreaterThanEq:
		return ">=", nil
	}
	return "", fmt.Errorf("bad operator")
}

func validateOperator(op string) error {
	switch op {
	case OperatorEq, OperatorNotEq, OperatorLowerThan, OperatorLowerThanEq, OperatorGreaterThan, OperatorGreaterThanEq:
		return nil
	}
	return fmt.Errorf("bad operator")
}
