package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsCurrencySupported(currency string) bool {
	switch currency {
	case CAD, EUR, USD:
		return true
	}
	return false
}
