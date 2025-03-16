package resources

import (
	"unsafe"

	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/uno"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type Condition struct {
	evaluator *uno.Evaluator
	slices    [][]unsafe.Pointer
}

func BuildCondition(ress *Resource, collection Collection, condition string) *Condition {
	stat := prome.NewStat("BuildCondition")
	defer stat.End()
	evaluator, err := uno.NewEvaluator(condition, ress.FieldDataType)
	if err != nil {
		zlog.LOG.Error("build condition error", zap.Error(err))
		stat.MarkErr()
		return nil
	}
	c := &Condition{
		evaluator: evaluator,
		slices:    make([][]unsafe.Pointer, ress.Pool.Len()),
	}

	for i := 0; i < len(collection); i++ {
		slice := evaluator.Allocate()
		itemID := collection[i]
		item := ress.Pool.GetById(itemID)
		evaluator.Fill(&item.Feats, slice)
		evaluator.PreEval(slice)
		c.slices[itemID] = slice
	}
	return c
}

func (c *Condition) Check(features sample.Features, collection Collection) Collection {
	stat := prome.NewStat("Condition.Check")
	defer stat.End()
	var slice []unsafe.Pointer
	var oldSlice []unsafe.Pointer
	var newSlice []unsafe.Pointer

	slices := make([][]unsafe.Pointer, 0, len(collection))
	slice = c.evaluator.Allocate()
	c.evaluator.Fill(features, slice)
	address := make([]uintptr, len(slice))
	for i := 0; i < len(slice); i++ {
		address[i] = uintptr(slice[i])
	}

	for i := 0; i < len(collection); i++ {
		oldSlice = c.slices[collection[i]]
		newSlice = make([]unsafe.Pointer, len(oldSlice))
		for j := 0; j < len(oldSlice); j++ {
			newSlice[j] = unsafe.Pointer(uintptr(oldSlice[j]) | address[j])
		}
		slices = append(slices, newSlice)
	}

	ret := make([]int, 0, len(collection))
	results := c.evaluator.BatchEval(slices)
	for i := 0; i < len(results); i++ {
		if results[i] == 1 {
			ret = append(ret, collection[i])
		}
	}
	stat.SetCounter(len(ret))
	return ret
}

func (c *Condition) CheckWithFillRuntime(features sample.Features, collection Collection, itemTableName string,
	onGetItem func(id int, indexInCollection int) sample.Features) Collection {
	stat := prome.NewStat("Condition.Check")
	defer stat.End()
	var slice []unsafe.Pointer
	var oldSlice []unsafe.Pointer
	var newSlice []unsafe.Pointer

	slices := make([][]unsafe.Pointer, 0, len(collection))
	slice = c.evaluator.Allocate()
	c.evaluator.Fill(features, slice)
	address := make([]uintptr, len(slice))
	for i := 0; i < len(slice); i++ {
		address[i] = uintptr(slice[i])
	}

	for i := 0; i < len(collection); i++ {
		id := collection[i]

		oldSlice = c.slices[id]
		newSlice = make([]unsafe.Pointer, len(oldSlice))
		for j := 0; j < len(oldSlice); j++ {
			newSlice[j] = unsafe.Pointer(uintptr(oldSlice[j]) | address[j])
		}

		if onGetItem != nil {
			itemF := onGetItem(id, i)
			if itemF != nil {

				c.evaluator.Fill(itemF, newSlice)
			}
		}
		slices = append(slices, newSlice)
	}

	ret := make([]int, 0, len(collection))
	results := c.evaluator.BatchEval(slices)
	for i := 0; i < len(results); i++ {
		if results[i] == 1 {
			ret = append(ret, collection[i])
		}
	}
	stat.SetCounter(len(ret))
	return ret
}

func (c *Condition) CheckAll(featureTable string, features sample.Features) Collection {
	stat := prome.NewStat("Condition.CheckAll")
	defer stat.End()
	var slice []unsafe.Pointer
	var newSlice []unsafe.Pointer

	slices := make([][]unsafe.Pointer, 0, len(c.slices))
	slice = c.evaluator.Allocate()
	address := make([]uintptr, len(slice))
	for i := 0; i < len(slice); i++ {
		address[i] = uintptr(slice[i])
	}
	c.evaluator.Fill(features, slice)

	for i := 0; i < len(c.slices); i++ {
		newSlice = make([]unsafe.Pointer, len(c.slices[i]))
		for j := 0; j < len(slice); j++ {
			newSlice[j] = unsafe.Pointer(uintptr(c.slices[i][j]) | address[j])
		}
		slices = append(slices, newSlice)
	}

	ret := make([]int, 0, len(c.slices))
	results := c.evaluator.BatchEval(slices)
	for i := 0; i < len(results); i++ {
		if results[i] == 1 {
			ret = append(ret, i)
		}
	}
	stat.SetCounter(len(ret))
	return ret
}

func (c *Condition) Release() {
	for i := 0; i < len(c.slices); i++ {
		c.evaluator.Clean(c.slices[i])
	}
	c.evaluator.Release()
}
