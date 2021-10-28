package watch

import (
	"binance/alert/wsclient"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Volume struct {
	Coin string
	Vol  float32
	Usr  User
	Id   string
}

func NewVolume(coin string, vol float32, user User, coinList *[]string) *Volume {
	return &Volume{
		Coin: coin,
		Vol:  vol,
		Usr:  user,
		Id:   uuid.NewString(),
	}
}

type PriceHighLow struct {
	High float32
	Low  float32
}

type Price struct {
	Coin  string
	Price float32
	Usr   User
	Id    string
}

func NewPrice(coin string, price float32, user *User) *Price {
	return &Price{
		Coin:  coin,
		Price: price,
		Usr:   *user,
		Id:    uuid.NewString(),
	}
}

type Watcher struct {
	Volumes []Volume
	Prices  []Price
	Coins   CoinList
}

func (w *Watcher) WatchVolume(volChan <-chan float32, volDone chan<- Volume) {
	fmt.Println("Watch started")
	for {
		vol := <-volChan
		for _, watchVol := range w.Volumes {
			if vol >= watchVol.Vol {
				go watchVol.Usr.Ln.SendMessage(fmt.Sprintf("😀 Coin: %s 😁 -> 🔊 Volume: %f 🔊", strings.ToUpper(watchVol.Coin), vol))
				// sent to channel for removing this volume watcher
				volDone <- watchVol
				// fmt.Println("vol:", vol)
			}
		}
	}
}

func (w *Watcher) WatchDoneVolume(volDone <-chan Volume, client *wsclient.Client) {
	for {
		vol := <-volDone
		w.RemoveVolumeWatcher(vol.Id, client)
	}
}

func (w *Watcher) IsCoinWatchForVolume(coin string) bool {
	for _, vol := range w.Volumes {
		if coin == vol.Coin {
			return true
		}
	}
	return false
}

func (w *Watcher) AddNewVolumeWatcher(coin string, vol float32, user User, client *wsclient.Client) {
	newVol := Volume{
		Coin: coin,
		Vol:  vol,
		Usr:  user,
		Id:   uuid.NewString(),
	}
	w.Volumes = append(w.Volumes, newVol)
	isAdded := w.Coins.AddNewCoin(coin)
	fmt.Printf("isAdded: %v\n", isAdded)
	if isAdded {
		client.Subscribe([]string{coin})
	}
}

func (w *Watcher) RemoveVolumeWatcher(volumeId string, client *wsclient.Client) {
	for index, vol := range w.Volumes {
		if volumeId == vol.Id {
			w.Volumes = append(w.Volumes[:index], w.Volumes[index+1:]...)
			if w.IsCoinWatchForPrice(vol.Coin) {
				return
			}
			isRemoved := w.Coins.RemoveCoin(vol.Coin)
			if isRemoved {
				client.Unsubscribe([]string{vol.Coin})
			}
		}
	}
}

func (w *Watcher) WatchPrice(priceChan <-chan PriceHighLow, priceDoneChan chan Price) {
	for {
		priceHL := <-priceChan
		// fmt.Println("price: ", priceHL)
		for _, price := range w.Prices {
			// fmt.Println("low", priceHL.Low <= price.Price)
			// fmt.Println("high", price.Price <= priceHL.High)
			if priceHL.Low <= price.Price && price.Price <= priceHL.High {
				// fmt.Println("Send Msg", price.Usr)

				go price.Usr.Ln.SendMessage(fmt.Sprintf("😀 Coin: %s 😄 -> ฿ Price : %v ฿", strings.ToUpper(price.Coin), price.Price))
				priceDoneChan <- price
			}
		}
	}
}

func (w *Watcher) WatchPriceDone(priceDoneChan chan Price, client *wsclient.Client) {
	for {
		price := <-priceDoneChan
		w.RemovePriceWatcher(price.Id, client)
	}
}

func (w *Watcher) IsCoinWatchForPrice(coin string) bool {
	for _, price := range w.Prices {
		if coin == price.Coin {
			return true
		}
	}
	return false
}

func (w *Watcher) AddNewPriceWatcher(coin string, price float32, user User, client *wsclient.Client) {
	newPrice := NewPrice(coin, price, &user)
	w.Prices = append(w.Prices, *newPrice)
	isAdded := w.Coins.AddNewCoin(coin)
	fmt.Printf("isAdded: %v\n", isAdded)
	if isAdded {
		client.Subscribe([]string{coin})
	}
}

func (w *Watcher) RemovePriceWatcher(priceId string, client *wsclient.Client) {
	for index, price := range w.Prices {
		if price.Id == priceId {
			w.Prices = append(w.Prices[:index], w.Prices[index+1:]...)
			if w.IsCoinWatchForVolume(price.Coin) {
				return
			}
			isRemove := w.Coins.RemoveCoin(price.Coin)
			if isRemove {
				client.Unsubscribe([]string{price.Coin})
			}
		}
	}
}
