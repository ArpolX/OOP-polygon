package exchanges

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const bybit = "https://api.bybit.com/v5/market/recent-trade?category=spot&"

type Bybit struct {
	Result struct {
		Bybit_Price []BybitPrice `json:"list"`
	} `json:"result"`
}

type BybitPrice struct {
	Price string `json:"price"`
	Time  string `json:"time"`
}

func (c *Client) TradesBybit(symbol string, limit int) (Bybit, error) {
	UrlRes := fmt.Sprintf("%ssymbol=%s&limit=%d", bybit, symbol, limit)

	req, err := c.Client.Get(UrlRes)
	if err != nil {
		return Bybit{}, err
	}
	defer req.Body.Close()

	var bybit Bybit
	if err := json.NewDecoder(req.Body).Decode(&bybit); err != nil {
		return Bybit{}, err
	}

	for index, i := range bybit.Result.Bybit_Price {
		i1, _ := strconv.ParseInt(string(i.Time), 10, 64)
		Time_Format := time.Unix(0, i1*int64(time.Millisecond))
		bybit.Result.Bybit_Price[index].Time = Time_Format.Format("2006-01-02 15:04:05")
	}
	return bybit, nil
}
