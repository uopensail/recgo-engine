/*
type RedisConfigure struct {
	URL          string            `json:"url" toml:"url"`
	MinIdleConns int               `json:"min_idle_conns" toml:"min_idle_conns"`
	Timeout      int               `json:"timeout" toml:"timeout"`
	Params       map[string]string `json:"params" toml:"params"`
}

type FilterResourceMeta struct {
	ID         int       `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	PluginName string    `json:"plugin_name" toml:"plugin_name" gorm:"column:plugin_name"` //插件模式，相当于类名
	Name       string    `json:"name" toml:"name" gorm:"column:name"`
	UpdateTime time.Time `json:"update_time" toml:"update_time" gorm:"column:update_time"`

	Redis RedisConfigure `json:"redis" toml:"redis"`
}
*/

CREATE TABLE IF NOT EXISTS `filter_source` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `plugin_name` varchar(255) NOT NULL COMMENT '插件名',
  `name` varchar(255) NOT NULL COMMENT '规则名',
  `params` JSON NOT NULL COMMENT '额外参数 json:map[string]string',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `source` JSON NOT NULL COMMENT 'redis config json',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='过滤数据源'