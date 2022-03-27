package products

type Product struct {
	Name     string `json:"name"`
	SellerID string `json:"seller_id"`
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
