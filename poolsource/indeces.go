package poolsource

import (
	"github.com/spf13/cast"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
)

type Index struct {
	data map[string]Collection
}

func (index *Index) Get(key interface{}) Collection {
	if array, ok := index.data[cast.ToString(key)]; ok {
		return array
	}
	return nil
}

func BuildIndex(pl *pool.Pool, collection Collection, column string) *Index {
	stat := prome.NewStat("Recall.BuildIndex." + column)
	defer stat.End()
	if len(collection) <= 0 {
		return nil
	}

	// generate func map
	funcs := make(map[sample.DataType]func(feature sample.Feature) []string)

	funcs[sample.Float32Type] = func(feature sample.Feature) []string {
		val, _ := feature.GetFloat32()
		return []string{cast.ToString(val)}
	}

	funcs[sample.Float32sType] = func(feature sample.Feature) []string {
		vals, _ := feature.GetFloat32s()
		ret := make([]string, len(vals))
		for i := 0; i < len(vals); i++ {
			ret[i] = cast.ToString(vals[i])
		}
		return ret
	}

	funcs[sample.Int64Type] = func(feature sample.Feature) []string {
		val, _ := feature.GetInt64()
		return []string{cast.ToString(val)}
	}

	funcs[sample.Int64sType] = func(feature sample.Feature) []string {
		vals, _ := feature.GetInt64s()
		ret := make([]string, len(vals))
		for i := 0; i < len(vals); i++ {
			ret[i] = cast.ToString(vals[i])
		}
		return ret
	}

	funcs[sample.StringType] = func(feature sample.Feature) []string {
		val, _ := feature.GetString()
		return []string{val}
	}

	funcs[sample.StringsType] = func(feature sample.Feature) []string {
		vals, _ := feature.GetStrings()
		return vals
	}
	firstRecord := pl.GetById(collection[0])
	foo := funcs[firstRecord.Get(column).Type()]
	dict := make(map[string]Collection)

	for i := 0; i < len(collection); i++ {
		itemID := collection[i]
		record := pl.GetById(itemID)
		if record != nil {
			strs := foo(record.Get(column))

			for j := 0; j < len(strs); j++ {
				if list, ok := dict[strs[j]]; ok {
					list = append(list, itemID)
					dict[strs[j]] = list
				} else {
					dict[strs[j]] = []int{itemID}
				}
			}
		}
	}

	return &Index{
		data: dict,
	}
}

type PoolIndeces struct {
	indeces map[string]*Index // pool column indeces
	pl      *pool.Pool
}

func NewPoolIndeces(pl *pool.Pool) *PoolIndeces {
	return &PoolIndeces{
		indeces: make(map[string]*Index),
		pl:      pl,
	}
}

func (pi *PoolIndeces) BuildColumnIndex(column string) {
	if _, ok := pi.indeces[column]; !ok {
		pi.indeces[column] = BuildIndex(pi.pl, pi.pl.WholeCollection, column)
	}
}

func (pi *PoolIndeces) GetIndex(column string) *Index {
	if v, ok := pi.indeces[column]; ok {
		return v
	}
	return nil
}

func (pi *PoolIndeces) GetIndexCollection(column string, key string) Collection {
	if v, ok := pi.indeces[column]; ok {
		if vc, ok2 := v.data[key]; ok2 {
			return vc
		}
	}
	return nil
}
