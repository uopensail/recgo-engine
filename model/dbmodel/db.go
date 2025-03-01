package dbmodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/prome"

	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBTables struct {
	Pools               []table.PoolMeta              `json:"pool" toml:"pool"`
	RecallResources     []table.RecallResourceMeta    `json:"recall_resources" toml:"recall_resources"`
	FilterResources     []table.FilterResourceMeta    `json:"filter_resources" toml:"filter_resources"`
	AbMetas             []table.ABMeta                `json:"ab_metas" toml:"ab_metas"`
	RecallEntities      []table.RecallEntityMeta      `json:"recall_entities" toml:"recall_entities"`
	RecallGroupEntities []table.RecallGroupEntityMeta `json:"recall_group_entities" toml:"recall_group_entities"`

	FilterEntities      []table.FilterEntityMeta      `json:"filter_entities" toml:"filter_entities"`
	FilterGroupEntities []table.FilterGroupEntityMeta `json:"filter_group_entities" toml:"filter_group_entities"`
	RankEntities        []table.RankEntityMeta        `json:"rank_entities" toml:"rank_entities"`
	WeightedEntities    []table.WeightedEntityMeta    `json:"weighted_entities" toml:"weighted_entities"`
	InsertEntities      []table.InsertEntityMeta      `json:"insert_entities" toml:"insert_entities"`
	InsertGroupEntities []table.InsertGroupEntityMeta `json:"insert_group_entities" toml:"insert_group_entities"`
	ScatterEntities     []table.ScatterEntityMeta     `json:"scatter_entities" toml:"scatter_entities"`

	StrategyEntities []table.StrategyEntityMeta `json:"strategy_entities" toml:"strategy_entities"`
}

func (tables *DBTables) Dump(path string) {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		panic(err)
	}
	defer fd.Close()
	if strings.HasSuffix(path, "yaml") {
		yaml.NewEncoder(fd).Encode(&tables)
	} else if strings.HasSuffix(path, "json") {
		json.NewEncoder(fd).Encode(&tables)
	} else {
		toml.NewEncoder(fd).Encode(&tables)
	}

}

func load(db *gorm.DB, tableName string, dest interface{}) {
	stat := prome.NewStat(fmt.Sprintf("load.%s", tableName))
	defer stat.End()
	result := db.Table(tableName).Where("status = 1").Find(dest)
	if result.Error != nil {
		zlog.LOG.Error("load", zap.Error(result.Error), zap.String("table_name", tableName))
		stat.MarkErr()
	}
	stat.SetCounter(int(result.RowsAffected))
}

func LoadDBTables(metaPath string) (DBTables, error) {
	var tables DBTables
	if strings.HasPrefix(metaPath, "oss://") || strings.HasPrefix(metaPath, "s3://") || strings.HasPrefix(metaPath, "/") {
		//TODO: read from object file
		tb, err := loadAllablesFromFile(metaPath)
		if err != nil {
			return DBTables{}, err
		}
		tables = tb
	} else {
		tb, err := loadAllablesFromDB(metaPath)
		if err != nil {
			return DBTables{}, err
		}
		tables = tb
	}
	return tables, nil
}
func LoadDBTabelModel(metaPath string) (DBTabelModel, error) {
	tables, err := LoadDBTables(metaPath)
	if err != nil {
		return DBTabelModel{}, err
	}
	dbModel := DBTabelModel{}
	dbModel.resourceTableModel.Init(tables.Pools)
	dbModel.RecallSourceTableModel.Init(tables.RecallResources)
	dbModel.ABEntityTableModel.Init(tables.AbMetas)
	dbModel.RecallEntityTableModel.Init(tables.RecallEntities)
	dbModel.RecallGroupEntityTableModel.Init(tables.RecallGroupEntities)
	dbModel.FilterResourceTableModel.Init(tables.FilterResources)
	dbModel.FilterEntityTableModel.Init(tables.FilterEntities)
	dbModel.FilterGroupEntityTableModel.Init(tables.FilterGroupEntities)
	dbModel.RankEntityTableModel.Init(tables.RankEntities)
	dbModel.WeightedEntityTableModel.Init(tables.WeightedEntities)
	dbModel.InsertEntityTableModel.Init(tables.InsertEntities)
	dbModel.InsertGroupEntityTableModel.Init(tables.InsertGroupEntities)
	dbModel.ScatterEntityTableModel.Init(tables.ScatterEntities)
	dbModel.StrategyEntityTableModel.Init(tables.StrategyEntities)

	return dbModel, nil
}

func loadAllablesFromDB(dsn string) (DBTables, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.LOG.Error("gorm Open", zap.Error(err))
		return DBTables{}, err
	}
	var dbModel DBTables
	load(db, "pools", &dbModel.Pools)
	load(db, "recall_resource", &dbModel.RecallResources)
	load(db, "filter_resource", &dbModel.FilterResources)
	load(db, "ab_entity", &dbModel.AbMetas)
	load(db, "recall_entity", &dbModel.RecallEntities)
	load(db, "recall_group_entity", &dbModel.RecallGroupEntities)

	load(db, "filter_entity", &dbModel.FilterEntities)
	load(db, "filter_group_entity", &dbModel.FilterGroupEntities)

	load(db, "rank_entity", &dbModel.RankEntities)
	load(db, "weighted_entity", &dbModel.WeightedEntities)
	load(db, "insert_entity", &dbModel.InsertEntities)
	load(db, "insert_group_entity", &dbModel.InsertGroupEntities)
	load(db, "scatter_entity", &dbModel.ScatterEntities)

	load(db, "strategy_entity", &dbModel.StrategyEntities)
	return dbModel, nil
}

func loadAllablesFromFile(filePath string) (DBTables, error) {
	fData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("os.ReadFile error: %s\n", err)
		return DBTables{}, err
	}
	buf := bytes.NewBuffer(fData)
	var dbModel DBTables
	if strings.HasSuffix(filePath, ".json") {

		err = json.NewDecoder(buf).Decode(&dbModel)
		if err != nil {
			fmt.Printf("Unmarshal error: %s\n", err)
			return DBTables{}, err
		}
	} else if strings.HasSuffix(filePath, ".yaml") {
		err = yaml.NewDecoder(buf).Decode(&dbModel)
		if err != nil {
			fmt.Printf("Unmarshal error: %s\n", err)
			return DBTables{}, err
		}
	} else {
		_, err = toml.NewDecoder(buf).Decode(&dbModel)
		if err != nil {
			fmt.Printf("Unmarshal error: %s\n", err)
			return DBTables{}, err
		}
	}

	return dbModel, nil
}
