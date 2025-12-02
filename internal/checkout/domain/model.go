package domain

type PaymentType string

const (
	PaymentPix        PaymentType = "pix"
	PaymentBoleto     PaymentType = "boleto"
	PaymentCreditCard PaymentType = "credit_card"
)

type CheckoutRequest struct {
	ProductIDs  []int        `json:"product_ids" binding:"required"`
	PaymentType PaymentType  `json:"payment_type" binding:"required"`
	Customer    CustomerInfo `json:"customer" binding:"required"`
}

type CustomerInfo struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type Order struct {
	ID          int64        `json:"id"`
	Total       float64      `json:"total"`
	PaymentType PaymentType  `json:"payment_type"`
	Customer    CustomerInfo `json:"customer"`
}
