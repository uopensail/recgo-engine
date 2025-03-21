package table

import (
	"encoding/json"
	"time"

	"github.com/uopensail/recgo-engine/model/utils"
)

type InvertInexRecallMeta struct {
	Resource          string   `json:"resource" toml:"resource"`
	UserFeatureFields []string `json:"user_feature_fields" toml:"user_feature_fields"`
	EachMaxCol        int      `json:"each_max_col" toml:"each_max_col"`
	TopK              int      `json:"top_k" toml:"top_k"`
}

type RecallEntityMeta struct {
	EntityMeta   `json:",inline" toml:",inline" gorm:"embedded"`
	Condition    string `json:"condition" toml:"condition" gorm:"column:condition"`
	SortKey      string `json:"sort_key" toml:"sort_key" gorm:"column:sort_key"`
	PluginParams XJSON  `json:"plugin_params" toml:"plugin_params" gorm:"column:plugin_params"`
}

func (c *RecallEntityMeta) ParseInvertInexRecallMeta() InvertInexRecallMeta {
	invertInexRecallMeta := InvertInexRecallMeta{}
	json.Unmarshal([]byte(c.PluginParams), &invertInexRecallMeta)
	return invertInexRecallMeta
}

// 召回组计算实体
type RecallGroupEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	RecallEntities utils.IntSlice `json:"recall_entities" toml:"recall_entities" gorm:"column:recall_entities"` //一份引用
	//TODO: 每个子召回的权重排序
	EntityWeights map[int]float32
}

// redis的配置
type BaseResourceMeta struct {
	ID         int       `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	PluginName string    `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"`
	Name       string    `json:"name" toml:"name" gorm:"column:name"`
	UpdateTime time.Time `json:"update_time" toml:"update_time" gorm:"column:update_time;autoUpdateTime"`

	Source XJSON `json:"source" toml:"source" gorm:"column:source"`
}

type RedisResourceConfig struct {
	URL          string            `json:"url" toml:"url"`
	MinIdleConns int               `json:"min_idle_conns" toml:"min_idle_conns"`
	Timeout      int               `json:"timeout" toml:"timeout"`
	Params       map[string]string `json:"params" toml:"params"`
}

type FileResourceConfig struct {
	Location string            `json:"location" toml:"location"`
	Params   map[string]string `json:"params" toml:"params"`
}

type RecallResourceMeta struct {
	BaseResourceMeta    `json:",inline" toml:",inline" gorm:"embedded"`
	FileResourceConfig  `json:"file" toml:"file"`
	RedisResourceConfig `json:"redis" toml:"redis"`
}

func (c *RecallResourceMeta) ParseFileSource() {
	json.Unmarshal([]byte(c.Source), &c.FileResourceConfig)
}

func (c *RecallResourceMeta) ParseRedisSource() {
	json.Unmarshal([]byte(c.Source), &c.RedisResourceConfig)
}

func (c RecallResourceMeta) GetID() int {
	return c.ID
}

func (c RecallResourceMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}
