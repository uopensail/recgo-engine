package table

import (
	"time"
)

type PoolMeta struct {
	ID            int            `json:"id" toml:"ID" gorm:"column:id;primaryKey"`
	Name          string         `json:"name" toml:"name" gorm:"column:name"`
	PrimaryField  string         `json:"primary_field" toml:"primary_field" gorm:"column:primary_field"`
	Location      string         `json:"location" toml:"location" gorm:"column:location"`
	UpdateTime    time.Time      `json:"update_time" toml:"update_time" gorm:"column:update_time"`
	FieldDataType map[string]int `json:"field_data_type" toml:"field_data_type"`
}

func (c PoolMeta) GetID() int {
	return c.ID
}

func (c PoolMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}
