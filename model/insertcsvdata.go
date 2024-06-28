package model

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func InsertCSVData(filePath string, tableName string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	switch tableName {
	case "trade_history":
		for _, record := range records[1:] {
			quantity, _ := strconv.Atoi(record[2])
			tradeHistory := TradeHistory{
				UserID:    record[0],
				FundID:    record[1],
				Quantity:  quantity,
				TradeDate: record[3],
			}
			db.Create(&tradeHistory)
		}
	case "reference_prices":
		for _, record := range records[1:] {
			referencePrice, _ := strconv.Atoi(record[2])
			referencePriceRecord := ReferencePrice{
				FundID:             record[0],
				ReferencePriceDate: record[1],
				ReferencePrice:     referencePrice,
			}
			db.Create(&referencePriceRecord)
		}
	default:
		log.Fatalf("Unknown table name: %s", tableName)
	}

	log.Printf("Data from %s inserted into %s table successfully.", filePath, tableName)
}
