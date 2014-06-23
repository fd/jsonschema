package jsonschema

import (
	"sort"
)

type uniqueItemsValidator struct {
	unique bool
}

func (v *uniqueItemsValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("uniqueItems"); found {
		if y, ok := x.(bool); ok && y {
			v.unique = true
		}
	}
	return nil
}

func (v *uniqueItemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	var (
		l       = len(y)
		skipbuf [32]int
		skip    = skipbuf[:0]
	)

	if l > cap(skip) {
		skip = make([]int, 0, l)
	}

	for i := 0; i < l; i++ {
		if containsInt(skip, i) {
			continue
		}
		for j := i + 1; j < l; j++ {
			if containsInt(skip, j) {
				continue
			}

			a, b := y[i], y[j]

			equal, err := isEqual(a, b)
			if err != nil {
				skip = append(skip, j)
				sort.Ints(skip)
				ctx.Report(err)
				continue
			}

			// other values
			if equal {
				skip = append(skip, j)
				sort.Ints(skip)
				ctx.Report(&ErrNotUnique{i, j, a})
			}
		}
	}
}

func containsInt(s []int, x int) bool {
	idx := sort.SearchInts(s, x)
	if idx == len(s) {
		return false
	}
	return s[idx] == x
}
