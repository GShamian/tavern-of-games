package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Special function that helps to validate original password to it's encrypted version
func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}
