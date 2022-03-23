package products

type Product struct {
	Name     string `json:"name"`
	SellerId string `json:"seller_id"`
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
type Update struct {
	Name   string `json:"name,omitempty" `
	Price  int    `json:"price,omitempty"`
	Amount int    `json:"amount,omitempty"`
}
