package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/RenanWinter/bank/util/currency"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return currency.IsSupported(value)
	}
	return false
}
