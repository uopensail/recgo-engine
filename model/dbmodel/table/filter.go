package table

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/uopensail/recgo-engine/model/utils"
)

type Entitiestatus int

const (
	OfflineEntiyStatus Entitiestatus = 0
	NormalEntiyStatus  Entitiestatus = 1 // 默认全量计算实体
	ExpingEntiyStatus  Entitiestatus = 2 //实验阶段
)

// 定义一个策略
type FilterEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	Condition string `json:"condition" toml:"condition" gorm:"column:condition"`
	MaxCount  int    `json:"max_count" toml:"max_count" gorm:"column:max_count"`
	Format    string `json:"format" toml:"format" gorm:"column:format"`

	SourceID int `json:"source_id" toml:"source_id" gorm:"column:source_id"`
}

func (cfg FilterEntityMeta) GetID() int {
	return cfg.ID
}
func (cfg FilterEntityMeta) GetUpdateTime() int64 {
	return cfg.UpdateTime.Unix()
}

type RedisConfigure struct {
	URL          string            `json:"url" toml:"url"`
	MinIdleConns int               `json:"min_idle_conns" toml:"min_idle_conns"`
	Timeout      int               `json:"timeout" toml:"timeout"`
	Params       map[string]string `json:"params" toml:"params"`
}

func (s *RedisConfigure) Scan(val interface{}) error {
	switch val := val.(type) {
	case string:
		return json.Unmarshal([]byte(val), s)
	case []byte:
		return json.Unmarshal(val, s)
	default:
		return errors.New("not support")
	}

}

func (s RedisConfigure) Value() (driver.Value, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}
func (s RedisConfigure) GormDataType() string {
	return "json"
}

type FilterResourceMeta struct {
	ID         int       `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	PluginName string    `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"` //插件模式，相当于类名
	Name       string    `json:"name" toml:"name" gorm:"column:name"`
	UpdateTime time.Time `json:"update_time" toml:"update_time" gorm:"column:update_time"`

	Redis RedisConfigure `json:"redis" toml:"redis" gorm:"column:redis"`
}

func (cfg FilterResourceMeta) GetID() int {
	return cfg.ID
}
func (cfg FilterResourceMeta) GetUpdateTime() int64 {
	return cfg.UpdateTime.Unix()
}

// 召回组计算实体
type FilterGroupEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	Timeout int `json:"timeout" toml:"timeout" gorm:"column:timeout"`

	FilterEntities utils.IntSlice `json:"filter_entities" toml:"filter_entities" gorm:"column:filter_entities"` //一份引用
}

func (cfg FilterGroupEntityMeta) GetID() int {
	return cfg.ID
}
func (cfg FilterGroupEntityMeta) GetUpdateTime() int64 {
	return cfg.UpdateTime.Unix()
}
