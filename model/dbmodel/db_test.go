package dbmodel

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/uopensail/recgo-engine/model/dbmodel/table"
)

func makeDefaultPool() table.PoolMeta {
	return table.PoolMeta{
		ID:           1,
		Name:         "default",
		Location:     "/tmp/sunmao/pool.json",
		PrimaryField: "d_s_id",

		UpdateTime: time.Now(),
	}
}

func makeFileRecallResource(id int, name, path string) table.RecallResourceMeta {
	return table.RecallResourceMeta{
		BaseResourceMeta: table.BaseResourceMeta{
			ID:         id,
			Name:       name,
			PluginName: "file",
			Source:     []byte(fmt.Sprintf("{\"location\":\"%s\"}", path)),
			UpdateTime: time.Now(),
		},
	}
}

func makeRedisRecallResource(id int) table.RecallResourceMeta {
	return table.RecallResourceMeta{
		BaseResourceMeta: table.BaseResourceMeta{
			ID:         id,
			Name:       "redisa",
			PluginName: "redis",
			Source:     []byte("{\"url\":\"redis://127.0.0.1:6379/2\"}"),
			UpdateTime: time.Now(),
		},
	}
}

func makeDSLJson(rawJson string) table.DSLMeta {

	ret := table.DSLMeta{}
	json.Unmarshal([]byte(rawJson), &ret)
	return ret
}

var recallUID int

func makeRecallEntity(name, dsl, paredDSL string) table.RecallEntityMeta {
	recallUID++
	return table.RecallEntityMeta{
		EntityMeta: table.EntityMeta{
			ID:         recallUID,
			Name:       name,
			UpdateTime: time.Now(),
			PluginName: "default",
		},
		DSL: dsl,

		DSLMeta: makeDSLJson(paredDSL),
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
			ID: 1,
			Params: map[string]string{
				"max_length": "150",
			},
			UpdateTime: time.Now(),
		},
		Condition: fmt.Sprintf("ts[int64] > %d", last7day),
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
		RecallGroupEntityID: 1,
		InsertGroupEntityID: 1,
		ScatterEntityID:     1,
	}
}

var insertUID int

func makeInsertEntity() (table.InsertEntityMeta, table.RecallEntityMeta) {
	insertUID++
	recallMeta := makeRecallEntity("__insert_recall_a__", `create recall __insert_a__ as select pool.id from sources.pool where concat("ctry",pool.d_s_country[string]) = user.u_s_country[string] and concat("lang",pool.d_s_language[string]) = user.u_s_language[string] limit 100
	`, `{"name":"__insert_a__","id":"id","paradigm":1,"from":{"resource":"pool"},"condition":{"index":{"indeces":[{"type":131,"value":{"left":{"type":8,"dtype":2,"value":{"func":"concat","args":[{"type":3,"value":"ctry"},{"type":7,"dtype":2,"value":{"table":"pool","column":"d_s_country"}}]}},"right":{"type":7,"dtype":2,"value":{"table":"user","column":"u_s_country"}},"op":"="}},{"type":131,"value":{"left":{"type":8,"dtype":2,"value":{"func":"concat","args":[{"type":3,"value":"lang"},{"type":7,"dtype":2,"value":{"table":"pool","column":"d_s_language"}}]}},"right":{"type":7,"dtype":2,"value":{"table":"user","column":"u_s_language"}},"op":"="}}]},"runtime_condition":"((concat(\"ctry\",pool.d_s_country[string]) = user.u_s_country[string]) AND (concat(\"lang\",pool.d_s_language[string]) = user.u_s_language[string]))"},"orderby":{"field":"id","desc":true},"limit":100}`)

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

		Condition: `concat("ctry",pool.d_s_country[string]) = user.u_s_country[string] and concat("lang",pool.d_s_language[string]) = user.u_s_language[string]`,
		Limit:     100,
		RecallID:  recallMeta.ID,
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
		PluginParams: []byte(`{"window_size":8, "group_limit": [{"field":"d_s_country","limit":2}]}`),
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
		Pools: []table.PoolMeta{makeDefaultPool()},
		RecallResources: []table.RecallResourceMeta{makeFileRecallResource(1, "cat_hotest", "/tmp/sunmao/cat_hotest_file.txt"),
			makeFileRecallResource(2, "w2v", "/tmp/sunmao/w2v_file.txt"),
			makeRedisRecallResource(3)},
		RecallEntities: []table.RecallEntityMeta{insertRecallMeta, makeRecallEntity("w2v",
			`
		create recall w2v as
		select w2v.id 
		from users.user join sources.w2v  on ( w2v.id[string] = user.u_d_click_list[strings] )
		where (not (w2v.d_s_cat1[string] in ("cat11","cat13","cat13"))) or (w2v.d_s_cat2[string] in ("cat21","cat22","cat23"))
		and concat("ctry",w2v.d_s_country[string]) = user.u_s_country[string]
		and concat("lang",w2v.d_s_language[string]) = user.u_s_language[string]
		having w2v.id not in filters.f1
		order by w2v.d_d_ctr asc 
		limit 100
		`,
			`
			{"name":"w2v","id":"id","paradigm":2,"from":{"resource":"w2v","key_format_exec_json":"{\"nodes\":[{\"id\":0,\"dtype\":5,\"value\":\"u_d_click_list\"}]}"},"condition":{"runtime_condition":"((w2v.d_s_cat1[string] NOT IN (\"cat11\",\"cat13\",\"cat13\")) OR (((w2v.d_s_cat2[string] IN (\"cat21\",\"cat22\",\"cat23\")) AND (concat(\"ctry\",w2v.d_s_country[string]) = user.u_s_country[string])) AND (concat(\"lang\",w2v.d_s_language[string]) = user.u_s_language[string])))"},"orderby":{"field":"d_d_ctr","desc":false},"limit":100,"filter":"f1"}
			`),
			makeRecallEntity("cat_hotest",
				`
		create recall cat_hotest as
select pool.id 
from sources.pool
where (not (pool.d_s_cat1[string] in ("cat11","cat13","cat13"))) or (pool.d_s_cat2[string] in ("cat21","cat22","cat23"))
and concat("ctry",pool.d_s_country[string]) = user.u_s_country[string]
and concat("lang",pool.d_s_language[string]) = user.u_s_language[string]
and pool.d_d_ctr[float32] > 0.3
having pool.id not in filters.f1
order by pool.d_d_ctr asc 
limit 100
		`,
				`
				{"name":"cat_hotest","id":"id","paradigm":1,"from":{"resource":"pool"},"condition":{"runtime_condition":"((pool.d_s_cat1[string] NOT IN (\"cat11\",\"cat13\",\"cat13\")) OR ((((pool.d_s_cat2[string] IN (\"cat21\",\"cat22\",\"cat23\")) AND (concat(\"ctry\",pool.d_s_country[string]) = user.u_s_country[string])) AND (concat(\"lang\",pool.d_s_language[string]) = user.u_s_language[string])) AND (pool.d_d_ctr[float32] > 0.3)))"},"orderby":{"field":"d_d_ctr","desc":false},"limit":100,"filter":"f1"}
								`,
			),
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
