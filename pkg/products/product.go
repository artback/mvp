package products

type Product struct {
	Name     string `json:"name"`
	SellerID string `json:"sellerId"`
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
