package main

import (
	"fmt"
	"os"

	"github.com/dkpcb/finance/data"
	"github.com/dkpcb/finance/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "import" {
		runImport()
	} else {
		runImport()
		runServer()
	}
}

func runImport() {
	db, err := data.DBConnection()
	if err != nil {
		panic("failed to connect to database")
	}

	data.InitializeDatabase(db)

	data.InsertCSVData("csv/trade_history.csv", "trade_history")
	data.InsertCSVData("csv/reference_prices.csv", "reference_prices")

	fmt.Println("Data import completed")
}

func runServer() {
	db, err := data.DBConnection()
	if err != nil {
		panic("failed to connect to database")
	}

	e := echo.New()
	router.SetRouter(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}
