package payment

type PaymentRequest struct {
	BankAccountId         string  `json:"bankAccountId" binding:"required" validate:"required"`
	PaymentProofImageUrl  string  `json:"paymentProofImageUrl" binding:"required" validate:"required"`
	Quantity 							int 	  `json:"quantity" binding:"required,min=1" validate:"required,min=1"`
}
