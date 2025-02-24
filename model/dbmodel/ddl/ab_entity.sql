/*
type ABEntity struct {
	ID       int `json:"id" toml:"id" gorm:"primaryKey;column:id"`
	RelateID int `json:"relate_id" toml:"relate_id" gorm:"column:relate_id"`
	LayerID  int `json:"layer_id" toml:"layer_id" gorm:"column:layer_id"` //同一层有多个实验说明需要替换变种
	CaseID   int `json:"case_id" toml:"case_id" gorm:"column:case_id"`
}
*/
CREATE TABLE IF NOT EXISTS `ab_entity` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `relate_id` int(10) NOT NULL COMMENT '关联的计算entity_id',
  `layer_id` int(10) NOT NULL COMMENT '层ID',
  `case_id` int(10) NOT NULL COMMENT '子实验ID',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='ab实验记录表'