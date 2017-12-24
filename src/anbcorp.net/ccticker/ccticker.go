package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type TickerResponse struct {
	Last      float64 `json:",string"`
	High      float64 `json:",string"`
	Low       float64 `json:",string"`
	Vwap      float64 `json:",string"`
	Volume    float64 `json:",string"`
	Bid       float64 `json:",string"`
	Ask       float64 `json:",string"`
	Timestamp uint64  `json:",string"`
	Open      float64 `json:",string"`
}

type Settings struct {
	Assets   map[string]float64 `json:assets`
	Currency string             `json:currency`
}

func getTicker(pair string) (ticker TickerResponse, err error) {

	url := fmt.Sprintf("https://www.bitstamp.net/api/v2/ticker_hour/%s/", pair)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&ticker)
	if err != nil {
		return
	}

	return
}

func showAssets(config Settings) {
	holdings := make(map[string]float64)

	for asset, units := range config.Assets {
		pair := strings.ToLower(fmt.Sprintf("%s%s", asset, config.Currency))
		ticker, err := getTicker(pair)
		if err != nil {
			fmt.Println(asset, err)
			continue
		}
		holdings[asset] = ticker.Last * units
	}

	var total float64
	for _, value := range holdings {
		total += value
	}

	fmt.Printf("Total value: $%.2f\n", total)
}

func getConfig() (config Settings, err error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
		return
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		fmt.Println("parsing config file", err.Error())
		return
	}

	return
}

func main() {
	config, err := getConfig()
	if err != nil {
		return
	}

	showAssets(config)
}
