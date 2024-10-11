package util

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
