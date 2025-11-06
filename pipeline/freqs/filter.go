package freqs

type Filter struct {
	dict map[int]struct{}
	ids  []int
	keys []string
}

func NewFilter() *Filter {
	return &Filter{
		dict: make(map[int]struct{}, 64),
		ids:  make([]int, 0, 64),
		keys: make([]string, 0, 64),
	}
}

func (f *Filter) Merge(f1 *Filter) {
	for i := range f1.ids {
		f.Add(f1.keys[i], f1.ids[i])
	}
}

func (f *Filter) Add(key string, id int) {
	if _, ok := f.dict[id]; !ok {
		f.dict[id] = struct{}{}
		f.ids = append(f.ids, id)
		f.keys = append(f.keys, key)
	}
}

func (f *Filter) Check(id int) bool {
	_, ok := f.dict[id]
	return !ok
}

func (f *Filter) Exclude() []string {
	return f.keys
}
