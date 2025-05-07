package exchanges

import (
	"net/http"
	"time"
)

//go:generate mockgen -source=client.go -destination=OOP-polygon/mocks/Exchanger.go -package=mocks
type Exchanger interface {
	TradesBybit(symbol string, limit int) (Bybit, error)
	TradesBinance(symbol string, limit int) (Binance, error)
	TradesMexc(symbol string, limit int) (Mexc, error)
}

type Client struct {
	Client http.Client
}

type OptionClient func(*Client)

func NewClient(opt ...OptionClient) *Client {
	client := &Client{}

	for _, op := range opt {
		op(client)
	}

	return client
}

func WithClient() OptionClient {
	return func(c *Client) {
		c.Client = http.Client{
			Timeout: 5 * time.Second,
		}
	}
}
