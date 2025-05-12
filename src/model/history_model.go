package model

import (
	"time"

	"gorm.io/gorm"

	"gorm.io/datatypes"
)

type History struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
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
	CreatedAt     *time.Time     `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

func (History) TableName() string {
	return "histories"
}
