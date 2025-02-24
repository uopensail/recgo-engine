package function

import (
	"encoding/json"
	"reflect"

	"github.com/tidwall/gjson"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type ArithType int

const (
	kVariable ArithType = iota
	kInt64
	kFloat32
	kString
	kFunction
	kError ArithType = 127
)

type Expression struct {
	vars  []Variable
	funcs []Function
	data  []reflect.Value
}

func NewExpression(code string) *Expression {
	stat := prome.NewStat("NewExpression")
	defer stat.End()
	type jsonNodes struct {
		Nodes []json.RawMessage `json:"nodes"`
	}

	nodes := &jsonNodes{}
	err := json.Unmarshal([]byte(code), nodes)
	if err != nil {
		stat.MarkErr()
		zlog.LOG.Error("express parse error", zap.String("code", code), zap.Error(err))
		return nil
	}
	vars := make([]Variable, 0, len(nodes.Nodes))
	data := make([]reflect.Value, len(nodes.Nodes))
	funcs := make([]Function, 0, len(nodes.Nodes))
	for i := 0; i < len(nodes.Nodes); i++ {
		bytes, _ := nodes.Nodes[i].MarshalJSON()
		atype := ArithType(gjson.Get(string(bytes), "type").Int())
		switch atype {
		case kVariable:
			node := Variable{}
			json.Unmarshal(bytes, &node)
			vars = append(vars, node)
		case kInt64:
			node := Int64{}
			json.Unmarshal(bytes, &node)
			data[node.Id] = reflect.ValueOf(node.Value)
		case kFloat32:
			node := Float32{}
			json.Unmarshal(bytes, &node)
			data[node.Id] = reflect.ValueOf(node.Value)
		case kString:
			node := String{}
			json.Unmarshal(bytes, &node)
			data[node.Id] = reflect.ValueOf(node.Value)
		case kFunction:
			node := Function{}
			json.Unmarshal(bytes, &node)
			funcs = append(funcs, node)
		}
	}
	return &Expression{
		vars:  vars,
		funcs: funcs,
		data:  data,
	}
}

func (e *Expression) Do(feas sample.Features) []string {
	stat := prome.NewStat("Expression.Do")
	defer stat.End()
	values := make([][]reflect.Value, len(e.vars))
	sizes := make([]int, len(e.vars))
	count := 1
	for i := 0; i < len(e.vars); i++ {
		tmp := getFeature(e.vars[i].Value, e.vars[i].DataType, feas)
		if len(tmp) == 0 {
			stat.MarkErr()
			zlog.LOG.Error("feature nil", zap.String("key", e.vars[i].Value))
			return nil
		}
		values[i] = tmp
		sizes[i] = len(tmp)
		count *= sizes[i]
	}
	indices := make([]int, len(sizes))
	index := count
	for i := 0; i < len(sizes); i++ {
		indices[i] = index / sizes[i]
		index = indices[i]
	}

	ret := make([]string, 0, count)
	for i := 0; i < count; i++ {
		idx := getIndeces(i, indices)
		for j := 0; j < len(idx); j++ {
			e.data[e.vars[j].Id] = values[j][idx[j]]
		}
		for j := 0; j < len(e.funcs); j++ {
			e.funcs[j].Call(e.data)
		}
		ret = append(ret, e.data[len(e.data)-1].String())
	}
	stat.SetCounter(len(ret))
	return ret
}

func getIndeces(index int, indices []int) []int {
	ret := make([]int, len(indices))
	for i := 0; i < len(indices); i++ {
		ret[i] = index / indices[i]
		index %= indices[i]
	}
	return ret
}

func getFeature(key string, dtype sample.DataType, feas sample.Features) []reflect.Value {
	fea := feas.Get(key)
	if fea == nil {
		return nil
	}
	ret := make([]reflect.Value, 0, 1)
	datatype := fea.Type()
	switch dtype {
	case sample.Float32Type, sample.Float32sType:
		if datatype == sample.Float32sType {
			val, err := fea.GetFloat32s()
			if err == nil {
				for i := 0; i < len(val); i++ {
					ret = append(ret, reflect.ValueOf(val[i]))
				}
			}
		} else if datatype == sample.Float32Type {
			val, err := fea.GetFloat32()
			if err == nil {
				ret = append(ret, reflect.ValueOf(val))
			}
		}
	case sample.Int64Type, sample.Int64sType:
		if datatype == sample.Int64sType {
			val, err := fea.GetInt64s()
			if err == nil {
				for i := 0; i < len(val); i++ {
					ret = append(ret, reflect.ValueOf(val[i]))
				}
			}
		} else if datatype == sample.Int64Type {
			val, err := fea.GetInt64()
			if err == nil {
				ret = append(ret, reflect.ValueOf(val))
			}
		}
	case sample.StringType, sample.StringsType:
		if datatype == sample.StringsType {
			val, err := fea.GetStrings()
			if err == nil {
				for i := 0; i < len(val); i++ {
					ret = append(ret, reflect.ValueOf(val[i]))
				}
			}
		} else if datatype == sample.StringType {
			val, err := fea.GetString()
			if err == nil {
				ret = append(ret, reflect.ValueOf(val))
			}
		}
	}
	return ret
}
