package payment

import "time"

type PaymentResponse struct {
	PurchaseID 						string 		`json:"purchaseId"`
	BuyerID    						string 		`json:"buyerId"`
	AccountID  						string 		`json:"bankAccountId"`
	ProductID  						string    `json:"productId"`
	Quantity   						int       `json:"quantity"`
	TotalPrice 						int       `json:"totalPrice"`
	ImageUrl   						string    `json:"paymentProofImageUrl"`
	CreatedAt 						time.Time `json:"createdAt"`
	UpdatedAt  						time.Time `json:"updatedAt"`
}
