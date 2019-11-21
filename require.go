package validator

import "fmt"

// RequireValidator 判断是否存在
type RequireValidator struct{}

// NewRequireValidator 生成数字校验器
// eg：validate:"require"
func NewRequireValidator() Validator {
	return RequireValidator{}
}

func (v RequireValidator) Validate(val interface{}) (bool, error) {
	if val == nil || val == "" {
		return false, fmt.Errorf("数据校验错误")
	}
	return true, nil
}
