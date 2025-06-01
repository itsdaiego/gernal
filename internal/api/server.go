package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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

type GekoCoinChartResponse struct {
	Prices [][]float64 `json:"prices"`
}

// sample for single day: https://api.coingecko.com/api/v3/coins/bitcoin/history?date=01-01-2022
// sample for chart: https://api.coingecko.com/api/v3/coins/bitcoin/market_chart/range?vs_currency=usd&from=1735694041&to=1746407641

var BASE_URL = "http://api.coingecko.com/api/v3/coins"

func convertToTimestamp(dateStr string) string {
	layout := "01-02-2006"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	return fmt.Sprintf("%d", t.Unix())
}

func makeRequest(coinID string, startDate string, endDate string) (*http.Response, error) {
	startDateTimestamp := convertToTimestamp(startDate)
	endDateTimestamp := convertToTimestamp(endDate)

	url := fmt.Sprintf("%s/%s/market_chart/range?vs_currency=usd&from=%s&to=%s", BASE_URL, strings.ToLower(coinID), startDateTimestamp, endDateTimestamp)

	return http.Get(url)
}

var FAKE_DATA = true

func fetchCoinCurrentPrice(coinID string, startDate string, endDate string) ([][]float64, error) {
	var data GekoCoinChartResponse

	if FAKE_DATA {
		file, err := os.ReadFile("mocks/btc.json")

		if err != nil {
			fmt.Println("Error reading mock data file:", err)
			return nil, err
		}

		if err := json.Unmarshal(file, &data); err != nil {
			fmt.Println("Error decoding mock JSON:", err)
			return nil, err
		}

		return data.Prices, nil
	}

	resp, err := makeRequest(coinID, startDate, endDate)

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	prices := data.Prices

	if len(prices) > 0 {
		return prices, nil
	}

	fmt.Println("Price not found")
	return nil, fmt.Errorf("price not found")
}

func FetchCoinCurrentPrice(coinID string, startDate string, endDate string) (float64, error) {
	prices, err := fetchCoinCurrentPrice(coinID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	if len(prices) > 0 {
		currentPrice := prices[len(prices)-1][1] // Get the last price in the list
		return currentPrice, nil
	}

	return 0, fmt.Errorf("no prices found for coin %s", coinID)
}

func FetchCoinPriceByDate(coinID string, startDate string, endDate string) ([][]float64, error) {
	prices, err := fetchCoinCurrentPrice(coinID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func FetchCoins() ([]table.Row, error) {
	// coins := []string{"bitcoin", "ethereum"}
	// coinPrices := map[string]float64{}
	//
	// for _, coin := range coins {
	// 	prices, err := fetchCoinCurrentPrice(coin, "01-01-2025", "05-05-2025")
	// 	if err != nil {
	// 		fmt.Println("Error fetching coin data:", err)
	// 		return nil, err
	// 	}
	//
	// 	for i := 0; i < len(prices); i++ {
	// 		currentPriceTimestamp := fmt.Sprintf("%v", prices[i][0])
	// 		currentPrice := prices[i][1]
	//
	// 		coinPrices[currentPriceTimestamp] = currentPrice
	// 	}
	// }

	// fmt.Println("Price:", coinPrices)

	coinRows := []table.Row{
		{"Bitcoin", "BTC", fmt.Sprintf("%.2f", 60000.0)},
		{"Ethereum", "ETH", fmt.Sprintf("%.2f", 4000.0)},
	}

	return coinRows, nil
}
