package validator

import (
	"regexp"
)
type Validator interface {
	Validate(input string) bool
}

type LengthValidator struct {
	min int
	max int
}

type PatternValidator struct {
	pattern *regexp.Regexp
}

func NewLengthValidator(min, max int) *LengthValidator {
	return &LengthValidator{min: min, max: max}
}

func (v *LengthValidator) Validate(input string) bool {
	length := len(input)
	return length >= v.min && length <= v.max
}

func NewPatternValidator(pattern string) *PatternValidator {
	return &PatternValidator{pattern: regexp.MustCompile(pattern)}
}

func (v *PatternValidator) Validate(input string) bool {
	return v.pattern.MatchString(input)
}