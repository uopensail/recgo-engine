package function

import (
	"fmt"
	"testing"

	"github.com/uopensail/ulib/sample"
)

func TestFunction(t *testing.T) {
	// concat(u4,concat(concat(lower("AAAA"), upper(u1)), u2),u3)
	code := `{
		"nodes":[
			{
				"type":0,
				"value":"u1",
				"dtype": 2,
				"id":0
			},
			{
				"type":4,
				"func":"upper",
				"args":[
					0
				],
				"id":1
			},
			{
				"type":3,
				"value":"AAAA",
				"id":2
			},
			{
				"type":4,
				"func":"lower",
				"args":[
					2
				],
				"id":3
			},
			{
				"type":4,
				"func":"concat",
				"args":[
					1,
					3
				],
				"id":4
			},
			{
				"type":0,
				"value":"u2",
				"dtype": 2,
				"id":5
			},
			{
				"type":4,
				"func":"concat",
				"args":[
					4,
					5
				],
				"id":6
			},
			{
				"type":0,
				"value":"u4",
				"dtype": 2,
				"id":7
			},
			{
				"type":0,
				"value":"u3",
				"dtype": 2,
				"id":8
			},
			{
				"type":4,
				"func":"concat",
				"args":[
					7,
					6,
					8
				],
				"id":9
			}
		]
	}`

	e := NewExpression(code)
	if e == nil {
		return
	}
	fmt.Println("test case 1:")
	feas1 := sample.NewMutableFeatures()
	feas1.Set("u1", &sample.Strings{Value: []string{"x", "y", "z"}})
	feas1.Set("u2", &sample.Strings{Value: []string{"-1", "-2", "-3", "-4"}})
	feas1.Set("u3", &sample.Strings{Value: []string{"-1", "-2", "-3", "-4", "-5"}})
	feas1.Set("u4", &sample.Strings{Value: []string{"S", "D"}})
	keys := e.Do(feas1)
	for _, k := range keys {
		fmt.Println(k)
	}

	fmt.Println("test case 2:")
	feas2 := sample.NewMutableFeatures()
	feas2.Set("u1", &sample.String{Value: "c"})
	feas2.Set("u2", &sample.Strings{Value: []string{"-1", "-2", "-3", "-4", "-5"}})
	feas2.Set("u3", &sample.Strings{Value: []string{"-a"}})
	feas2.Set("u4", &sample.Strings{Value: []string{"S"}})
	keys = e.Do(feas2)
	for _, k := range keys {
		fmt.Println(k)
	}

	// u1
	code1 := `{
		"nodes":[
			{
				"type":0,
				"value":"u1",
				"dtype": 0,
				"id":0
			},
			{
				"type":4,
				"func":"cast2str",
				"args":[
					0
				],
				"id":1
			}
		]
	}`

	e = NewExpression(code1)
	if e == nil {
		return
	}
	fmt.Println("test case 3:")
	feas1 = sample.NewMutableFeatures()
	feas1.Set("u1", &sample.Int64s{Value: []int64{1, 2, 3, 4}})
	keys = e.Do(feas1)
	for _, k := range keys {
		fmt.Println(k)
	}
}
