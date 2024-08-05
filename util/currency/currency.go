package currency

// Constants for all valid supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	JPY = "JPY"
	BRL = "BRL"
)

func IsSupported(currency string) bool {
	switch currency {
	case USD, EUR, GBP, JPY, BRL:
		return true
	}
	return false
}
