package data

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/gorm"
)

const batchSize = 300 //2000

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

	var tradeHistories []TradeHistory
	var referencePrices []ReferencePrice
	cnt := 0

	err = db.Transaction(func(tx *gorm.DB) error {

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
				tradeHistories = append(tradeHistories, tradeHistory)
				cnt++

				if cnt >= batchSize {
					if err := tx.Create(&tradeHistories).Error; err != nil {
						return err
					}
					tradeHistories = nil
					cnt = 0
				}
			}
		case "reference_prices":
			for _, record := range records[1:] {
				referencePrice, _ := strconv.Atoi(record[1])
				referencePriceRecord := ReferencePrice{
					FundID:             record[0],
					ReferencePrice:     referencePrice,
					ReferencePriceDate: record[2],
				}
				referencePrices = append(referencePrices, referencePriceRecord)
				cnt++

				if cnt >= batchSize {
					if err := tx.Create(&referencePrices).Error; err != nil {
						return err
					}
					referencePrices = nil
					cnt = 0
				}
			}
		default:
			return fmt.Errorf("Unknown table name: %s", tableName)
		}

		// 残りのデータを挿入
		if len(tradeHistories) > 0 {
			if err := tx.Create(&tradeHistories).Error; err != nil {
				return err
			}
		}
		if len(referencePrices) > 0 {
			if err := tx.Create(&referencePrices).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	log.Printf("Data from %s inserted into %s table successfully.", filePath, tableName)
}

// package data

// import (
// 	"encoding/csv"
// 	"log"
// 	"os"
// 	"strconv"
// )

// func InsertCSVData(filePath string, tableName string) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Fatalf("Error opening file: %v", err)
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		log.Fatalf("Error reading CSV file: %v", err)
// 	}

// 	switch tableName {
// 	case "trade_history":
// 		for _, record := range records[1:] {
// 			quantity, _ := strconv.Atoi(record[2])
// 			tradeHistory := TradeHistory{
// 				UserID:    record[0],
// 				FundID:    record[1],
// 				Quantity:  quantity,
// 				TradeDate: record[3],
// 			}
// 			db.Create(&tradeHistory)
// 		}
// 	case "reference_prices":
// 		for _, record := range records[1:] {
// 			referencePrice, _ := strconv.Atoi(record[1])
// 			referencePriceRecord := ReferencePrice{
// 				FundID:             record[0],
// 				ReferencePrice:     referencePrice,
// 				ReferencePriceDate: record[2],
// 			}
// 			db.Create(&referencePriceRecord)
// 		}
// 	default:
// 		log.Fatalf("Unknown table name: %s", tableName)
// 	}

// 	log.Printf("Data from %s inserted into %s table successfully.", filePath, tableName)
// }
