package change

import (
	"github.com/artback/mvp/pkg/coin"
	"sort"
)

type Deposit map[coin.Coin]int

func (d Deposit) ToAmount() int {
	var amount int
	for c, a := range d {
		amount += int(c) * a
	}
	return amount
}
func New(coins coin.Coins, amount int) Deposit {
	d := make(Deposit, len(coins))
	sort.Sort(coins)
	for _, c := range coins {
		if amount/c > 0 {
			d[coin.Coin(c)], amount = amount/c, amount%c
		}
	}
	return d
}
