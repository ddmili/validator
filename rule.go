package validator

import (
	"fmt"
	"regexp"
)

// ruleArr 校验规则
var ruleArr = map[string]*regexp.Regexp{
	"require": regexp.MustCompile(`.+`),                                                                                                                                                                 //是否必须
	"email":   regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`),                                                                                                                      //email
	"url":     regexp.MustCompile(`^[a-zA-z]+://[^\s]*`),                                                                                                                                                //url
	"ip":      regexp.MustCompile(`((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`),                                                                                                  //ip
	"chinese": regexp.MustCompile("^[\u4e00-\u9fa5]$"),                                                                                                                                                  //中文
	"mobile":  regexp.MustCompile(`^((\(\d{2,3}\))|(\d{3}\-))?13\d{9}$`),                                                                                                                                //手机号
	"qq":      regexp.MustCompile(`^[1-9]*[1-9][0-9]*$`),                                                                                                                                                //qq号
	"account": regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,20}$`),                                                                                                                                       //帐号(字母开头，允许5-16字节，允许字母数字下划线)
	"string":  regexp.MustCompile("^[a-zA-Z0-9_,\u4e00-\u9fa5-%|]{1,200}$"),                                                                                                                             //字符串(允许字母数字下划线)
	"app_id":  regexp.MustCompile(`^[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+$`),                                                                                                                      //字符串(允许字母数字下划线)
	"number":  regexp.MustCompile(`^\d+$`),                                                                                                                                                              //数字
	"int":     regexp.MustCompile(`^\d+$`),                                                                                                                                                              //数字
	"time":    regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{1,2}\s+[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}$`),                                                                                                  //日期
	"sql":     regexp.MustCompile(`(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`), //判断是否有sql 注入
}

type RuleValidator struct {
	re  *regexp.Regexp
	Val string
}

// NewStringValidator 生成字符串校验器
// eg：validate:"string,min=1,max=10"
func NewRuleValidator(tag string) (validator Validator) {
	validator = DefaultValidator{}
	if r, ok := ruleArr[tag]; ok {
		validator = RuleValidator{r, ""}
	}

	//将structTag中的min和max解析到结构体中
	fmt.Sscanf(tag, "min=%d,max=%d", &validator.Min, &validator.Max)
	return validator
}

func (v RuleValidator) Validate(val interface{}) (bool, error) {
	if val == "" || val == nil {
		return true, fmt.Errorf("数据校验错误", v.re)
	}
	if !v.re.MatchString(val.(string)) {
		return false, fmt.Errorf("数据校验错误", v.re)
	}
	return true, nil
}
