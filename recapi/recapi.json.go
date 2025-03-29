package recapi

import (
	"encoding/json"

	"github.com/bytedance/sonic"
	"github.com/uopensail/ulib/sample"
)

type RecRequestWrapper struct {
	*RecRequest `json:",inline"`
	UFeat       *sample.MutableFeatures    `json:"-"`
	FieldsType  map[string]sample.DataType `json:"-"`
}

func (d *RecRequestWrapper) FromRecRequest(in *RecRequest) {
	feat := sample.NewMutableFeatures()

	if in != nil {
		for k, v := range in.UserFeature {
			if v != nil {
				switch sample.DataType(v.Type) {
				case sample.Int64Type:
					if v.Value != nil {
						feat.Set(k, &sample.Int64{Value: v.Value.Iv})
					}
				case sample.Float32Type:
					if v.Value != nil {
						feat.Set(k, &sample.Float32{Value: v.Value.Fv})
					}
				case sample.StringType:
					if v.Value != nil {
						feat.Set(k, &sample.String{Value: v.Value.Sv})
					}
				case sample.Int64sType:
					if v.Value != nil {
						feat.Set(k, &sample.Int64s{Value: v.Value.Ivs})
					}
				case sample.Float32sType:
					if v.Value != nil {
						feat.Set(k, &sample.Float32s{Value: v.Value.Fvs})
					}
				case sample.StringsType:
					if v.Value != nil {
						feat.Set(k, &sample.Strings{Value: v.Value.Svs})
					}
				}
			}
		}
		d.RecRequest = in
		d.UFeat = feat
	}

}
func (d *RecRequestWrapper) UnmarshalJSON(data []byte) error {
	// 首先解析type和原始value
	type Alias RecRequestWrapper
	aux := &struct {
		UserFeature map[string]json.RawMessage `json:"user_feature"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	feat := sample.NewMutableFeatures()
	for k, v := range aux.UserFeature {
		if fieldType, ok := d.FieldsType[k]; ok {
			switch fieldType {
			case sample.Int64Type:
				var num int64
				err := sonic.Unmarshal(v, &num)
				if err == nil {
					feat.Set(k, &sample.Int64{Value: num})
				}

			case sample.Float32Type:
				var num float32
				err := sonic.Unmarshal(v, &num)
				if err == nil {
					feat.Set(k, &sample.Float32{Value: num})
				}
			case sample.StringType:
				var str string
				err := sonic.Unmarshal(v, &str)
				if err == nil {
					feat.Set(k, &sample.String{Value: str})
				}
			case sample.Int64sType:
				var nums []int64
				err := sonic.Unmarshal(v, &nums)
				if err == nil {
					feat.Set(k, &sample.Int64s{Value: nums})
				}
			case sample.Float32sType:
				var nums []float32
				err := sonic.Unmarshal(v, &nums)
				if err == nil {
					feat.Set(k, &sample.Float32s{Value: nums})
				}
			case sample.StringsType:
				var strs []string
				err := sonic.Unmarshal(v, &strs)
				if err == nil {
					feat.Set(k, &sample.Strings{Value: strs})
				}

			}
		}
	}
	return nil
}
