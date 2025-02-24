package table

import (
	"time"

	"github.com/uopensail/recgo-engine/model/utils"
)

type EntityMeta struct {
	ID         int             `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	UpdateTime time.Time       `json:"update_time" toml:"update_time" gorm:"column:update_time"`
	Name       string          `json:"name" toml:"name" gorm:"column:name"`
	Status     Entitiestatus   `json:"status,omitempty" toml:"status" gorm:"column:status"`
	PluginName string          `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"`
	ABLayerID  int             `json:"ab_layer_id" toml:"ab_layer_id" gorm:"column:ab_layer_id"` //绑定实验层
	Params     utils.StringMap `json:"params" toml:"params" gorm:"column:params"`
}
