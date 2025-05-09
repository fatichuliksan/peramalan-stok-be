package handler

import (
	"math"
	"peramalan-stok-be/src/helper/response"
	"peramalan-stok-be/src/model"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// PingHandler ...
type HistoryHandler struct {
	Response response.Interface
	DB       *gorm.DB
}

func (t *HistoryHandler) GetHistorySalesOrder(c echo.Context) error {
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

	db := t.DB.Model(&model.SalesOrder{}).
		Select("to_char(posting_date,'YYYY-MM-DD') as posting_date",
			"warehouse_code", "warehouse_name", "sales_code", "sales_name", "customer_code", "customer_name", "item_code", "item_name", "item_unit", "quantity")

	if req.WarehouseCode != "" {
		db = db.Where("warehouse_code = ?", req.WarehouseCode)
	}

	if req.ItemCode != "" {
		db = db.Where("item_code = ?", req.ItemCode)
	}

	if req.DateStart != "" && req.DateEnd != "" {
		db = db.Where("posting_date >= ? AND posting_date <= ?", req.DateStart, req.DateEnd)
	}

	dbCountAll := db

	dbCountAll.Count(&totalRecord)

	if req.Search != "" {
		listSearch := []string{}
		for i := 0; i < 8; i++ {
			listSearch = append(listSearch, "%"+req.Search+"%")
		}

		db = db.Where(`(warehouse_code ILIKE ? 
		or warehouse_name ILIKE ?  
		or sales_code ILIKE ?  
		or sales_name ILIKE ?  
		or customer_code ILIKE ?  
		or customer_name ILIKE ?  
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
		db = db.Order("warehouse_code, posting_date")
	} else {
		db = db.Order(req.Order + " " + req.Sort)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 0 {
		totalPages = 0
	}
	offset := (limit * int64(req.Page)) - limit

	data := []model.SalesOrder{}

	db.Debug().Offset(int(offset)).Limit(int(limit)).Find(&data)

	return t.Response.SendSuccess(c, "", map[string]interface{}{
		"records":               data,
		"total_record":          totalRecord,
		"total_record_per_page": limit,
		"total_record_search":   totalSearch,
		"total_page":            totalPages,
	})
}

func (t *HistoryHandler) GetHistorySalesOrderMonthly(c echo.Context) error {
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

	db := t.DB.Debug().Model(&model.SalesOrderMonthly{})

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

		db = db.Where("year >= ? and month >= ?", dateStart.Year(), int(dateStart.Month()))
		db = db.Where("year <= ? and month <= ?", dateEnd.Year(), int(dateEnd.Month()))

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
		db = db.Order("warehouse_code, posting_date")
	} else {
		db = db.Order(req.Order + " " + req.Sort)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 0 {
		totalPages = 0
	}
	offset := (limit * int64(req.Page)) - limit

	data := []model.SalesOrderMonthly{}

	db.Offset(int(offset)).Limit(int(limit)).Find(&data)

	return t.Response.SendSuccess(c, "", map[string]interface{}{
		"records":               data,
		"total_record":          totalRecord,
		"total_record_per_page": limit,
		"total_record_search":   totalSearch,
		"total_page":            totalPages,
	})
}
