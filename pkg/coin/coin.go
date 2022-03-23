package coin

type Coin int
type Coins []int

func (c Coins) Len() int {
	return len(c)
}

func (c Coins) Less(i, j int) bool {
	return c[i] > c[j]
}

func (c Coins) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
