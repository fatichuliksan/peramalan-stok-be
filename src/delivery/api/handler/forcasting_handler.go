package handler

import (
	"peramalan-stok-be/src/helper/response"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ForecastingHandler struct {
	Response response.Interface
	DB       *gorm.DB
}

func (t *ForecastingHandler) PostGenerate(c echo.Context) error {
	response := []map[string]interface{}{}

	// Get the request body
	type Req struct {
		WarehouseCode string  `json:"warehouse_code"`
		ItemCode      string  `json:"item_code"`
		DateStart     string  `json:"date_start"`
		DateEnd       string  `json:"date_end"`
		Alpha         float64 `json:"alpha"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.WarehouseCode == "" {
		return t.Response.SendBadRequest(c, "warehouse is required", nil)
	}

	if req.ItemCode == "" {
		return t.Response.SendBadRequest(c, "item is required", nil)
	}
	if req.DateStart == "" {
		return t.Response.SendBadRequest(c, "period start is required", nil)
	}
	if req.DateEnd == "" {
		return t.Response.SendBadRequest(c, "period end is required", nil)
	}
	if req.Alpha == 0 {
		return t.Response.SendBadRequest(c, "alpha is required", nil)
	}
	if req.Alpha < 0 || req.Alpha > 1 {
		return t.Response.SendBadRequest(c, "alpha must be between 0 and 1", nil)
	}

	dateStart, err := time.Parse("2006-01-02", req.DateStart)
	if err != nil {
		return t.Response.SendBadRequest(c, "invalid date start", nil)

	}

	dateEnd, err := time.Parse("2006-01-02", req.DateEnd)
	if err != nil {
		return t.Response.SendBadRequest(c, "invalid date end", nil)
	}

	if dateStart.After(dateEnd) {
		return t.Response.SendBadRequest(c, "date start must be before date end", nil)
	}

	startYear := dateStart.Year()
	startMonth := int(dateStart.Month())
	endYear := dateEnd.Year()
	endMonth := int(dateEnd.Month())

	err = t.DB.Debug().Raw("select * from generate(?,?,?,?,?,?,?)", req.WarehouseCode, req.ItemCode, req.Alpha, startYear, startMonth, endYear, endMonth).Scan(&response).Error
	if err != nil {
		return t.Response.SendError(c, "error generate data", nil)
	}

	if len(response) == 0 {
		return t.Response.SendError(c, "no data generated, please check actual data in choosed period", nil)
	}

	return t.Response.SendSuccess(c, "data generated", response)

}
