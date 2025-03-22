package dbmodel

import (
	"fmt"
	"testing"
	"time"

	"github.com/uopensail/recgo-engine/model/dbmodel/table"
)

var recallUID int

func makeConditionRecallEntity(name, condition string) table.RecallEntityMeta {
	recallUID++
	return table.RecallEntityMeta{
		Condition: condition,
		EntityMeta: table.EntityMeta{
			ID:         recallUID,
			Name:       name,
			UpdateTime: time.Now(),
			PluginName: "condition",
		},
	}
}

func makeInvertIndexRecallEntity(name, condition string, fields []string) table.RecallEntityMeta {
	recallUID++
	return table.RecallEntityMeta{
		Condition: condition,
		EntityMeta: table.EntityMeta{
			ID:         recallUID,
			Name:       name,
			UpdateTime: time.Now(),
			PluginName: "invert_index",
		},
		PluginParams: []byte("{\"resource\":\"d_s_country|d_s_cat\",\"user_feature_fields\":[\"u_s_country\",\"u_s_cat\"]}"),
	}
}

func makeRecallGroupEntity() table.RecallGroupEntityMeta {
	return table.RecallGroupEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         1,
			Name:       "home",
			UpdateTime: time.Now(),
		},
		RecallEntities: []int{2, 3},
	}
}

func makeRedisFilterResource() table.FilterResourceMeta {
	return table.FilterResourceMeta{
		ID:         1,
		Name:       "f1",
		PluginName: "redis",
		Redis: table.RedisConfigure{
			URL: "redis://127.0.0.1:6379/1",
		},
		UpdateTime: time.Now(),
	}
}

func makeFilterEntity() table.FilterEntityMeta {
	last7day := time.Now().Unix() - 3*86400
	return table.FilterEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         1,
			Params:     map[string]string{"max_length": "150"},
			UpdateTime: time.Now(),
		},
		Condition: fmt.Sprintf("ts > %d", last7day),
		Format:    "imp_%s",
		MaxCount:  2,
		SourceID:  1,
	}
}

func makeFilterGroupEntity() table.FilterGroupEntityMeta {

	return table.FilterGroupEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         1,
			Name:       "f1",
			UpdateTime: time.Now(),
		},
		FilterEntities: []int{1},
	}
}

func makeStrategyEntity() table.StrategyEntityMeta {
	return table.StrategyEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         1,
			Name:       "home",
			PluginName: "default",
			UpdateTime: time.Now(),
		},
		FilterGroupEntityID: 1,
		RecallGroupEntityID: 1,
		InsertGroupEntityID: 1,
		ScatterEntityID:     1,
	}
}

var insertUID int

func makeInsertEntity() (table.InsertEntityMeta, table.RecallEntityMeta) {
	insertUID++
	recallMeta := makeConditionRecallEntity("__insert_recall_a__", `i_s_country = u_s_country and i_s_language =  u_s_language`)

	insertMeta := table.InsertEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         insertUID,
			Name:       "a",
			UpdateTime: time.Now(),
			PluginName: "default",
		},
		Bengin:   0,
		End:      5,
		Priority: 0.1,
		Prob:     1.0,
		Limit:    100,
		RecallID: recallMeta.ID,
	}
	return insertMeta, recallMeta
}

func makeScatterEntity() table.ScatterEntityMeta {
	return table.ScatterEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         1,
			Name:       "cat1",
			UpdateTime: time.Now(),
			PluginName: "window",
		},
		PluginParams: []byte(`{"window_size":8, "group_limit": [{"field":"i_s_country","limit":2}]}`),
	}
}

func makeInsertGroupEntity() table.InsertGroupEntityMeta {

	return table.InsertGroupEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         1,
			Name:       "home",
			UpdateTime: time.Now(),
		},
		InsertEntities: []int{1},
	}
}

func Test_dBTablesDump(t *testing.T) {
	insertMeta, insertRecallMeta := makeInsertEntity()
	tables := DBTables{
		RecallEntities: []table.RecallEntityMeta{insertRecallMeta,
			makeConditionRecallEntity("w2v",
				`(not (i_s_cat1 in ("cat11","cat13","cat13"))) or (i_s_cat2 in ("cat21","cat22","cat23"))
		and  i_s_country  = u_s_country
		and  i_s_language  = u_s_language`),
			makeInvertIndexRecallEntity("country_invert_index", "i_s_language =  u_s_language", []string{"i_s_coutry"}),
		},
		RecallGroupEntities: []table.RecallGroupEntityMeta{makeRecallGroupEntity()},

		FilterResources:     []table.FilterResourceMeta{makeRedisFilterResource()},
		FilterEntities:      []table.FilterEntityMeta{makeFilterEntity()},
		FilterGroupEntities: []table.FilterGroupEntityMeta{makeFilterGroupEntity()},
		InsertEntities:      []table.InsertEntityMeta{insertMeta},
		InsertGroupEntities: []table.InsertGroupEntityMeta{makeInsertGroupEntity()},
		ScatterEntities:     []table.ScatterEntityMeta{makeScatterEntity()},

		StrategyEntities: []table.StrategyEntityMeta{makeStrategyEntity()},
	}
	tables.Dump("/tmp/dbmodel.toml")
	tableModel, _ := LoadDBTabelModel("/tmp/dbmodel.toml")
	fmt.Println(tableModel)
}
