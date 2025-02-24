package table

type ScatterEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	PluginParams XJSON `json:"plugin_params" toml:"plugin_params" gorm:"column:plugin_params"`
}

func (c ScatterEntityMeta) GetID() int {
	return c.ID
}

func (c ScatterEntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}
