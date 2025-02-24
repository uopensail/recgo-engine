package table

type RankEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`
}

func (c RankEntityMeta) GetID() int {
	return c.ID
}

func (c RankEntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}
