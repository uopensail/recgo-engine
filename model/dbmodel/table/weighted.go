package table

type WeightedEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`
}
