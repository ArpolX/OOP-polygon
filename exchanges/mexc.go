package exchanges

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const mexc = "https://api.mexc.com/api/v3/trades?"

type Mexc struct {
	Mexc_Price []MexcPrice
}

type MexcPrice struct {
	Price string  `json:"price"`
	Time  float64 `json:"time"`
	Times string
}

func (c *Client) TradesMexc(symbol string, limit int) (Mexc, error) {
	UrlRes := fmt.Sprintf("%ssymbol=%s&limit=%d", mexc, symbol, limit)
	fmt.Println(UrlRes)

	req, err := c.Client.Get(UrlRes)
	if err != nil {
		return Mexc{}, err
	}
	defer req.Body.Close()

	var mexc Mexc
	if err := json.NewDecoder(req.Body).Decode(&mexc.Mexc_Price); err != nil {
		return Mexc{}, err
	}

	for index, i := range mexc.Mexc_Price {
		i1_1 := strconv.Itoa(int(i.Time))
		i1, _ := strconv.ParseInt(i1_1, 10, 64)
		Time_Format := time.Unix(0, i1*int64(time.Millisecond))
		mexc.Mexc_Price[index].Times = Time_Format.Format("2006-01-02 15:04:05")
	}
	return mexc, nil
}
