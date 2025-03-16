package table

import (
	"github.com/uopensail/ulib/sample"
)

type PoolMeta struct {
	FieldDataType map[string]sample.DataType `json:"field_data_type" toml:"field_data_type"`
}
