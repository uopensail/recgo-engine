package table

type ScatterEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	PluginParams XJSON `json:"plugin_params" toml:"plugin_params" gorm:"column:plugin_params"`
}
