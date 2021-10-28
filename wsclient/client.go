package wsclient

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	Url         url.URL
	Conn        *websocket.Conn
	Subscribers *[]SubscribeModel
}

// {
// 	"method": "SUBSCRIBE",
// 	"params":
// 	[
// 	"btcusdt@aggTrade",
// 	"btcusdt@depth"
// 	],
// 	"id": 1
// 	}

type Kline struct {
	StartTime                int    `json:"t"`
	CloseTime                int    `json:"T"`
	Symbol                   string `json:"s"`
	Interval                 string `json:"i"`
	FirstTradeId             int    `json:"f"`
	LastTradeId              int    `json:"L"`
	OpenPrice                string `json:"o"`
	ClosePrice               string `json:"c"`
	HighPrice                string `json:"h"`
	LowPrice                 string `json:"l"`
	BaseAssetVolume          string `json:"v"`
	NumberOfTrades           int    `json:"n"`
	IsThisKlineClose         bool   `json:"x"`
	QuoteAssetVolume         string `json:"q"`
	TakerBuyBaseAssetVolume  string `json:"V"`
	TakerBuyQuoteAssetVolume string `json:"Q"`
	Ingore                   string `json:"B"`
	// {

	// 		"t": 123400000, // Kline start time
	// 		"T": 123460000, // Kline close time
	// 		"s": "BTCUSDT",  // Symbol
	// 		"i": "1m",      // Interval
	// 		"f": 100,       // First trade ID
	// 		"L": 200,       // Last trade ID
	// 		"o": "0.0010",  // Open price
	// 		"c": "0.0020",  // Close price
	// 		"h": "0.0025",  // High price
	// 		"l": "0.0015",  // Low price
	// 		"v": "1000",    // Base asset volume
	// 		"n": 100,       // Number of trades
	// 		"x": false,     // Is this kline closed?
	// 		"q": "1.0000",  // Quote asset volume
	// 		"V": "500",     // Taker buy base asset volume
	// 		"Q": "0.500",   // Taker buy quote asset volume
	// 		"B": "123456"   // Ignore
	// 	}
}

type KLineResponse struct {
	EventType string `json:"e"`
	EventTime int    `json:"E"`
	Symbol    string `json:"s"`
	Kln       Kline  `json:"k"`
	// 	"e": "kline",     // Event type
	// 	"E": 123456789,   // Event time
	// 	"s": "BTCUSDT",    // Symbol
	// 	"k": {}
}

func (k KLineResponse) PrettyPrint() []byte {
	jsn, err := json.MarshalIndent(k, "", "\t")
	if err != nil {
		fmt.Println("Could not print")
	}
	return jsn
}

type WsMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

func (m WsMessage) toString() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error while convert WsMessage to json")
		return "", err
	}
	return string(b), nil
}

func (c *Client) New() {
	client, _, err := websocket.DefaultDialer.Dial(c.Url.String(), nil)
	c.Conn = client
	if err != nil {
		fmt.Println(err)
	}
	// defer client.Close()

	// done := make(chan struct{})
}

func (c *Client) SendMessage(msg WsMessage) {
	msgStr, err := msg.toString()
	if err != nil {
		fmt.Println("cannot convert WsMessage to JSON byte", err)
	}
	fmt.Printf("msgStr: %v\n", msgStr)
	err = c.Conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Error while sending msg", err)
	}
}

func (c *Client) GetNextSubscribeId() int {
	if len((*c.Subscribers)) <= 0 {
		return 1
	}
	lastSub := (*c.Subscribers)[len((*c.Subscribers))-1]
	id := lastSub.Id
	if id > 9999 {
		id = 0
	} else {
		id++
	}
	return id
}

func (c *Client) IsExistingSubscribe(coin string) bool {
	if len((*c.Subscribers)) <= 0 {
		return false
	}
	isSub := false
	for _, sub := range *c.Subscribers {
		if coin == sub.Coin {
			isSub = true
			break
		}
	}
	fmt.Println("IsSub = ", isSub)
	return isSub
}

func (c *Client) GetSubscribeFromCoinName(coin string) (*SubscribeModel, int, bool) {
	for index, sub := range *c.Subscribers {
		if sub.Coin == coin {
			return &sub, index, true
		}
	}
	return nil, -1, false
}

func (c *Client) RemoveSubscribeAtIndex(index int) {
	start := (*c.Subscribers)[:index]
	end := (*c.Subscribers)[index+1:]
	*c.Subscribers = append(start, end...)
}
