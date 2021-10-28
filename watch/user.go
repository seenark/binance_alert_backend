package watch

import "binance/alert/line"

type User struct {
	Ln line.Line
}

func NewUser(lineToken string) *User {
	return &User{
		Ln: line.Line{
			Token: lineToken,
		},
	}
}
