package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BTCAverage struct {
	Daily float64 `json:"24h_avg"`
}

var rate float64

func mBtc(price float64) string {
	return fmt.Sprintf("%f mBTC", price/rate*1000*(1-(*btcDiscount/100)))
}

func GetRate() error {
	resp, err := http.Get("https://api.bitcoinaverage.com/ticker/global/USD/")
	if err != nil {
		return err
	}
	buf := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		return err
	}

	var avg BTCAverage
	err = json.Unmarshal(buf, &avg)
	if err != nil {
		return err
	}

	rate = avg.Daily
	return nil
}
