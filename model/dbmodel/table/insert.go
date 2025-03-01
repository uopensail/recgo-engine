package table

import (
	"github.com/uopensail/recgo-engine/model/utils"
)

type InsertEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	Bengin   int     `json:"begin" toml:"begin" gorm:"column:begin"`          //开始位置
	End      int     `json:"end" toml:"end" gorm:"column:end"`                //结束位置
	Prob     float32 `json:"prob" toml:"prob" gorm:"column:prob"`             //是否强插的概率
	Priority float32 `json:"priority" toml:"priority" gorm:"column:priority"` //优先级

	//for recall entity
	RecallID  int    `json:"recall_id" toml:"recall_id" gorm:"column:recall_id"`
	Condition string `json:"condition" toml:"condition" gorm:"column:condition"`
	Limit     int    `json:"limit" toml:"limit" gorm:"column:limit"`
}

// 召回组计算实体
type InsertGroupEntityMeta struct {
	EntityMeta `json:",inline" toml:",inline" gorm:"embedded"`

	InsertEntities utils.IntSlice `json:"insert_entities" toml:"insert_entities" gorm:"column:insert_entities"` //一份引用
}
