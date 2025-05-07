package main

import (
	"OOP-polygon/exchanges"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
)

type PriceAggregator interface {
	CreaterDeviationMap() (float64, DeviationMap)
	GetPrices(avp float64, DeviationMap DeviationMap, ab_deviation float64) error
	GetAverage(avp float64, DeviationMap DeviationMap, ab_deviation float64) error
}

type DeviationMap = map[string]float64

type TradePrice struct {
	Bybit_Price   string
	Binance_Price string
	Mexc_Price    string
}

func NewTradePrice(opt ...func(t *TradePrice)) *TradePrice {
	tradePrice := &TradePrice{}

	for _, op := range opt {
		op(tradePrice)
	}

	return tradePrice
}

func WithBybit(bybit exchanges.Bybit) func(t *TradePrice) {
	return func(t *TradePrice) {
		t.Bybit_Price = bybit.Result.Bybit_Price[0].Price
	}
}

func WithBinance(binance exchanges.Binance) func(t *TradePrice) {
	return func(t *TradePrice) {
		t.Binance_Price = binance.Binance_Price[0].Price
	}
}

func WithMexc(mexc exchanges.Mexc) func(t *TradePrice) {
	return func(t *TradePrice) {
		t.Mexc_Price = mexc.Mexc_Price[0].Price
	}
}

func main() {
	client := exchanges.NewClient(exchanges.WithClient())
	PriceAggregator, err := ProcessData(client)
	if err != nil {
		log.Println(err)
		return
	}

	avp, DeviationMap := PriceAggregator.CreaterDeviationMap()

	err = PriceAggregator.GetPrices(avp, DeviationMap, 0.5)
	if err != nil {
		log.Println(err)
		return
	}

	err = PriceAggregator.GetAverage(avp, DeviationMap, 0.5)
	if err != nil {
		log.Println(err)
		return
	}
}

func ProcessData(e exchanges.Exchanger) (PriceAggregator, error) {
	bybit, err := e.TradesBybit("ETHUSDT", 10)
	if err != nil {
		return nil, err
	}

	binance, err := e.TradesBinance("ETHUSDT", 10)
	if err != nil {
		return nil, err
	}

	mexc, err := e.TradesMexc("ETHUSDT", 10)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(bybit, exchanges.Bybit{}) || reflect.DeepEqual(binance, exchanges.Binance{}) || reflect.DeepEqual(mexc, exchanges.Mexc{}) {
		return nil, fmt.Errorf("ошибка парсинга")
	}

	tradePrice := NewTradePrice(WithBybit(bybit), WithBinance(binance), WithMexc(mexc))

	return tradePrice, nil
}

func (t TradePrice) CreaterDeviationMap() (float64, DeviationMap) {
	DeviationMap := make(map[string]float64)
	DeviationMap["Bybit"], _ = strconv.ParseFloat(t.Bybit_Price, 64)
	DeviationMap["Binance"], _ = strconv.ParseFloat(t.Binance_Price, 64)
	DeviationMap["Mexc"], _ = strconv.ParseFloat(t.Mexc_Price, 64)
	avp := (DeviationMap["Bybit"] + DeviationMap["Binance"] + DeviationMap["Mexc"]) / 3

	return avp, DeviationMap
}

func (t TradePrice) GetPrices(avp float64, DeviationMap DeviationMap, ab_deviation float64) error {
	for key, value := range DeviationMap {
		deviation := math.Abs((value - avp) / avp * 100)
		if deviation > ab_deviation {
			fmt.Printf("%s: значительное отклонение от средней цены трёх бирж - %.6f%%, допустимое %.2f%%\n", key, deviation, ab_deviation)
		} else {
			fmt.Printf("%s: незначительное отклонение от средней цены трёх бирж - %.6f%%, допустимое %.2f%%\n", key, deviation, ab_deviation)
		}
	}
	return nil
}

func (t TradePrice) GetAverage(avp float64, DeviationMap DeviationMap, ab_deviation float64) error {
	fmt.Printf("\nСредняя цена: %.2f\n", avp)
	for key, value := range DeviationMap {
		deviation := (value - avp) / avp * 100
		fmt.Printf("%s (Прайс %.2f): отклонение от среднего прайса трёх бирж %.4f%%\n", key, DeviationMap[key], deviation)
	}

	return nil
}
