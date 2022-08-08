package model

var (
	Alloted             string
	Coinid              string
	Amount              string
	Buyinterval         string
	Fee                 string
	Bapi                string
	Bsecret             string
	Testvalue           string
	Priceperbuy         string
	Config_allocated    float64
	Config_amount       float64
	Config_buyintervals int
	Config_fee          float64
)

type Config struct {
	Coins        string  `json:"coins"`
	Buyintervals int     `json:"buyintervals"`
	Amount       float64 `json:"amount"`
	Percent      float64 `json:"percent"`
	Fee          float64 `json:"fee"`
	Testing      bool    `json:"testing"`
} // struct to hold config data

type CoinPercentageData struct {
	Symbol             string `json:"symbol"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	OpenPrice          string `json:"openPrice"`
}
