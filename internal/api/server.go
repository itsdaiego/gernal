package api

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/charmbracelet/bubbles/table"
)

type GeckoCoinResponse struct {
	MarketData struct {
		CurrentPrice struct {
			USD float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

// var BASE_URL = "https://api.coingecko.com/api/v3/coins/bitcoin/history?date=01-01-2022"
var BASE_URL = "https://api.coingecko.com/api/v3/coins"

func makeRequest(url string, coinID string, date string) (*http.Response, error) {
	return http.Get(fmt.Sprintf("%s/%s/history?date=%s", BASE_URL, coinID, date))
}

func fetchCoinCurrentPrice(coinID string, date string) (float64, error) {
	resp, err := makeRequest(BASE_URL, "bitcoin", "01-01-2025")

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return 0.0, err
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var data GeckoCoinResponse

	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0.0, err
	}

	price := data.MarketData.CurrentPrice.USD

	if price != 0.0 {
		return price, nil
	}

	fmt.Println("Price not found")
	return 0.0, nil
}

func FetchCoins() ([]table.Row, error) {
	price, err := fetchCoinCurrentPrice("bitcoin", "01-01-2025")
	if err != nil {
		return nil, err
	}
	fmt.Println("Price:", price)

	coins := []table.Row{
		{"Bitcoin", "BTC", fmt.Sprintf("%.2f", 60000.0)},
		{"Ethereum", "ETH", fmt.Sprintf("%.2f", 4000.0)},
		{"Litecoin", "LTC", fmt.Sprintf("%.2f", 200.0)},
		{"Cardano", "ADA", fmt.Sprintf("%.2f", 1.20)},
		{"Polkadot", "DOT", fmt.Sprintf("%.2f", 30.50)},
	}

	return coins, nil
}
