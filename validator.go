package validator

import (
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
	// args[0] 为校验类型
	//判断自定义校验函数里面有没有
	fun,ok := validateFun[args[0]]
	if ok{
		return NewCustomValidator(fun)
	}
	switch args[0] {
	case "number": //数字
		return NewNumberValidator(strings.Join(args[1:], ","))
	case "string": //字符串
		return NewStringValidator(strings.Join(args[1:], ","))
	default:
		return NewRuleValidator(args[0])
	}
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
