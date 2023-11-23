package attributes

type Map map[Attribute]any

func (r Map) Keys() (keys Attributes) {
	for a := range r {
		keys.Append(a)
	}

	keys.Sort()
	return
}

type Maps []map[Attribute]any
