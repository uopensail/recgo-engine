CREATE TABLE IF NOT EXISTS `weighted_entity` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `ab_layer_id` int(10) NOT NULL COMMENT '实验层ID',
  `plugin_name` varchar(255) NOT NULL COMMENT '插件名',
  `name` varchar(255) NOT NULL COMMENT '规则名',
  `params` JSON NOT NULL COMMENT '额外参数 json:map[string]string',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='调权计算实体'