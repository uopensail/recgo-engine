package table

type RankEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`
}
