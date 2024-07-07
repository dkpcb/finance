package main

import (
	"github.com/dkpcb/finatext/data"
	"github.com/dkpcb/finatext/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {

	db, err := data.DBConnection()
	if err != nil {
		panic("failed to connect to database")
	}

	data.InitializeDatabase(db)

	data.InsertCSVData("csv/trade_history.csv", "trade_history")
	data.InsertCSVData("csv/reference_prices.csv", "reference_prices")

	e := echo.New()
	router.SetRouter(e, db)

}
