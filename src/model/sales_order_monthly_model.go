package model

type SalesOrderMonthly struct {
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
	ItemCode      string `json:"item_code"`
	ItemName      string `json:"item_name"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	Qty           int    `json:"qty"`
}

func (SalesOrderMonthly) TableName() string {
	return "sales_order_monthly"
}
