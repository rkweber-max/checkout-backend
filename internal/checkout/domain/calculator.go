package domain

func CalculateTotalPrice(products []float64, paymentType PaymentType) float64 {
	var total float64
	for _, price := range products {
		total += price
	}

	if paymentType == PaymentCreditCard {
		total *= 1.03
	}

	return total
}
