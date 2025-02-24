package table

import "time"

// 单路召回计算实体
type ABMeta struct {
	ID         int       `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	RelateID   int       `json:"relate_id" toml:"relate_id" gorm:"column:relate_id"` //实体entity_id
	LayerID    int       `json:"layer_id" toml:"layer_id" gorm:"column:layer_id"`    //同一层有多个实验说明需要替换变种
	CaseID     int       `json:"case_id" toml:"case_id" gorm:"column:case_id"`
	UpdateTime time.Time `json:"update_time" toml:"update_time" gorm:"column:update_time"`
}

func (cfg ABMeta) GetID() int {
	return cfg.ID
}
func (cfg ABMeta) GetUpdateTime() int64 {
	return cfg.UpdateTime.Unix()
}

type ABEntityTableModel struct {
	Entities []ABMeta

	caseEntities map[int]*ABMeta
}

func (ab *ABEntityTableModel) Init(Entities []ABMeta) {

	ab.caseEntities = make(map[int]*ABMeta, len(Entities))
	for i := 0; i < len(Entities); i++ {
		entiy := &Entities[i]
		ab.caseEntities[entiy.CaseID] = entiy
	}
	ab.Entities = Entities
}

func (ab *ABEntityTableModel) Get(caseID int) *ABMeta {
	if v, ok := ab.caseEntities[caseID]; ok {
		return v
	}
	return nil
}
