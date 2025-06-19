package handler

import (
	"log"
	"math"
	"peramalan-stok-be/src/helper/response"
	"peramalan-stok-be/src/model"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ForecastingHandler struct {
	Response response.Interface
	DB       *gorm.DB
}

func (t *ForecastingHandler) PostGenerate(c echo.Context) error {

	// Get the request body
	type Req struct {
		WarehouseCode string  `json:"warehouse_code"`
		ItemCode      string  `json:"item_code"`
		DateStart     string  `json:"date_start"`
		DateEnd       string  `json:"date_end"`
		Alpha         float64 `json:"alpha"`
		ForcastPeriod int     `json:"forcast_period"`
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

	// validate minimal range of 4 months, from month range of dateStart and dateEnd
	if dateStart.Year() == dateEnd.Year() && int(dateEnd.Month())-int(dateStart.Month()) < 4 {
		return t.Response.SendBadRequest(c, "period must be at least 4 months", nil)
	}

	startYear := dateStart.Year()
	startMonth := int(dateStart.Month())
	endYear := dateEnd.Year()
	endMonth := int(dateEnd.Month())

	generateLines := []map[string]interface{}{}
	err = t.DB.Debug().Raw("select * from generate(?,?,?,?,?,?,?)", req.WarehouseCode, req.ItemCode, req.Alpha, startYear, startMonth, endYear, endMonth).Scan(&generateLines).Error
	if err != nil {
		return t.Response.SendError(c, "error generate data", nil)
	}

	if len(generateLines) == 0 {
		return t.Response.SendError(c, "no data generated, please check actual data in choosed period", nil)
	}

	// check generateLines for no missing months
	isValid := true
	nextYear := int32(0)
	nextMonth := int32(0)
	for index, gl := range generateLines {
		if index == 0 {
			if gl["month"].(int32) == 12 {
				nextYear = gl["year"].(int32) + 1
				nextMonth = 1
			} else {
				nextYear = gl["year"].(int32)
				nextMonth = gl["month"].(int32) + 1
			}
			continue
		}

		log.Println("checking month at index", index, "year", gl["year"], "month", gl["month"], "nextYear", nextYear, "nextMonth", nextMonth)

		if gl["year"].(int32) != nextYear && gl["month"].(int32) != nextMonth {
			isValid = false
			return t.Response.SendBadRequest(c, "generate has missing data at year: "+string(gl["year"].(int32))+", month: "+string(gl["month"].(int32)), nil)
		}

		if gl["month"].(int32) == 12 {
			nextYear = gl["year"].(int32) + 1
			nextMonth = 1
		} else {
			nextYear = gl["year"].(int32)
			nextMonth = gl["month"].(int32) + 1
		}

	}

	mape := generateLines[len(generateLines)-1]["e_sig"].(float64)
	mapeCriteria := "Invalid"
	if len(generateLines) > 3 && isValid {
		if mape >= 0 && mape < 10 {
			mapeCriteria = "Highly Accurate"
		} else if mape >= 10 && mape < 20 {
			mapeCriteria = "Accurate"
		} else if mape >= 20 && mape < 50 {
			mapeCriteria = "Less Accurate"
		} else if mape >= 50 {
			mapeCriteria = "Inaccurate"
		}
	}

	req.ForcastPeriod = 1

	var forecastLines = []map[string]interface{}{}
	aa := generateLines[len(generateLines)-1]["a"].(float64)
	bb := generateLines[len(generateLines)-1]["b"].(float64)
	cc := generateLines[len(generateLines)-1]["c"].(float64)
	for i := 0; i < req.ForcastPeriod; i++ {
		forecastLines = append(forecastLines, map[string]interface{}{
			"t": i + 1,
			"a": aa,
			"b": bb,
			"c": cc,
			"f": aa + bb*float64(i+1) + 0.5*cc*float64(i+1),
		})
	}

	return t.Response.SendSuccess(c, "data generated", map[string]interface{}{
		"generate_lines": generateLines,
		"forcast_lines":  forecastLines,
		"forcast_period": req.ForcastPeriod,
		"mape":           mape,
		"mape_criteria":  mapeCriteria,
	})
}

func (t *ForecastingHandler) PostHistory(c echo.Context) error {

	response := []map[string]interface{}{}

	// Get the request body
	type Req struct {
		WarehouseCode string         `json:"warehouse_code"`
		WarehouseName string         `json:"warehouse_name"`
		ItemCode      string         `json:"item_code"`
		ItemName      string         `json:"item_name"`
		Alpha         float32        `json:"alpha"`
		DateStart     string         `json:"date_start"`
		DateEnd       string         `json:"date_end"`
		GenerateLines datatypes.JSON `json:"generate_lines"`
		ForcastLines  datatypes.JSON `json:"forcast_lines"`
		ForcastPeriod int            `json:"forcast_period"`
		Mape          float32        `json:"mape"`
		MapeCriteria  string         `json:"mape_criteria"`
	}

	var req Req
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	if req.WarehouseCode == "" {
		return t.Response.SendBadRequest(c, "warehouse is required", nil)
	}
	if req.ItemCode == "" {
		return t.Response.SendBadRequest(c, "item is required", nil)
	}
	if req.Alpha <= 0 || req.Alpha > 1 {
		return t.Response.SendBadRequest(c, "alpha must be between 0 and 1", nil)
	}

	if req.DateStart == "" {
		return t.Response.SendBadRequest(c, "period start is required", nil)
	}

	if req.DateEnd == "" {
		return t.Response.SendBadRequest(c, "period end is required", nil)
	}

	// check if data aleady exist
	var count int64
	err = t.DB.Model(&model.History{}).Where("warehouse_code = ? and item_code = ? and date_start = ? and date_end = ? and alpha =?", req.WarehouseCode, req.ItemCode, req.DateStart, req.DateEnd, req.Alpha).Count(&count).Error
	if err != nil {
		return t.Response.SendError(c, "error check history", nil)
	}
	if count > 0 {
		return t.Response.SendBadRequest(c, "data already exist", nil)
	}

	if len(req.GenerateLines) == 0 {
		return t.Response.SendBadRequest(c, "generate data is required", nil)
	}

	var histories = model.History{}
	histories.WarehouseCode = req.WarehouseCode
	histories.WarehouseName = req.WarehouseName
	histories.ItemCode = req.ItemCode
	histories.ItemName = req.ItemName
	histories.Alpha = req.Alpha
	histories.DateStart = req.DateStart
	histories.DateEnd = req.DateEnd
	histories.GenerateLines = req.GenerateLines
	histories.ForcastLines = req.ForcastLines
	histories.ForcastPeriod = req.ForcastPeriod
	histories.Mape = req.Mape
	histories.MapeCriteria = req.MapeCriteria

	log.Println("histories", histories)

	err = t.DB.Debug().Create(&histories).Error
	if err != nil {
		return t.Response.SendError(c, "error create history", nil)
	}

	return t.Response.SendSuccess(c, "success to save generated data", response)

}

func (t *ForecastingHandler) GetHistory(c echo.Context) error {
	type Req struct {
		Search        string `query:"search"`
		Length        int    `query:"length"`
		Page          int    `query:"page"`
		Sort          string `query:"sort"`
		Order         string `query:"order"`
		WarehouseCode string `query:"warehouse_code"`
		ItemCode      string `query:"item_code"`
		DateStart     string `query:"date_start"`
		DateEnd       string `query:"date_end"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	var totalRecord int64
	var limit int64
	var totalSearch int64
	var total int64

	db := t.DB.Debug().Model(&model.History{})

	if req.WarehouseCode != "" {
		db = db.Where("warehouse_code = ?", req.WarehouseCode)
	}

	if req.ItemCode != "" {
		db = db.Where("item_code = ?", req.ItemCode)
	}

	if req.DateStart != "" && req.DateEnd != "" {
		dateStart, err := time.Parse("2006-01-02", req.DateStart)
		if err != nil {
			return t.Response.SendError(c, err.Error(), nil)
		}

		dateEnd, err := time.Parse("2006-01-02", req.DateEnd)
		if err != nil {
			return t.Response.SendError(c, err.Error(), nil)
		}

		db = db.Where("date_start >= ? and date_end <= ?", dateStart.Format("2006-01-02"), dateEnd.Format("2006-01-02"))

	}

	dbCountAll := db

	dbCountAll.Count(&totalRecord)

	if req.Search != "" {
		listSearch := []string{}
		for i := 0; i < 4; i++ {
			listSearch = append(listSearch, "%"+req.Search+"%")
		}

		db = db.Where(`(warehouse_code ILIKE ? 
		or warehouse_name ILIKE ?   
		or item_code ILIKE ? 
		or item_name ILIKE ? )`, listSearch)

		dbCountSearch := db
		dbCountSearch.Count(&totalSearch)
	} else {
		total = totalRecord
	}

	if req.Length == 0 {
		limit = total
	} else {
		limit = int64(req.Length)
	}

	if req.Order == "" {
		db = db.Order("warehouse_code, alpha")
	} else {
		db = db.Order(req.Order + " " + req.Sort)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 0 {
		totalPages = 0
	}
	offset := (limit * int64(req.Page)) - limit

	data := []model.History{}

	db.Offset(int(offset)).Limit(int(limit)).Find(&data)

	return t.Response.SendSuccess(c, "", map[string]interface{}{
		"records":               data,
		"total_record":          totalRecord,
		"total_record_per_page": limit,
		"total_record_search":   totalSearch,
		"total_page":            totalPages,
	})
}

func (t *ForecastingHandler) DeleteHistory(c echo.Context) error {

	response := []map[string]interface{}{}

	// Get the request body
	type Req struct {
		ID uint64 `json:"id" query:"id"`
	}

	var req Req
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	if req.ID == 0 {
		return t.Response.SendBadRequest(c, "id is required", nil)
	}

	var histories = model.History{}

	err = t.DB.Debug().Delete(&histories, req.ID).Error
	if err != nil {
		return t.Response.SendError(c, "error delete history", nil)
	}

	return t.Response.SendSuccess(c, "success to delete history data", response)

}
