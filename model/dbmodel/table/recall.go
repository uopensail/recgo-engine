package table

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uopensail/recgo-engine/model/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OrderByMeta struct {
	Field string `json:"field" toml:"field"`
	Desc  bool   `json:"desc" toml:"desc"`
}

type IndexColumnMeta struct {
	Column         string `json:"column" toml:"column"`
	FormatExecJson string `json:"format_exec_json" toml:"format_exec_json"`
}

type IndexMeta struct {
}

type FromMeta struct {
	Resource          string `json:"resource" toml:"resource"`
	KeyFormatExecJson string `json:"key_format_exec_json" toml:"key_format_exec_json"`
}

type ConditionMeta struct {
	StaticCondition  string `json:"static_condition" toml:"static_condition"`   //只包含物料特征的
	RuntimeCondition string `json:"runtime_condition" toml:"runtime_condition"` //需要实时计算的
	IndexMeta        `json:"index" toml:"index"`
}

type DSLMeta struct {
	Name          string `json:"name" toml:"name"`
	FromMeta      `json:"from" toml:"from"`
	ConditionMeta `json:"condition" toml:"condition"`
	OrderByMeta   `json:"orderby" toml:"orderby"`

	Filter string `json:"filter" toml:"filter"`
	Limit  int    `json:"limit" toml:"limit"`
}

func (s *DSLMeta) Scan(val interface{}) error {
	switch val := val.(type) {
	case string:
		return json.Unmarshal([]byte(val), s)
	case []byte:
		return json.Unmarshal(val, s)
	default:
		return errors.New("not support")
	}

}

func (s DSLMeta) Value() (driver.Value, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}
func (s DSLMeta) GormDataType() string {
	return "json"
}

type RecallEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	DSL     string `json:"dsl" toml:"dsl" gorm:"column:dsl"`
	DSLMeta `json:"dsl_json" toml:"dsl_json" gorm:"column:dsl_json"`
}

func (c RecallEntityMeta) GetID() int {
	return c.ID
}

func (c RecallEntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}

// 召回组计算实体
type RecallGroupEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	RecallEntities utils.IntSlice `json:"recall_entities" toml:"recall_entities" gorm:"column:recall_entities"` //一份引用
	//TODO: 每个子召回的权重排序
	EntityWeights map[int]float32
}

func (c RecallGroupEntityMeta) GetID() int {
	return c.ID
}

func (c RecallGroupEntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}

// redis的配置
type BaseResourceMeta struct {
	ID         int       `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	PluginName string    `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"`
	Name       string    `json:"name" toml:"name" gorm:"column:name"`
	UpdateTime time.Time `json:"update_time" toml:"update_time" gorm:"column:update_time;autoUpdateTime"`

	Source XJSON `json:"source" toml:"source" gorm:"column:source"`
}
type XJSON datatypes.JSON

// MarshalJSON returns m as the JSON encoding of m.
func (m XJSON) MarshalTOML() ([]byte, error) {
	s := string(m)
	return json.Marshal(s)
}

// UnmarshalJSON sets *m to a copy of data.
func (m *XJSON) UnmarshalText(data []byte) error {
	if m == nil {
		return errors.New("XJSON: UnmarshalTOML on nil pointer")
	}

	(*m) = append((*m)[0:0], data...)
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j XJSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *XJSON) Scan(value interface{}) error {
	if value == nil {
		*j = XJSON("null")
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage(bytes)
	*j = XJSON(result)
	return nil
}

// MarshalJSON to output non base64 encoded []byte
func (j XJSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (j *XJSON) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*j = XJSON(result)
	return err
}

func (j XJSON) String() string {
	return string(j)
}

// GormDataType gorm common data type
func (XJSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (XJSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (c BaseResourceMeta) GetID() int {
	return c.ID
}

func (c BaseResourceMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
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
