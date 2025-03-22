package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DefaultAttribute definition
type DefaultAttribute struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedBy *uint          `json:"-"`
	UpdatedBy *uint          `json:"-"`
	DeletedBy *uint          `json:"-"`
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type TimestampAttribute struct {
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UserstampAttribute struct {
	CreatedBy *uint `json:"-"`
	UpdatedBy *uint `json:"-"`
	DeletedBy *uint `json:"-"`
}

type MapStringInterface map[string]interface{}

func (msi MapStringInterface) Value() (driver.Value, error) {
	valueString, err := json.Marshal(msi)
	return string(valueString), err
}

func (msi *MapStringInterface) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &msi); err != nil {
		return err
	}
	return nil
}

type ArrayString []string

func (as ArrayString) Value() (driver.Value, error) {
	valueString, err := json.Marshal(as)
	return string(valueString), err
}

func (as *ArrayString) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &as); err != nil {
		return err
	}
	return nil
}

type ArrayMapStringInterface []map[string]interface{}

func (amsi ArrayMapStringInterface) Value() (driver.Value, error) {
	valueString, err := json.Marshal(amsi)
	return string(valueString), err
}

func (amsi *ArrayMapStringInterface) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &amsi); err != nil {
		return err
	}
	return nil
}
