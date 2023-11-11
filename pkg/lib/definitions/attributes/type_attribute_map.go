package attributes

type Map map[Attribute]any

func (r Map) Keys() (keys Attributes) {
	for a := range r {
		keys = append(keys, a)
	}

	keys.Sort()
	return
}

type Maps []map[Attribute]any
