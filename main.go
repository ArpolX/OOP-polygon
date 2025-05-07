package main

import (
	"OOP-polygon/exchanges"
	"fmt"
	"log"
	"math"
	"strconv"
)

type PriceAggregator interface {
	CreaterDeviationMap() DeviationMap
	GetPrices(DeviationMap) error
	GetAverage(DeviationMap) error
}

type DeviationMap = map[string]float64

type TradePrice struct {
	Bybit_Price   string
	Binance_Price string
	Mexc_Price    string
}

func main() {
	client := exchanges.NewClient(exchanges.WithClient())
	err := ProcessData(client)
	if err != nil {
		log.Println(err)
		return
	}
}

func ProcessData(e exchanges.Exchanger) error {
	var tradePrice TradePrice
	bybit, err := e.TradesBybit("ETHUSDT", 10)
	if err != nil {
		return err
	}
	tradePrice.Bybit_Price = bybit.Result.Bybit_Price[0].Price
	fmt.Println(bybit.Result.Bybit_Price)

	binance, err := e.TradesBinance("ETHUSDT", 10)
	if err != nil {
		return err
	}
	tradePrice.Binance_Price = binance.Binance_Price[0].Price
	fmt.Println(binance.Binance_Price)

	mexc, err := e.TradesMexc("ETHUSDT", 10)
	if err != nil {
		return err
	}
	tradePrice.Mexc_Price = mexc.Mexc_Price[0].Price
	fmt.Println(mexc)

	avp, DeviationMap := tradePrice.CreaterDeviationMap()

	err = tradePrice.GetPrices(avp, DeviationMap, 0.5)
	if err != nil {
		return err
	}

	err = tradePrice.GetAverage(avp, DeviationMap, 0.5)
	if err != nil {
		return err
	}

	return nil
}

func (t TradePrice) CreaterDeviationMap() (float64, DeviationMap) {
	DeviationMap := make(map[string]float64)
	DeviationMap["Bybit"], _ = strconv.ParseFloat(t.Bybit_Price, 64)
	DeviationMap["Binance"], _ = strconv.ParseFloat(t.Binance_Price, 64)
	DeviationMap["Maxc"], _ = strconv.ParseFloat(t.Mexc_Price, 64)
	avp := (DeviationMap["Bybit"] + DeviationMap["Binance"] + DeviationMap["Maxc"]) / 3

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
