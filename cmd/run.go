package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/adshao/go-binance/v2"
	api_service "github.com/saintmalik/dca-tool/api_services"
	"github.com/saintmalik/dca-tool/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "This command runs the program and execute orders",
	Long:  `This command runs the program and checks the percentage change of the coin and then purchases the coin if the percentage change is less than the set percentage.`,
	Run: func(cmd *cobra.Command, args []string) {
		perform()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Println("Config file successfully read", viper.ConfigFileUsed())
}

func perform() {
	model.Testvalue = viper.GetString("testing")
	if model.Testvalue == "true" {
		fmt.Println("Real orders cannot be made in testing mode")
		crap()
		return
	} else { // if testing is false, then we are in production

		model.Coinid = viper.GetString("coins") + "USDT"
		percentage_data_res_body, err := api_service.PercentageData(model.Coinid) //get current data of coin
		if err != nil {
			log.Fatal("No response from request")
		}

		var reading model.CoinPercentageData
		err = json.Unmarshal(percentage_data_res_body, &reading) // unmarshal JSON into Go value
		if err != nil {
			fmt.Println("error:", err)
		}
		var lik string = reading.PriceChangePercent
		pil, _ := strconv.ParseFloat(lik, 64)
		fmt.Println("percentage change", pil)

		//check type of percentage change
		switch {
		case pil < 0:
			fmt.Println("Coin is down")
			dcaBuy()
		case pil >= 0 && pil < 1.5:
			fmt.Println("Your weight is normal")
			dcaBuy()
		case pil >= 1.6 && pil < 2:
			fmt.Println("Coin is up")
			dcaBuy()
		default:
			fmt.Println("You're obese")
		}
	}
}

func dcaBuy() {
	model.Coinid = viper.GetString("coins") + "USDT"
	model.Config_buyintervals = viper.GetInt("buyintervals")
	model.Config_amount = viper.GetFloat64("amount")
	model.Config_allocated = viper.GetFloat64("percent")
	model.Config_fee = viper.GetFloat64("fee")

	purchasing_amount := (model.Config_allocated / 100) * model.Config_amount // amount of coins to buy per interval

	purchasing_fee := (model.Config_fee / 100) * model.Config_amount //total price of coins to buy based on percentage of coins to buy

	purchasing_amount_perbuy := purchasing_amount / float64(model.Config_buyintervals) // amount of coins to buy per interval

	purchasing_fee_perbuy := purchasing_fee / float64(model.Config_buyintervals) //total price of coins to buy based on percentage of coins to buy

	buying_pricing := purchasing_amount_perbuy + purchasing_fee_perbuy //total price of coins to buy based on percentage of coins to buy

	model.Priceperbuy = strconv.FormatFloat(buying_pricing, 'f', 2, 64) //convert float to string

	fmt.Println("Price per buy:", model.Priceperbuy)

	viper.SetConfigFile("./config.yaml")
	viper.ReadInConfig()

	client := binance.NewClient(viper.GetString("api"), viper.GetString("secretkey"))
	order, err := client.NewCreateOrderService().Symbol(model.Coinid).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(model.Priceperbuy).Do(context.Background())
	if err != nil {
		fmt.Println("Order Failed", err)
		return
	}
	fmt.Println(order)
}

func crap() {
	model.Coinid = viper.GetString("coins") + "USDT"
	model.Config_buyintervals = viper.GetInt("buyintervals")
	model.Config_amount = viper.GetFloat64("amount")
	model.Config_allocated = viper.GetFloat64("percent")
	model.Config_fee = viper.GetFloat64("fee")

	purchasing_amount := (model.Config_allocated / 100) * model.Config_amount // amount of coins to buy per interval

	purchasing_fee := (model.Config_fee / 100) * model.Config_amount //total price of coins to buy based on percentage of coins to buy

	purchasing_amount_perbuy := purchasing_amount / float64(model.Config_buyintervals) // amount of coins to buy per interval

	purchasing_fee_perbuy := purchasing_fee / float64(model.Config_buyintervals) //total price of coins to buy based on percentage of coins to buy

	buying_pricing := purchasing_amount_perbuy + purchasing_fee_perbuy //total price of coins to buy based on percentage of coins to buy

	model.Priceperbuy = strconv.FormatFloat(buying_pricing, 'f', 2, 64) //convert float to string

	fmt.Println("Price per buy:", model.Priceperbuy)

	client := binance.NewClient(viper.GetString("api"), viper.GetString("secretkey"))

	err := client.NewCreateOrderService().Symbol(model.Coinid).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(model.Priceperbuy).Test(context.Background())
	if err != nil {
		fmt.Println("Test Order Failed", err)
		return
	}
}
