package main

import (
	"binance/alert/handlers"
	"binance/alert/watch"
	"binance/alert/wsclient"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	lineToken = "sBV0pqTRkK5NhYtSUnQHLevNWAVJuBsjRsyPfG5nEMy"
)

var Users = []watch.User{}
var watcher = watch.Watcher{
	Volumes: []watch.Volume{},
	Prices:  []watch.Price{},
	Coins:   watch.CoinList{},
}
var wsClient = wsclient.Client{}

func main() {
	HadesGod := watch.NewUser(lineToken)
	Users = append(Users, *HadesGod)

	NewWs()
	// watcher.AddNewVolumeWatcher("dotusdt", 1000, *HadesGod, &wsClient)
	// watcher.AddNewPriceWatcher("dotusdt", 36.1, *HadesGod, &wsClient)

	// watcher.AddNewVolumeWatcher("btcusdt", 1000, *HadesGod, &wsClient)
	NewEcho()

}

func NewWs() {
	volumeChan := make(chan float32)
	volumeDone := make(chan watch.Volume)
	go watcher.WatchVolume(volumeChan, volumeDone)
	go watcher.WatchDoneVolume(volumeDone, &wsClient)

	priceChan := make(chan watch.PriceHighLow)
	priceDoneChan := make(chan watch.Price)
	go watcher.WatchPrice(priceChan, priceDoneChan)
	go watcher.WatchPriceDone(priceDoneChan, &wsClient)

	u := url.URL{
		Scheme: "wss",
		Host:   "fstream.binance.com",
		Path:   "/ws",
	}
	fmt.Println("url:", u.String())

	wsClient = wsclient.Client{
		Url:         u,
		Subscribers: &[]wsclient.SubscribeModel{},
	}
	wsClient.New()
	go HandleMessage(&wsClient, volumeChan, priceChan)

	// wsClient.Subscribe([]string{"dotusdt"})
	// wsClient.Subscribe([]string{"btcusdt", "ethusdt"})
	// wsClient.Unsubscribe([]string{"btcusdt"})

}

func NewEcho() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Homepage")
	})

	e.POST("/volume", func(c echo.Context) error {
		return handlers.HandleVolume(c, &watcher, &Users[0], &wsClient)
	})

	e.POST("/price", func(c echo.Context) error {
		return handlers.HandlerPrice(c, &watcher, &Users[0], &wsClient)
	})

	e.GET("/watcher", func(c echo.Context) error {
		return c.JSON(http.StatusOK, watcher)
	})

	e.DELETE("/volume", func(c echo.Context) error {
		id := c.QueryParam("id")
		watcher.RemoveVolumeWatcher(id, &wsClient)
		return c.JSON(http.StatusOK, id)
	})

	e.DELETE("/price", func(c echo.Context) error {
		id := c.QueryParam("id")
		watcher.RemovePriceWatcher(id, &wsClient)
		return c.JSON(http.StatusOK, id)
	})

	e.Static("/public", "./public")
	e.Start("127.0.0.1:8080")
}

func HandleMessage(c *wsclient.Client, volChan chan<- float32, priceChan chan<- watch.PriceHighLow) {
	// defer close(done)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("error while read message", err)
			break
		}
		kln := wsclient.KLineResponse{}
		err = json.Unmarshal([]byte(message), &kln)
		if err != nil {
			fmt.Println("Could not unmarshall")
		}

		baseVol, err := strconv.ParseFloat(kln.Kln.BaseAssetVolume, 32)
		if err != nil {
			fmt.Println("Cannot convert string to float32", err)
		}
		volChan <- float32(baseVol)

		priceH, err := strconv.ParseFloat(kln.Kln.HighPrice, 32)
		if err != nil {
			fmt.Println("Cannot convert high price to float32")
		}
		priceL, err := strconv.ParseFloat(kln.Kln.LowPrice, 32)
		if err != nil {
			fmt.Println("Cannot convert low price to float32")
		}
		priceHL := watch.PriceHighLow{
			High: float32(priceH),
			Low:  float32(priceL),
		}

		priceChan <- priceHL

	}
}
