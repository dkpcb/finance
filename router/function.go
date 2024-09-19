package router

import (
	"log"
	"net/http"
	"time"

	"github.com/dkpcb/finance/data"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AssetResponse struct {
	Date         string `json:"date"`
	CurrentValue int    `json:"current_value"`
	CurrentPL    int    `json:"current_pl"`
}

type Handler struct {
	DB *gorm.DB
}

const unitsPerFund = 10000

func (h *Handler) GetTradesHandler(c echo.Context) error {
	userID := c.Param("user_id")

	var count int64
	query := "SELECT COUNT(*) FROM trade_histories WHERE user_id = ?"
	result := h.DB.Raw(query, userID).Scan(&count)
	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}

	return c.JSON(http.StatusOK, map[string]int64{"count": count})
}

// func (h *Handler) GetAssetsHandler(c echo.Context) error {
// 	userID := c.Param("user_id")

// 	currentValue, currentPL, err := data.CalculateAssetsAndPL(h.DB, userID)
// 	if err != nil {
// 		log.Printf("Error calculating assets and PL: %v", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
// 	}

// 	response := AssetResponse{
// 		Date:         time.Now().Format("2006-01-02"),
// 		CurrentValue: currentValue,
// 		CurrentPL:    currentPL,
// 	}

// 	return c.JSON(http.StatusOK, response)
// }

// func (h *Handler) GetAssetsdataHandler(c echo.Context) error {
// 	userID := c.Param("user_id")
// 	date := c.QueryParam("date")

// 	log.Printf("Received request - userID: %s, date: %s\n", userID, date)

// 	currentValue, currentPL, err := data.CalculateAssetsAndPL_date(h.DB, userID, date)
// 	if err != nil {
// 		log.Printf("Error calculating assets and PL: %v", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}

// 	response := AssetResponse{
// 		Date:         date,
// 		CurrentValue: currentValue,
// 		CurrentPL:    currentPL,
// 	}

// 	return c.JSON(http.StatusOK, response)
// }

func (h *Handler) GetAssetsHandler(c echo.Context) error {
	userID := c.Param("user_id")
	date := c.QueryParam("date")

	var currentValue, currentPL int
	var err error

	if date != "" {
		log.Printf("Received request - userID: %s, date: %s\n", userID, date)
		currentValue, currentPL, err = data.CalculateAssetsAndPL_date(h.DB, userID, date)
	} else {
		currentValue, currentPL, err = data.CalculateAssetsAndPL(h.DB, userID)
		date = time.Now().Format("2006-01-02")
	}

	if err != nil {
		log.Printf("Error calculating assets and PL: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := AssetResponse{
		Date:         date,
		CurrentValue: currentValue,
		CurrentPL:    currentPL,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetAssetsByYearHandler(c echo.Context) error {
	userID := c.Param("user_id")
	currentDate := time.Now().Format("2006-01-02")

	log.Printf("Received request - userID: %s, currentDate: %s\n", userID, currentDate)

	assetsByYear, err := data.CalculateAssetsAndPLByYear(h.DB, userID, currentDate)
	if err != nil {
		log.Printf("Error calculating assets and PL by year: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := data.AssetResponseByYear{
		Date:   currentDate,
		Assets: assetsByYear,
	}

	return c.JSON(http.StatusOK, response)
}

// func calculateAssetsAndPL(db *gorm.DB, userID string) (int, int, error) {
// 	queryTrade := "SELECT fund_id, quantity FROM trade_histories WHERE user_id = ?"
// 	queryPrice := "SELECT reference_price FROM reference_prices WHERE fund_id = ? ORDER BY reference_price_date DESC LIMIT 1"
// 	queryBuyPrice := "SELECT quantity, trade_date FROM trade_histories WHERE user_id = ? AND fund_id = ?"

// 	rows, err := db.Raw(queryTrade, userID).Rows()
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	defer rows.Close()

// 	var totalCurrentValue int
// 	var totalBuyPrice int

// 	for rows.Next() {
// 		var fundID string
// 		var quantity int
// 		if err := rows.Scan(&fundID, &quantity); err != nil {
// 			return 0, 0, err
// 		}

// 		var referencePrice int
// 		if err := db.Raw(queryPrice, fundID).Scan(&referencePrice).Error; err != nil {
// 			return 0, 0, err
// 		}

// 		currentValue := (referencePrice * quantity) / unitsPerFund
// 		totalCurrentValue += currentValue

// 		buyRows, err := db.Raw(queryBuyPrice, userID, fundID).Rows()
// 		if err != nil {
// 			return 0, 0, err
// 		}
// 		defer buyRows.Close()

// 		var totalFundBuyPrice int
// 		for buyRows.Next() {
// 			var buyQuantity int
// 			var tradeDate time.Time
// 			if err := buyRows.Scan(&buyQuantity, &tradeDate); err != nil {
// 				return 0, 0, err
// 			}

// 			var buyReferencePrice int
// 			if err := db.Raw(queryPrice, fundID).Scan(&buyReferencePrice).Error; err != nil {
// 				return 0, 0, err
// 			}

// 			buyPrice := (buyReferencePrice * buyQuantity) / unitsPerFund
// 			totalFundBuyPrice += buyPrice
// 		}

// 		totalBuyPrice += totalFundBuyPrice
// 	}

// 	currentPL := totalCurrentValue - totalBuyPrice
// 	return totalCurrentValue, currentPL, nil
// }
