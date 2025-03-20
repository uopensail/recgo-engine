package resources

import (
	"fmt"
	"testing"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/sample"
	"golang.org/x/exp/slices"
)

func Test_BuildCollection(t *testing.T) {
	pl, _ := pool.NewPool("/tmp/sunmao/pool.json")
	ress := &Resource{
		ResourceMeta: ResourceMeta{
			FieldDataType: map[string]sample.DataType{
				"d_s_cat1":     sample.StringType,
				"d_s_language": sample.StringType,
				"d_s_country":  sample.StringType,
				"u_s_coutry":   sample.StringType,
			},
		},
		Pool: pl,
	}

	collection := BuildCollection(ress, pl.WholeCollection, "d_s_cat1  in (\"cat11\")")
	collection = BuildCollection(ress, collection, "d_s_language  = \"en\"")
	condition := BuildCondition(ress, collection, "d_s_country =u_s_coutry")
	for i := 0; i < len(collection); i++ {
		item := pl.GetById(collection[i])
		if v, _ := item.Get("d_s_cat1").GetStrings(); slices.Contains(v, "cat11") {
			t.FailNow()
		}
		if v, _ := item.Get("d_s_language").GetString(); v != "en" {
			t.FailNow()
		}
		itemData, _ := item.Feats.MarshalJSON()
		fmt.Println(string(itemData))
	}
	userF := sample.NewMutableFeatures()
	userF.Set("u_s_coutry", &sample.String{Value: "us"})
	needCheck := []int{5520, 5583, 1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 9}

	reuslt := condition.Check(userF, needCheck)
	for i := 0; i < len(reuslt); i++ {
		item := pl.GetById(reuslt[i])

		if v, _ := item.Get("d_s_country").GetString(); v != "us" {
			t.FailNow()
		}
		itemData, _ := item.Feats.MarshalJSON()
		fmt.Println(string(itemData))
	}
}

func Test_TableFill(t *testing.T) {

}
