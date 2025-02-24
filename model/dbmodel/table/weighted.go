package table

type WeightedEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`
}

func (c WeightedEntityMeta) GetID() int {
	return c.ID
}

func (c WeightedEntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}
