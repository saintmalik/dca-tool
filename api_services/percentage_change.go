package api_service

import (
	"fmt"
	"io"
	"net/http"
)

func PercentageData(id string) ([]byte, error) {
	//get historical data of coin
	res, err := http.Get("https://api.binance.com/api/v3/ticker/24hr?symbol=" + id)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer res.Body.Close()
	percentage_data_body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error:", err)
	}

	return percentage_data_body, err
}
