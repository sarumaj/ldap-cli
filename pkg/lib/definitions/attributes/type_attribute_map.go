package attributes

import (
	"bytes"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type Map map[Attribute]any

func (r Map) String() string {
	m := make(map[string]any)
	for a, v := range r {
		m[a.String()] = v
	}

	buffer := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buffer)
	enc.SetIndent(2)
	if err := enc.Encode(m); err != nil {
		return fmt.Sprint(m)
	}

	return buffer.String()
}

type Maps []map[Attribute]any

func (rs Maps) String() string {
	var maps []map[string]any
	for _, r := range rs {
		m := make(map[string]any)
		for a, v := range r {
			m[a.String()] = v
		}

		maps = append(maps, m)
	}

	buffer := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buffer)
	enc.SetIndent(2)
	if err := enc.Encode(maps); err != nil {
		return fmt.Sprint(maps)
	}

	return buffer.String()
}
