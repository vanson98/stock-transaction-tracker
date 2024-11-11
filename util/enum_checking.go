package util

import db "stt/database/postgres/sqlc"

const (
	USD = "USD"
	VND = "VND"
	EUR = "EUR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, VND, EUR:
		{
			return true
		}
	}
	return false
}

func IsSupportedTradeType(trade string) bool {
	switch trade {
	case string(db.TradeTypeBUY), string(db.TradeTypeSELL):
		{
			return true
		}
	}
	return false
}
