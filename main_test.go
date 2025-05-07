package main

import (
	"OOP-polygon/exchanges"
	"OOP-polygon/mocks"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestProcessData_MockTest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "mock_test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var bybit exchanges.Bybit
			MockExchanger := mocks.NewMockExchanger(ctrl)
			MockExchanger.EXPECT().
				TradesBybit("ETHUSDT", 10).
				Times(1).
				Return(bybit, nil)

			var binance exchanges.Binance
			MockExchanger.EXPECT().
				TradesBinance("ETHUSDT", 10).
				Times(1).
				Return(binance, nil)

			var mexc exchanges.Mexc
			MockExchanger.EXPECT().
				TradesMexc("ETHUSDT", 10).
				Times(1).
				Return(mexc, nil)

			priceAggregator, err := ProcessData(MockExchanger)
			if priceAggregator != nil {
				t.Fail()
			} else if error.Error(err) != "ошибка парсинга" {
				t.Fail()
			}
		})
	}
}

func TestTradePrice_CreaterDeviationMap(t *testing.T) {
	tests := []struct {
		name  string
		tr    TradePrice
		want  float64
		want1 DeviationMap
	}{
		{name: "DeviationMap", tr: TradePrice{
			Bybit_Price:   "100",
			Binance_Price: "100",
			Mexc_Price:    "100",
		}, want: 100, want1: map[string]float64{
			"Bybit":   100,
			"Binance": 100,
			"Mexc":    100,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tr.CreaterDeviationMap()
			if got != tt.want {
				t.Errorf("TradePrice.CreaterDeviationMap() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TradePrice.CreaterDeviationMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTradePrice_GetPrices(t *testing.T) {
	type args struct {
		avp          float64
		DeviationMap DeviationMap
	}
	tests := []struct {
		name    string
		tr      TradePrice
		args    args
		wantErr bool
	}{
		{name: "GetPrices", tr: TradePrice{
			Bybit_Price:   "100",
			Binance_Price: "100",
			Mexc_Price:    "100",
		}, args: args{
			avp: 100,
			DeviationMap: map[string]float64{
				"Bybit":   100,
				"Binance": 100,
				"Mexc":    100,
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.GetPrices(tt.args.avp, tt.args.DeviationMap, 0); (err != nil) != tt.wantErr {
				t.Errorf("TradePrice.GetPrices() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTradePrice_GetAverage(t *testing.T) {
	type args struct {
		avp          float64
		DeviationMap DeviationMap
	}
	tests := []struct {
		name    string
		tr      TradePrice
		args    args
		wantErr bool
	}{
		{name: "GetPrices", tr: TradePrice{
			Bybit_Price:   "100",
			Binance_Price: "100",
			Mexc_Price:    "100",
		}, args: args{
			avp: 100,
			DeviationMap: map[string]float64{
				"Bybit":   100,
				"Binance": 100,
				"Mexc":    100,
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.GetAverage(tt.args.avp, tt.args.DeviationMap, 0); (err != nil) != tt.wantErr {
				t.Errorf("TradePrice.GetAverage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
