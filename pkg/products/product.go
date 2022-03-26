package products

type Product struct {
	Name     string `json:"name"`
	SellerId string `json:"seller_id"`
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
