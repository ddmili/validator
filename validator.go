package validators

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// tagName 标签字段
const tagName = "validate"

// ValidatorHandler 校验基类
type ValidatorHandler struct {
}

// Validator 验证接口
type Validator interface {
	Validate() (bool, error)
}

type ValidateFun func(string, string, ...string) (bool, error)
type clipperFun func(reflect.Value, ...string) (bool, error)

// Validate 校验
func (v CustomValidator) Validate() (bool, error) {
	t := v.val.Field(v.i).Type()
	val := strings.TrimSpace(v.val.Field(v.i).String())

	if t.Kind() == reflect.Int {
		val = strconv.FormatInt(v.val.Field(v.i).Int(), 10)
	}
	return v.run(v.val.Type().Field(v.i).Name, val, v.params...)
}

// CustomValidator 自定义校验函数
type CustomValidator struct {
	params []string
	i      int
	val    reflect.Value
	run    ValidateFun
}

// clipperValidator 裁剪器
type clipperValidator struct {
	params []string
	i      int
	val    reflect.Value
	run    clipperFun
}

// Validate 校验
func (v clipperValidator) Validate() (bool, error) {
	return v.run(v.val.Field(v.i), v.params...)
}

// NewCustomValidator 生成自定义校验器
// eg：validate:"require"
func NewCustomValidator(i int, val reflect.Value, fun ValidateFun, params ...string) Validator {
	return CustomValidator{run: fun, val: val, i: i, params: params}
}

// NewClipperValidator 生成裁剪校验器
// eg：validate:"require"
func NewClipperValidator(i int, val reflect.Value, fun clipperFun, params ...string) Validator {
	return clipperValidator{run: fun, val: val, i: i, params: params}
}

//getValidation 获取校验器 eg: enum=1,2,3,4,5
func getValidation(i int, val reflect.Value) (Validator, error) {
	tag := val.Type().Field(i).Tag.Get(tagName)
	a := strings.Split(tag, "=")
	switch a[0] {
	case "require": // validate:"require"
		return NewCustomValidator(i, val, regHandler(`.+`, errors.New("参数不能为空"))), nil
	case "time": // validate:"time"
		return NewCustomValidator(i, val, regHandler(`^[0-9]{4}-[0-9]{2}-[0-9]{1,2}\s+[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}$`, errors.New("参数需为时间格式"))), nil
	case "mobile": // validate:"mobile"
		return NewCustomValidator(i, val, regHandler(`^\d+$`, errors.New("参数需为手机号"))), nil
	case "enum": // validate:"enum=1,2,3,4"
		params := strings.Split(a[1], ",")
		return NewCustomValidator(i, val, enum, params...), nil
	case "trimSpace": // validate:"trimSpace"
		return NewClipperValidator(i, val, trimSpace), nil
	case "default": // validate:"default=1"
		params := strings.Split(a[1], ",")
		return NewClipperValidator(i, val, defaultVal, params...), nil

	}
	return CustomValidator{}, errors.New("未找到校验器")
}

// regHandler 正则匹配
func regHandler(ruleStr string, err error) ValidateFun {
	return func(tag, val string, params ...string) (bool, error) {
		reg := regexp.MustCompile(ruleStr)
		if !reg.MatchString(val) {
			return false, fmt.Errorf("%s|%s 参数错误", tag, val)
		}
		return true, nil
	}
}

// Verification 验证
func Verification(s interface{}) error {
	obj := reflect.ValueOf(s)
	if obj.Kind() != reflect.Struct {
		obj = reflect.Indirect(obj)
	}
	return execute(obj)
}

// execute 执行
func execute(obj reflect.Value) error {
	for i := 0; i < obj.NumField(); i++ {
		if obj.Field(i).Kind() == reflect.Struct {
			//嵌套检查
			err := execute(obj.Field(i))
			if err != nil {
				return err
			}
		}
		// fmt.Println(obj.Type().Field(i).Name)
		// 利用反射获取structTag
		tag := obj.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		//多个校验器用 | 分割
		if strings.Contains(tag, "|") {
			a := strings.Split(tag, "|")
			for _, k := range a {
				err := checking(obj, i, k)
				if err != nil {
					return err
				}
			}
		} else {
			err := checking(obj, i, tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// checking 执行检查
func checking(obj reflect.Value, i int, tag string) error {
	validator, e := getValidation(i, obj)
	if e == nil {
		valid, err := validator.Validate()

		if !valid && err != nil {
			return err
		}
		return nil
	}
	//匹配是否有函数 validate:"selfFunc"
	m := obj.MethodByName(tag)
	if !m.IsValid() {
		return errors.New("未匹配到校验器")
	}
	params := []reflect.Value{}
	rs := m.Call(params)
	if len(rs) > 0 {
		e, ok := rs[0].Interface().(error)
		if ok && e != nil {
			return e
		}
	}
	return nil
}

// enum 枚举校验
func enum(name, val string, params ...string) (bool, error) {
	if val == "" {
		return true, nil
	}
	for _, v := range params {
		if v == val {
			return true, nil
		}
	}
	return false, fmt.Errorf("%s|%s 参数错误", name, val)
}

// defaultVal 默认值
func defaultVal(val reflect.Value, params ...string) (bool, error) {
	if !val.CanSet() || !val.IsZero() {
		return true, nil
	}
	switch val.Kind() {
	case reflect.Uint:
		result, _ := strconv.ParseUint(params[0], 10, 64)
		val.SetUint(result)
	case reflect.String:
		val.SetString(params[0])
	case reflect.Int, reflect.Int8, reflect.Int64, reflect.Int32:
		result, _ := strconv.Atoi(params[0])
		val.SetInt(int64(result))
	}
	return true, nil
}

// trimSpace 删除首尾空格
func trimSpace(val reflect.Value, params ...string) (bool, error) {
	if val.CanSet() && val.Kind() == reflect.String {
		val.SetString(strings.TrimSpace(val.String()))
	}
	return true, nil
}

func Validate(v ...error) error {
	for _, e := range v {
		if e != nil {
			return e
		}
	}
	return nil
}
