package wsclient

import (
	"fmt"
)

type SubscribeModel struct {
	Coin string
	Id   int
}

func (c *Client) Subscribe(coins []string) {
	kline := "@kline_1m"
	for _, coin := range coins {
		if c.IsExistingSubscribe(coin) {
			continue
		}
		newSub := SubscribeModel{
			Coin: coin,
			Id:   c.GetNextSubscribeId(),
		}
		params := []string{}
		params = append(params, fmt.Sprintf("%s%s", coin, kline))
		msg := WsMessage{
			Method: "SUBSCRIBE",
			Params: params,
			Id:     newSub.Id,
		}
		fmt.Println("Subscribe")
		c.SendMessage(msg)
		*c.Subscribers = append(*c.Subscribers, newSub)
		fmt.Printf("subs: %v", *c.Subscribers)
	}

}
