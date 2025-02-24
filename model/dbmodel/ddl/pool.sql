/*

type ContentTableMeta struct {
	ID           int    `json:"id" toml:"ID" gorm:"column:id;primaryKey"`
	Name         string `json:"name" toml:"name" gorm:"column:name"`
	PrimaryField string `json:"primary_field" toml:"primary_field" gorm:"column:primary_field"`
	Location     string `json:"location" toml:"location" gorm:"column:location"`
	UpdateTime   int64  `json:"update_time" toml:"update_time" gorm:"column:update_time;autoCreateTime"`
}

*/
CREATE TABLE `content_source` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `primary_field` varchar(255) NOT NULL COMMENT '内容主键字段名',
  `name` varchar(255) NOT NULL COMMENT '规则名',
  `params` json NOT NULL COMMENT '额外参数 json:map[string]string',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `location` varchar(1024) NOT NULL COMMENT '对象存储地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='物料源'