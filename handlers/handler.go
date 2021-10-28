package handlers

import (
	"binance/alert/watch"
	"binance/alert/wsclient"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type VolumeRequest struct {
	Coin   string  `json:"coin" validate:"lowercase,required"`
	Volume float32 `json:"volume" validate:"numeric,required"`
}

func HandleVolume(c echo.Context, watcher *watch.Watcher, user *watch.User, client *wsclient.Client) error {
	volumeReq := VolumeRequest{}
	err := c.Bind(&volumeReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	validate := validator.New()
	err = validate.Struct(volumeReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	watcher.AddNewVolumeWatcher(volumeReq.Coin, volumeReq.Volume, *user, client)

	return c.JSON(http.StatusOK, volumeReq)
}

type PriceRequest struct {
	Coin  string  `json:"coin" validate:"lowercase,required"`
	Price float32 `json:"price" validate:"numeric,required"`
}

func HandlerPrice(c echo.Context, watcher *watch.Watcher, user *watch.User, client *wsclient.Client) error {
	priceReq := PriceRequest{}
	err := c.Bind(&priceReq)
	if err != nil {
		return echo.ErrBadRequest
	}
	validate := validator.New()
	err = validate.Struct(priceReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	watcher.AddNewPriceWatcher(priceReq.Coin, priceReq.Price, *user, client)
	return c.JSON(http.StatusOK, priceReq)
}
