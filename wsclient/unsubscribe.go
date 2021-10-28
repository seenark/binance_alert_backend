package wsclient

import "fmt"

func (c *Client) Unsubscribe(coins []string) {
	kline := "@kline_1m"
	for _, coin := range coins {
		sub, index, ok := c.GetSubscribeFromCoinName(coin)
		if !ok {
			continue
		}
		params := []string{}
		params = append(params, fmt.Sprintf("%s%s", sub.Coin, kline))
		msg := WsMessage{
			Method: "UNSUBSCRIBE",
			Params: params,
			Id:     sub.Id,
		}
		fmt.Println("Unsubscribe")
		c.SendMessage(msg)
		c.RemoveSubscribeAtIndex(index)
		fmt.Printf("subs: %v", *c.Subscribers)
	}
}
