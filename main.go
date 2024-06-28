package main

import (
	"github.com/dkpcb/step1/model"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	sqlDB := model.DBConnection()
	defer sqlDB.Close()

	model.InsertCSVData("csv/trade_history.csv", "trade_history")
	model.InsertCSVData("csv/reference_prices.csv", "reference_prices")
}
