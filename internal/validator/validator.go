package validator

import (
	"regexp"
)

// sanity check email address
// expression from: https://html.spec.whatwg.org/#valid-e-mail-address
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator type to hold map of errors
type Validator struct {
	Errors map[string]string
}

// creates a new Validator instance
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid() returns true if errors map is empty
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Adds error message to map if not already exists
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Adds error message if validation check is not ok
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Returns true if specific value is in list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Returns true if string matches regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// returns true if all values in slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
