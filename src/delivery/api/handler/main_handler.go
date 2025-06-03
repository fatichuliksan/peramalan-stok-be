package handler

import (
	"peramalan-stok-be/src/helper/response"
	"peramalan-stok-be/src/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// PingHandler ...
type MainHandler struct {
	Response response.Interface
	DB       *gorm.DB
}

// Get ...
func (t *MainHandler) GetWarehouses(c echo.Context) error {
	query := c.QueryParam("query")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit == 0 {
		limit = 10
	}

	var warehouses []model.Warehouse
	db := t.DB.Debug().Model(&model.Warehouse{})

	if query != "" {
		db.Where("(warehouse_name ILIKE ? or warehouse_code ILIKE ? )", "'%"+query+"%'", "'%"+query+"%'")
	}

	// db.Limit(limit).Find(&warehouses)
	db.Find(&warehouses)

	return t.Response.SendSuccess(c, "", warehouses)
}

func (t *MainHandler) GetItems(c echo.Context) error {
	query := c.QueryParam("query")
	warehouse_code := c.QueryParam("warehouse_code")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit == 0 {
		limit = 10
	}

	var items []model.Item
	db := t.DB.Model(&model.Item{})

	if query != "" {
		db.Where("(item_code ILIKE ? or item_name ILIKE ? )", "'%"+query+"%'", "'%"+query+"%'")
	}

	if warehouse_code != "" {
		db.Where("warehouse_code = ?", warehouse_code)
	}

	db.Find(&items)

	return t.Response.SendSuccess(c, "", items)
}
