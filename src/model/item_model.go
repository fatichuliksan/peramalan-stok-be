package model

type Item struct {
	ItemCode string `gorm:"primaryKey" json:"item_code"`
	ItemName string `json:"item_name"`
}

func (Item) TableName() string {
	return "items"
}
