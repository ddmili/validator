//注册自定义验证函数
package validator

import (
	"sync"
)

// 自定义校验函数
var validateFun = map[string]ValidateFun{}

type ValidateFun func(interface{}) (bool, error)

// CustomValidator 自定义校验函数
type CustomValidator struct{
	run func(interface{}) (bool, error)
}

// NewCustomValidator 生成自定义校验器
// eg：validate:"require"
func NewCustomValidator(fun ValidateFun) Validator {
	return CustomValidator{run:fun}
}

// Validate 校验
func (v CustomValidator) Validate(val interface{}) (bool, error) {
	return v.run(val)
}

// RegisterVerification 注册验证函数
func RegisterVerification(name string,fun ValidateFun)  {
	lock := new(sync.RWMutex)
	lock.Lock()
	defer lock.Unlock()
	validateFun[name] = fun
}