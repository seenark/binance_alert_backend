package watch

type CoinList struct {
	coins []string
}

func (c *CoinList) AddNewCoin(newCoin string) bool {
	result := false
	if len(c.coins) == 0 {
		c.coins = append(c.coins, newCoin)
		return true
	}

	for _, coin := range c.coins {
		if coin != newCoin {
			c.coins = append(c.coins, newCoin)
			return true
		}
	}
	return result
}

func (c *CoinList) RemoveCoin(coinName string) bool {
	if len(c.coins) == 0 {
		return false
	}
	for index, coin := range c.coins {
		if coinName == coin {
			c.coins = append(c.coins[:index], c.coins[index+1:]...)
			return true
		}
	}
	return false
}
