package model

type Warehouse struct {
	WarehouseCode string `gorm:"primaryKey" json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}
