package function

import (
	"reflect"

	"github.com/uopensail/ulib/sample"
)

type Float32 struct {
	Value float32   `json:"value" toml:"value"`
	Id    int       `json:"id" toml:"id"`
	Type  ArithType `json:"type" toml:"type"`
}

func (expr *Float32) Call(data []reflect.Value) {
	data[expr.Id] = reflect.ValueOf(expr.Value)
}

type Int64 struct {
	Value int64     `json:"value" toml:"value"`
	Id    int       `json:"id" toml:"id"`
	Type  ArithType `json:"type" toml:"type"`
}

func (expr *Int64) Call(data []reflect.Value) {
	data[expr.Id] = reflect.ValueOf(expr.Value)
}

type String struct {
	Value string    `json:"value" toml:"value"`
	Id    int       `json:"id" toml:"id"`
	Type  ArithType `json:"type" toml:"type"`
}

func (expr *String) Call(data []reflect.Value) {
	data[expr.Id] = reflect.ValueOf(expr.Value)
}

type Function struct {
	Args []int     `json:"args" toml:"args"`
	Func string    `json:"func" toml:"func"`
	Id   int       `json:"id" toml:"id"`
	Type ArithType `json:"type" toml:"type"`
}

func (expr *Function) Call(data []reflect.Value) {
	parms := make([]reflect.Value, len(expr.Args))
	for i := 0; i < len(expr.Args); i++ {
		parms[i] = data[expr.Args[i]]
	}
	if f, ok := Functions[expr.Func]; ok {
		ret := f.Call(parms)
		data[expr.Id] = ret[0]
	}
}

type Variable struct {
	Value    string          `json:"value" toml:"value"`
	Id       int             `json:"id" toml:"id"`
	Type     ArithType       `json:"type" toml:"type"`
	DataType sample.DataType `json:"dtype" toml:"dtype"`
}

func (expr *Variable) Call(data []reflect.Value) {}
