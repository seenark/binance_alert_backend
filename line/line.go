package line

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Line struct {
	Token string
}

const (
	lineEndPoint = "https://notify-api.line.me/api/notify"
)

func (l Line) SendMessage(message string) {
	auth := fmt.Sprintf("Bearer %s", l.Token)
	data := url.Values{}
	data.Set("message", message)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, lineEndPoint, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", auth)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Something error while send message")
	}
	defer response.Body.Close()
}
