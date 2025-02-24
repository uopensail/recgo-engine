package poolsource

import (
	"fmt"
	"testing"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/sample"
	"golang.org/x/exp/slices"
)

func Test_BuildCollection(t *testing.T) {
	pl, _ := pool.NewPool("/tmp/sunmao/pool.json")
	collection := BuildCollection(pl, pl.WholeCollection, "", "d_s_cat1[strings] in (\"cat11\")")
	collection = BuildCollection(pl, collection, "", "d_s_language[string] = \"en\"")
	condition := BuildCondition(pl, collection, "", "d_s_country[string]=user.u_s_coutry[string]")
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

	reuslt := condition.Check("user", userF, needCheck)
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
