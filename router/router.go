package router

import (
	"os"

	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	_ "net/http"

	"github.com/labstack/echo/v4"
)

func SetRouter(e *echo.Echo, db *gorm.DB) error {

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${host} ${method} ${uri} ${status} ${header}\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	handler := &Handler{DB: db}

	e.GET("/:user_id/trades", handler.GetTradesHandler)
	e.GET("/:user_id/assets", handler.GetAssetsHandler)
	e.GET("/:user_id/assets/byYear", handler.GetAssetsByYearHandler)

	err := e.Start(":8080")
	return err
}
