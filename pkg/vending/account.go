package vending

import (
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/products"
)

type AccountResponse struct {
	Deposit  change.Deposit    `json:"deposit,omitempty"`
	Products []products.Update `json:"products,omitempty"`
	Spent    int               `json:"spent"`
}
