package freqs

// Filter stores a unique set of IDs with their corresponding keys.
// It maintains insertion order for both IDs and keys.
type Filter struct {
	dict map[int]struct{} // set of IDs for fast existence checking
	ids  []int            // ordered list of IDs
	keys []string         // ordered list of keys corresponding to IDs
}

// NewFilter creates and returns a new Filter with pre-allocated capacity.
func NewFilter() *Filter {
	return &Filter{
		dict: make(map[int]struct{}, 64),
		ids:  make([]int, 0, 64),
		keys: make([]string, 0, 64),
	}
}

// Merge adds all (key, id) pairs from another Filter into the current one.
// Only IDs not already present in the current Filter are added.
func (f *Filter) Merge(f1 *Filter) {
	for i := range f1.ids {
		f.Add(f1.keys[i], f1.ids[i])
	}
}

// Add inserts a new (key, id) pair if the ID is not already present.
// Keys and IDs are stored in the order they are added.
func (f *Filter) Add(key string, id int) {
	if _, exists := f.dict[id]; exists {
		return
	}
	f.dict[id] = struct{}{}
	f.ids = append(f.ids, id)
	f.keys = append(f.keys, key)
}

// Exists checks whether the given ID exists in the Filter.
func (f *Filter) Exists(id int) bool {
	_, exists := f.dict[id]
	return exists
}

// Missing checks whether the given ID is not present in the Filter.
// This is the inverse of Exists.
func (f *Filter) Missing(id int) bool {
	return !f.Exists(id)
}

// Exclude returns all stored keys in the order they were added.
func (f *Filter) Exclude() []string {
	return f.keys
}
