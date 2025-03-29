package userctx

import (
	"github.com/uopensail/ulib/sample"
)

type UserFeatures struct {
	UFeat *sample.MutableFeatures
}

func (feat UserFeatures) AddTestClickList() {
	feat.UFeat.Set("u_d_click_list", &sample.Strings{Value: []string{"item_id_4589", "item_id_4408"}})

}
