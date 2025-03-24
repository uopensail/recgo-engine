package table

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/uopensail/recgo-engine/model/utils"
	"gopkg.in/yaml.v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type EntityMeta struct {
	ID         int             `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	UpdateTime time.Time       `json:"update_time" toml:"update_time" gorm:"column:update_time"`
	Name       string          `json:"name" toml:"name" gorm:"column:name"`
	Status     Entitiestatus   `json:"status,omitempty" toml:"status" gorm:"column:status"`
	PluginName string          `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"`
	ABLayerID  string          `json:"ab_layer_id" toml:"ab_layer_id" gorm:"column:ab_layer_id"` //绑定实验层
	Params     utils.StringMap `json:"params" toml:"params" gorm:"column:params"`
}

func (c EntityMeta) GetID() int {
	return c.ID
}

func (c EntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
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
func LoadMeta[T any](filePath string, config *T) error {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 根据文件扩展名选择解析器
	switch ext := filepath.Ext(filePath); ext {
	case ".json":
		return json.Unmarshal(data, config)
	case ".toml":
		return toml.Unmarshal(data, config)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, config)
	default:
		return errors.New("unsupported config format: " + ext)
	}
}
