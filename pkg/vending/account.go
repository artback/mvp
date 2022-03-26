package vending

import (
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/products"
)

type Account struct {
	Deposit  int                `json:"deposit,omitempty"`
	Products []products.Product `json:"products,omitempty"`
	Spent    int                `json:"spent"`
}
type Response struct {
	Deposit  change.Deposit     `json:"deposit,omitempty"`
	Products []products.Product `json:"products,omitempty"`
	Spent    int                `json:"spent"`
}
