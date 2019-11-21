package validator

import (
	"fmt"
)

// NumberValidator 数字校验器
type NumberValidator struct {
	Min int
	Max int
}

// NewNumberValidator 生成数字校验器
// eg：validate:"number,min=1,max=10"
func NewNumberValidator(tag string) Validator {
	validator := NumberValidator{}
	fmt.Sscanf(tag, "min=%d,max=%d", &validator.Min, &validator.Max)
	return validator
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	num := val.(int)

	if num < v.Min {
		return false, fmt.Errorf("应大于 %v", v.Min)
	}

	if v.Max >= v.Min && num > v.Max {
		return false, fmt.Errorf("应小于 %v", v.Max)
	}

	return true, nil
}
