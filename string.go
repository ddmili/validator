package validator

import (
	"fmt"
)

// StringValidator 字符串校验器
type StringValidator struct {
	Min int //最小长度
	Max int //最大长度
}

// NewStringValidator 生成字符串校验器
// eg：validate:"string,min=1,max=10"
func NewStringValidator(tag string) Validator {
	validator := StringValidator{}
	//将structTag中的min和max解析到结构体中
	fmt.Sscanf(tag, "min=%d,max=%d", &validator.Min, &validator.Max)
	return validator
}

// Validate 校验函数
func (v StringValidator) Validate(val interface{}) (bool, error) {
	a,ok := val.(string)
	if !ok{
		return false, fmt.Errorf("不")
	}
	l := len(a)

	if l == 0 {
		return false, fmt.Errorf("不能为空")
	}

	if l <= v.Min {
		return false, fmt.Errorf("字符串太短")
	}

	if v.Max >= v.Min && l > v.Max {
		return false, fmt.Errorf("字符串太长")
	}

	return true, nil
}
