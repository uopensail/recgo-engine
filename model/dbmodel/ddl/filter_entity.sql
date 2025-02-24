/*
type FilterEntityMeta struct {
	ID         int           `json:"id" toml:"id" gorm:"primaryKey;column:id"`                 // 对象实体ID
	ABLayerID  int           `json:"ab_layer_id" toml:"ab_layer_id" gorm:"column:ab_layer_id"` //绑定实验层
	Status     Entitiestatus `json:"status" toml:"status" gorm:"column:status"`
	Name       string        `json:"name" toml:"name" gorm:"column:name"` // 对象实体名
	UpdateTime time.Time     `json:"update_time" toml:"update_time" gorm:"column:update_time"`

	PluginName string `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"` //插件模式，相当于类名

	Condition string `json:"condition" toml:"condition" gorm:"column:condition"`
	MaxCount  int    `json:"max_count" toml:"max_count" gorm:"column:max_count"`
	Format    string `json:"format" toml:"format" gorm:"column:format"`

	SourceID int             `json:"source_id" toml:"source_id" gorm:"column:source_id"`
	Params   utils.StringMap `json:"params" toml:"params" gorm:"column:params"`
}

*/
CREATE TABLE IF NOT EXISTS `filter_entity` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `ab_layer_id` int(10) NOT NULL COMMENT '实验层ID',
  `plugin_name` varchar(255) NOT NULL COMMENT '插件名',
  `name` varchar(255) NOT NULL COMMENT '规则名',
  `params` JSON NOT NULL COMMENT '额外参数 json:map[string]string',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `condition` varchar(1024) COMMENT '过滤条件',
  `max_count` int(10) NOT NULL COMMENT '过滤聚合截取',
  `format` varchar(255)  COMMENT '用户查找key的format',
  `source_id` int(10) NOT NULL COMMENT 'filter_source.id 过滤源ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='过滤计算实体'