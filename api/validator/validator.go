package controler_validator

import (
	"stt/util"

	"github.com/go-playground/validator/v10"
)

var ValidCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var ValidTradeType validator.Func = func(fl validator.FieldLevel) bool {
	if trade, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedTradeType(trade)
	}
	return false
}
