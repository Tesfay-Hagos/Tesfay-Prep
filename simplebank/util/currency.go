package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	BIR = "BIR"
	NKF = "NAKFA"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, BIR, NKF:
		return true
	}
	return false
}
