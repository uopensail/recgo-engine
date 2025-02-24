/*

type RecallEntityMeta struct {
	ID         int             `json:"id" toml:"id"  gorm:"primaryKey;column:id"`
	Name       string          `json:"name" toml:"name" gorm:"column:name"`
	ABLayerID  int             `json:"ab_layer_id" toml:"ab_layer_id" gorm:"column:ab_layer_id"` //绑定实验层
	Status     Entitiestatus   `json:"status" toml:"status" gorm:"column:status"`
	UpdateTime time.Time       `json:"update_time" toml:"update_time" gorm:"column:update_time"`
	Params     utils.StringMap `json:"params" toml:"params" gorm:"column:params"`

	PluginName string `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"`

	DSL     string `json:"dsl" toml:"dsl" gorm:"column:dsl"`
	DSLMeta `json:"dsl_json" toml:"dsl_json" gorm:"column:dsl_json"`
}

*/

CREATE TABLE IF NOT EXISTS `recall_entity` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `ab_layer_id` int(10) NOT NULL COMMENT '实验层ID',
  `plugin_name` varchar(255) NOT NULL COMMENT '插件名',
  `name` varchar(255) NOT NULL COMMENT '规则名',
  `params` JSON NOT NULL COMMENT '额外参数 json:map[string]string',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  
  `dsl` varchar(1024) NOT NULL COMMENT '自定义dsl',
  `dsl_json` JSON NOT NULL COMMENT '解析后的dsl',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='匹配召回计算实体'