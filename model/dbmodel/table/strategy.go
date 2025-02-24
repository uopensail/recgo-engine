package table

type StrategyEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	FilterGroupEntityID int `json:"filter_group_entity_id" toml:"filter_group_entity_id" gorm:"column:filter_group_entity_id"`
	RecallGroupEntityID int `json:"recall_group_entity_id" toml:"recall_group_entity_id" gorm:"column:recall_group_entity_id"`
	RankEntityID        int `json:"rank_entity_id" toml:"rank_entity_id" gorm:"column:rank_entity_id"`
	WeightedEntityID    int `json:"weighted_entity_id" toml:"weighted_entity_id" gorm:"column:weighted_entity_id"`
	InsertGroupEntityID int `json:"insert_group_entity_id" toml:"insert_group_entity_id" gorm:"column:insert_group_entity_id"`
	ScatterEntityID     int `json:"scatter_entity_id" toml:"scatter_entity_id" gorm:"column:scatter_entity_id"`
}

func (c StrategyEntityMeta) GetID() int {
	return c.ID
}

func (c StrategyEntityMeta) GetUpdateTime() int64 {
	return c.UpdateTime.Unix()
}
