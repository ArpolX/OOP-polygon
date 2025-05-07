package exchanges

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const binance = "https://api.binance.com/api/v3/trades?"

type Binance struct {
	Binance_Price []BinancePrice
}

type BinancePrice struct {
	Price string  `json:"price"`
	Time  float64 `json:"time"`
	Times string
}

func (c *Client) TradesBinance(symbol string, limit int) (Binance, error) {
	UrlRes := fmt.Sprintf("%ssymbol=%s&limit=%d", binance, symbol, limit)

	req, err := c.Client.Get(UrlRes)
	if err != nil {
		return Binance{}, err
	}
	defer req.Body.Close()

	var binance Binance
	if err := json.NewDecoder(req.Body).Decode(&binance.Binance_Price); err != nil {
		return Binance{}, err
	}

	for index, i := range binance.Binance_Price {
		i1_1 := strconv.Itoa(int(i.Time))
		i1, _ := strconv.ParseInt(i1_1, 10, 64)
		Time_Format := time.Unix(0, i1*int64(time.Millisecond))
		binance.Binance_Price[index].Times = Time_Format.Format("2006-01-02 15:04:05")
	}

	return binance, nil
}
