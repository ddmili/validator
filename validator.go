package validator

import (
	"fmt"
	"reflect"
	"strings"
)

// 标签字段
const tagName = "validate"

// Validator 验证接口
type Validator interface {
	Validate(interface{}) (bool, error)
}

// getValidatorFromTag 获取验证规则
// eg：validate:"string,min=1,max=10|require"
func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number": //数字
		validator := NumberValidator{}
		//将structTag中的min和max解析到结构体中
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "string": //字符串
		validator := StringValidator{}
		//将structTag中的min和max解析到结构体中
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "require":
		return RequireValidator{}
	default:
		return NewRuleValidator(args[0])
	}
	//return DefaultValidator{}
}

// Verification 验证
func Verification(s interface{}) []error {
	errs := []error{}
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		// 利用反射获取structTag
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		//多个校验器用 | 分割
		if strings.Contains(tag, "|") {
			a := strings.Split(tag, "|")
			for _, k := range a {
				errs = append(errs, checking(v.Field(i).Interface(), k)...)
			}
		} else {
			errs = append(errs, checking(v.Field(i).Interface(), tag)...)
		}

	}
	return errs
}

// checking 执行检查
func checking(val interface{}, tag string) (errs []error) {
	validator := getValidatorFromTag(tag)
	valid, err := validator.Validate(val)

	if !valid && err != nil {
		errs = append(errs, err)
	}
	return
}
