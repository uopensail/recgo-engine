/*
type FilterGroupEntityMeta struct {
	ID         int           `json:"id" toml:"id" gorm:"primaryKey;column:id"`                 // 对象实体ID
	ABLayerID  int           `json:"ab_layer_id" toml:"ab_layer_id" gorm:"column:ab_layer_id"` //绑定实验层
	Status     Entitiestatus `json:"status" toml:"status" gorm:"column:status"`
	Name       string        `json:"name" toml:"name" gorm:"column:name"` // 对象实体名
	UpdateTime time.Time     `json:"update_time" toml:"update_time" gorm:"column:update_time"`

	Params utils.StringMap `json:"params" toml:"params" gorm:"column:params"`

	Timeout int `json:"timeout" toml:"timeout" gorm:"column:timeout"`

	FilterEntitys utils.IntSlice `json:"filter_entitys" toml:"filter_entitys" gorm:"column:filter_entitys"` //一份引用
}
*/

CREATE TABLE IF NOT EXISTS `filter_group_entity` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `ab_layer_id` int(10) NOT NULL COMMENT '实验层ID',
  `name` varchar(255) NOT NULL COMMENT '规则名',
  `params` JSON NOT NULL COMMENT '额外参数 json:map[string]string',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `filter_Entities` JSON NOT NULL COMMENT 'json: int array filter_entity.id  ',
  `timeout` int(10) NOT NULL COMMENT '超时时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='过滤组计算实体'