package sliceutils

import (
	"reflect"
)

type Slice []any

func (sl *Slice) Pop() any {
	if len(*sl) == 0 {
		return nil
	}
	lastIndex := len(*sl) - 1
	lastElement := (*sl)[lastIndex]

	*sl = (*sl)[:lastIndex]

	return lastElement
}

func (sl *Slice) PushBack(v any) {
	if sl == nil {
		return
	}
	*sl = append(*sl, v)
}

func (sl *Slice) PushFront(v any) {
	if sl == nil {
		return
	}
	*sl = append(Slice{v}, *sl...)
}

func (sl *Slice) Clear() {
	*sl = Slice{}
}

func (sl Slice) Count(v any) int {
	count := 0
	for _, val := range sl {
		if val == v {
			count++
		}
	}
	return count
}

func (sl Slice) Contains(v any) bool {
	for _, val := range sl {
		if val == v {
			return true
		}
	}
	return false
}

func (sl Slice) ForEach(f func(any)) {
	for _, v := range sl {
		f(v)
	}
}

func (sl Slice) Map(f func(v any) any) Slice {
	var mappedSlice Slice
	for _, v := range sl {
		mappedSlice = append(mappedSlice, f(v))
	}
	return mappedSlice
}

func (sl Slice) Filter(f func(v any) bool) Slice {
	var filteredSlice Slice
	for _, v := range sl {
		if f(v) {
			filteredSlice = append(filteredSlice, v)
		}
	}
	return filteredSlice
}

func (sl Slice) FilterMap(f func(v any) bool, f2 func(v any) any) Slice {
	return sl.Filter(f).Map(f2)
}

func isSlice(v any) bool {
	k := reflect.TypeOf(v).Kind()
	return k == reflect.Slice
}

func isArray(v any) bool {
	k := reflect.TypeOf(v).Kind()
	return k == reflect.Array
}

func (sl *Slice) Flatten() {
	if sl == nil {
		return
	}
	var flattenedSlice Slice
	for _, val := range *sl {
		switch {
		case isSlice(val) || isArray(val):
			flattenedSlice = append(flattenedSlice, val.(Slice)...)
		default:
			flattenedSlice = append(flattenedSlice, val)
		}
	}

	*sl = flattenedSlice
}

func (sl *Slice) FlattenDepth(depth int) {
	if sl == nil {
		return
	}
	for i := 0; i < depth; i++ {
		var flattenedSlice Slice
		flattened := false
		for _, val := range *sl {
			switch {
			case isSlice(val) || isArray(val):
				flattened = true
				flattenedSlice = append(flattenedSlice, val.(Slice)...)
			default:
				flattenedSlice = append(flattenedSlice, val)
			}
		}
		if !flattened {
			return
		}
		*sl = flattenedSlice
	}
}

func (sl *Slice) FlattenFully() {
	if sl == nil {
		return
	}
	for {
		var flattenedSlice Slice
		flattened := false
		for _, val := range *sl {
			switch {
			case isSlice(val) || isArray(val):
				flattened = true
				flattenedSlice = append(flattenedSlice, val.(Slice)...)
			default:
				flattenedSlice = append(flattenedSlice, val)
			}
		}
		if !flattened {
			return
		}
		*sl = flattenedSlice
	}
}

func (sl *Slice) Insert(v any, n int) {
	for len(*sl) < n {
		*sl = append(*sl, nil) // Safety first bro
	}
	rightSide := append(Slice{v}, (*sl)[n:]...)
	*sl = append((*sl)[:n], rightSide...)
}

func (sl *Slice) Reverse() {
	for i, j := 0, len(*sl)-1; i < j; i, j = i+1, j-1 {
		(*sl)[i], (*sl)[j] = (*sl)[j], (*sl)[i]
	}

}
