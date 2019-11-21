package validator

// DefaultValidator 默认校验器，直接通过
type DefaultValidator struct{}

func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}
