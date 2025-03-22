package model

type SalesOrder struct {
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
	SalesCode     string `json:"sales_code"`
	SalesName     string `json:"sales_name"`
	CustomerCode  string `json:"customer_code"`
	CustomerName  string `json:"customer_name"`
	ItemCode      string `json:"item_code"`
	ItemName      string `json:"item_name"`
	PostingDate   string `json:"posting_date"`
	ItemUnit      string `json:"item_unit"`
	Quantity      int    `json:"quantity"`
}

func (SalesOrder) TableName() string {
	return "sales_orders"
}
