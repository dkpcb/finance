package data

import (
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

type AssetResponse struct {
	Date         string `json:"date"`
	CurrentValue int    `json:"current_value"`
	CurrentPL    int    `json:"current_pl"`
}

type AssetByYear struct {
	Year         int `json:"year"`
	CurrentValue int `json:"current_value"`
	CurrentPL    int `json:"current_pl"`
}

type AssetResponseByYear struct {
	Date   string        `json:"date"`
	Assets []AssetByYear `json:"assets"`
}

type Handler struct {
	DB *gorm.DB
}

const unitsPerFund = 10000

func CalculateAssetsAndPL(db *gorm.DB, userID string) (int, int, error) {
	queryTrade := "SELECT fund_id, quantity FROM trade_histories WHERE user_id = ?"
	queryRecentPrice := "SELECT reference_price FROM reference_prices WHERE fund_id = ? ORDER BY reference_price_date DESC LIMIT 1"
	queryBuyHistory := "SELECT fund_id, trade_date, quantity FROM trade_histories WHERE user_id = ?"
	queryBuyPrice := "SELECT reference_price FROM reference_prices WHERE fund_id = ? AND reference_price_date = ?"

	//user_idを指定してtrade_history.csvからfund_id, quantityを取得
	rows, err := db.Raw(queryTrade, userID).Rows()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query trade_histories: %w", err)
	}
	defer rows.Close()

	var totalCurrentValue int

	//fund_idを用いてreference_prices.csvから最も最近のreference_priceを取得
	//reference_priceとquantityを用いてcurrentValueを計算
	//for文を使ってuser_idの中の全てのfund_idに対してtotalCurrentValue += currentValueを求める
	for rows.Next() {
		var fundID string
		var quantity int
		if err := rows.Scan(&fundID, &quantity); err != nil {
			return 0, 0, fmt.Errorf("failed to scan trade_histories row: %w", err)
		}

		var referencePrice int
		if err := db.Raw(queryRecentPrice, fundID).Scan(&referencePrice).Error; err != nil {
			return 0, 0, fmt.Errorf("failed to query recent reference_price for fund %s: %w", fundID, err)
		}

		currentValue := (referencePrice * quantity) / unitsPerFund
		totalCurrentValue += currentValue

		// デバッグ用出力：整数の切り捨てが正しく行われているか確認
		fmt.Printf("Current Value Calculation - fundID: %s, quantity: %d, referencePrice: %d, currentValue: %d, totalCurrentValue: %d\n", fundID, quantity, referencePrice, currentValue, totalCurrentValue)
	}

	// user_idを指定してtrade_history.csvからfund_id, trade_date, quantityを取得
	buyRows, err := db.Raw(queryBuyHistory, userID).Rows()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query buy history: %w", err)
	}
	defer buyRows.Close()

	var totalBuyPrice int

	//fund_id, trade_dateを用いてreference_priceを取得
	//trade_dateとreference_priceを用いてbuyPrice を計算
	//for文を使ってuser_idの中の全てのfund_idに対してtotalFundBuyPrice += buyPriceを求める
	for buyRows.Next() {
		var fundID string
		var tradeDateStr string
		var buyQuantity int
		if err := buyRows.Scan(&fundID, &tradeDateStr, &buyQuantity); err != nil {
			return 0, 0, fmt.Errorf("failed to scan buy history row: %w", err)
		}

		var buyReferencePrice int
		if err := db.Raw(queryBuyPrice, fundID, tradeDateStr).Scan(&buyReferencePrice).Error; err != nil {
			return 0, 0, fmt.Errorf("failed to query buy reference_price for fund %s on date %s: %w", fundID, tradeDateStr, err)
		}

		buyPrice := (buyReferencePrice * buyQuantity) / unitsPerFund
		totalBuyPrice += buyPrice

		// デバッグ用出力：整数の切り捨てが正しく行われているか確認
		fmt.Printf("Buy Price Calculation - fundID: %s, tradeDate: %s, buyQuantity: %d, buyReferencePrice: %d, buyPrice: %d, totalBuyPrice: %d\n", fundID, tradeDateStr, buyQuantity, buyReferencePrice, buyPrice, totalBuyPrice)
	}

	if err = rows.Err(); err != nil {
		return 0, 0, fmt.Errorf("error occurred during iteration of trade_histories rows: %w", err)
	}

	//最後にcurrentPL := totalCurrentValue - totalBuyPriceを求める
	currentPL := totalCurrentValue - totalBuyPrice

	// デバッグ用出力
	fmt.Printf("Final Calculation - totalCurrentValue: %d, totalBuyPrice: %d, currentPL: %d\n", totalCurrentValue, totalBuyPrice, currentPL)

	return totalCurrentValue, currentPL, nil
}

func CalculateAssetsAndPL_date(db *gorm.DB, userID string, date string) (int, int, error) {
	queryTrade := "SELECT fund_id, SUM(quantity) as total_quantity FROM trade_histories WHERE user_id = ? AND trade_date <= ? GROUP BY fund_id"
	queryReferencePrice := "SELECT reference_price FROM reference_prices WHERE fund_id = ? AND reference_price_date = ?"
	queryBuyHistory := "SELECT fund_id, trade_date, quantity FROM trade_histories WHERE user_id = ? AND trade_date <= ?"

	// 指定日付までのユーザーの全ての取引を取得
	rows, err := db.Raw(queryTrade, userID, date).Rows()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query trade_histories: %w", err)
	}
	defer rows.Close()

	var totalDateValue int

	for rows.Next() {
		var fundID string
		var quantity int
		if err := rows.Scan(&fundID, &quantity); err != nil {
			return 0, 0, fmt.Errorf("failed to scan trade_histories row: %w", err)
		}

		var referencePrice int
		if err := db.Raw(queryReferencePrice, fundID, date).Scan(&referencePrice).Error; err != nil {
			return 0, 0, fmt.Errorf("failed to query reference_price for fund %s on date %s: %w", fundID, date, err)
		}

		dateValue := (referencePrice * quantity) / unitsPerFund
		totalDateValue += dateValue

		fmt.Printf("Date Value Calculation - fundID: %s, quantity: %d, referencePrice: %d, dateValue: %d, totalDateValue: %d\n", fundID, quantity, referencePrice, dateValue, totalDateValue)
	}

	// 指定日付までの購入履歴を取得
	buyRows, err := db.Raw(queryBuyHistory, userID, date).Rows()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query buy history: %w", err)
	}
	defer buyRows.Close()

	var totalBuyPrice int

	for buyRows.Next() {
		var fundID string
		var tradeDateStr string
		var buyQuantity int
		if err := buyRows.Scan(&fundID, &tradeDateStr, &buyQuantity); err != nil {
			return 0, 0, fmt.Errorf("failed to scan buy history row: %w", err)
		}

		var buyReferencePrice int
		if err := db.Raw(queryReferencePrice, fundID, tradeDateStr).Scan(&buyReferencePrice).Error; err != nil {
			return 0, 0, fmt.Errorf("failed to query buy reference_price for fund %s on date %s: %w", fundID, tradeDateStr, err)
		}

		buyPrice := (buyReferencePrice * buyQuantity) / unitsPerFund
		totalBuyPrice += buyPrice

		fmt.Printf("Buy Price Calculation - fundID: %s, tradeDate: %s, buyQuantity: %d, buyReferencePrice: %d, buyPrice: %d, totalBuyPrice: %d\n", fundID, tradeDateStr, buyQuantity, buyReferencePrice, buyPrice, totalBuyPrice)
	}

	if err = rows.Err(); err != nil {
		return 0, 0, fmt.Errorf("error occurred during iteration of trade_histories rows: %w", err)
	}

	currentPL := totalDateValue - totalBuyPrice

	fmt.Printf("Final Calculation - totalDateValue: %d, totalBuyPrice: %d, currentPL: %d\n", totalDateValue, totalBuyPrice, currentPL)

	return totalDateValue, currentPL, nil
}

func CalculateAssetsAndPLByYear(db *gorm.DB, userID string, currentDate string) ([]AssetByYear, error) {
	queryTrade := "SELECT fund_id, quantity, YEAR(trade_date) as year FROM trade_histories WHERE user_id = ?"
	queryReferencePrice := "SELECT reference_price FROM reference_prices WHERE fund_id = ? ORDER BY reference_price_date DESC LIMIT 1"
	queryBuyHistory := "SELECT fund_id, trade_date, quantity, YEAR(trade_date) as year FROM trade_histories WHERE user_id = ?"
	queryBuyPrice := "SELECT reference_price FROM reference_prices WHERE fund_id = ? AND reference_price_date = ?"

	// デバッグ用出力
	log.Printf("Executing query: %s with user_id: %s", queryTrade, userID)

	rows, err := db.Raw(queryTrade, userID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query trade_histories: %w", err)
	}
	defer rows.Close()

	yearMap := make(map[int]map[string]int)

	for rows.Next() {
		var fundID string
		var quantity int
		var year int
		if err := rows.Scan(&fundID, &quantity, &year); err != nil {
			return nil, fmt.Errorf("failed to scan trade_histories row: %w", err)
		}

		// デバッグ用出力
		log.Printf("Trade row - fundID: %s, quantity: %d, year: %d", fundID, quantity, year)

		if _, exists := yearMap[year]; !exists {
			yearMap[year] = map[string]int{"current_value": 0, "buy_price": 0}
		}

		var referencePrice int
		if err := db.Raw(queryReferencePrice, fundID).Scan(&referencePrice).Error; err != nil {
			return nil, fmt.Errorf("failed to query reference_price for fund %s: %w", fundID, err)
		}

		// デバッグ用出力
		log.Printf("Reference price - fundID: %s, referencePrice: %d", fundID, referencePrice)

		currentValue := (referencePrice * quantity) / unitsPerFund
		yearMap[year]["current_value"] += currentValue

		// デバッグ用出力
		log.Printf("Current value calculation - year: %d, currentValue: %d, totalCurrentValue: %d", year, currentValue, yearMap[year]["current_value"])
	}

	buyRows, err := db.Raw(queryBuyHistory, userID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query buy history: %w", err)
	}
	defer buyRows.Close()

	for buyRows.Next() {
		var fundID string
		var tradeDateStr string // 修正 まず文字列として読み取る
		var buyQuantity int
		var year int
		if err := buyRows.Scan(&fundID, &tradeDateStr, &buyQuantity, &year); err != nil {
			return nil, fmt.Errorf("failed to scan buy history row: %w", err)
		}

		// 文字列からtime.Timeに変換
		tradeDate, err := time.Parse("2006-01-02", tradeDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse trade date: %w", err)
		}

		// デバッグ用出力
		log.Printf("Buy row - fundID: %s, tradeDate: %s, quantity: %d", fundID, tradeDate.Format("2006-01-02"), buyQuantity)

		var buyReferencePrice int
		if err := db.Raw(queryBuyPrice, fundID, tradeDate.Format("2006-01-02")).Scan(&buyReferencePrice).Error; err != nil {
			return nil, fmt.Errorf("failed to query buy reference_price for fund %s on date %s: %w", fundID, tradeDate.Format("2006-01-02"), err)
		}

		// デバッグ用出力
		log.Printf("Buy reference price - fundID: %s, tradeDate: %s, referencePrice: %d", fundID, tradeDate.Format("2006-01-02"), buyReferencePrice)

		buyPrice := (buyReferencePrice * buyQuantity) / unitsPerFund
		yearMap[year]["buy_price"] += buyPrice

		// デバッグ用出力
		log.Printf("Buy price calculation - year: %d, buyPrice: %d, totalBuyPrice: %d", year, buyPrice, yearMap[year]["buy_price"])
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during iteration of trade_histories rows: %w", err)
	}

	var assetsByYear []AssetByYear
	for year, values := range yearMap {
		currentPL := values["current_value"] - values["buy_price"]
		assetsByYear = append(assetsByYear, AssetByYear{
			Year:         year,
			CurrentValue: values["current_value"],
			CurrentPL:    currentPL,
		})

		// デバッグ用出力
		log.Printf("Yearly calculation - year: %d, currentValue: %d, buyPrice: %d, currentPL: %d", year, values["current_value"], values["buy_price"], currentPL)
	}

	// 年で降順にソート
	sort.Slice(assetsByYear, func(i, j int) bool {
		return assetsByYear[i].Year > assetsByYear[j].Year
	})

	return assetsByYear, nil
}
