package sliceutils

import (
	"reflect"
)

type Slice[T Eq[any]] []T

type Default[T Eq[any]] interface {
	Default() T
}

func (sl Slice[T]) Default() T {
	var out T
	return out
}

func New[T Eq[any]](values ...T) Slice[T] {
	return values
}

func (sl *Slice[T]) Pop() T {
	if len(*sl) == 0 || sl == nil {
		return sl.Default()
	}
	lastIndex := len(*sl) - 1
	lastElement := (*sl)[lastIndex]

	*sl = (*sl)[:lastIndex]
	return lastElement
}

func (sl *Slice[T]) PopFront() T {
	if len(*sl) == 0 || sl == nil {
		return sl.Default()
	}
	firstElement := (*sl)[0]
	if len(*sl) == 1 {
		*sl = New[T]()
		return firstElement
	}

	*sl = (*sl)[1:]
	return firstElement
}

func (sl *Slice[T]) PopN(n int) T {
	if len(*sl) == 0 || len(*sl) < int(n) || sl == nil {
		return sl.Default()
	}
	for n < 0 {
		n += len(*sl)
	}
	copy := *sl
	value := copy[n]
	if n == len(copy)-1 {
		copy = copy[:n]
	} else {
		copy = append(copy[:n], copy[n+1:]...)
	}
	*sl = copy
	return value
}

func (sl *Slice[T]) Push(values ...T) {
	if sl == nil {
		return
	}
	*sl = append(*sl, values...)
}

func (sl *Slice[T]) PushFront(values ...T) {
	if sl == nil {
		return
	}
	*sl = append(values, *sl...)
}

func (sl *Slice[T]) PushN(n int, values ...T) {
	if sl == nil || len(*sl) < n {
		return
	}
	copy := *sl
	if n < 0 {
		for n < 0 {
			n += len(copy)
		}
		n++
	}
	var result Slice[T]
	if n == len(copy) {
		result = append(copy, values...)
	} else {
		for i, val := range copy {
			if i == n {
				result.Push(values...)
				result.Push(val)
			} else {
				result.Push(val)
			}
		}
	}
	*sl = result
}

func (sl *Slice[T]) Clear() {
	*sl = Slice[T]{}
}

func (sl Slice[T]) Count(v any) int {
	count := 0
	for _, val := range sl {
		if val.Eq(v) {
			count++
		}
	}
	return count
}

func (sl Slice[T]) Contains(v any) bool {
	for _, val := range sl {
		if val.Eq(v) {
			return true
		}
	}
	return false
}

func (sl Slice[T]) ForEach(f func(any)) {
	for _, v := range sl {
		f(v)
	}
}

func (sl Slice[T]) Map(f func(v T) T) Slice[T] {
	var mappedSlice Slice[T]
	for _, v := range sl {
		mappedSlice = append(mappedSlice, f(v))
	}
	return mappedSlice
}

func (sl Slice[T]) Filter(f func(v T) bool) Slice[T] {
	var filteredSlice Slice[T]
	for _, v := range sl {
		if f(v) {
			filteredSlice = append(filteredSlice, v)
		}
	}
	return filteredSlice
}

func (sl Slice[T]) FilterMap(f func(v T) bool, f2 func(v T) T) Slice[T] {
	return sl.Filter(f).Map(f2)
}

func (sl Slice[T]) IsNested() bool {
	if len(sl) == 0 {
		return false
	}
	return reflect.TypeOf(sl[0]).Kind() == reflect.Slice
}

type E Eq[any]

func (sl Slice[T]) Flatten() Slice[E] {
	var result Slice[E]
	if sl.IsNested() {
		for _, v := range sl {
			nestedSlice := reflect.ValueOf(v)
			for i := 0; i < nestedSlice.Len(); i++ {
				result = append(result, nestedSlice.Index(i).Interface().(E))
			}
		}
	} else {
		for _, v := range sl {
			result = append(result, v)
		}
	}
	return result
}

func (sl *Slice[T]) Dedup() {
	seen := Slice[T]{}
	for _, value := range *sl {
		if !seen.Contains(value) {
			seen.Push(value)
		}
	}
	*sl = seen
}

func (sl Slice[T]) FlattenN(n uint) Slice[E] {
	if !sl.IsNested() || n == 0 {
		var result Slice[E]
		for _, v := range sl {
			result.Push(v)
		}
		return result
	} else {
		return sl.Flatten().FlattenN(n - 1)
	}
}

func (sl Slice[T]) FlattenAll() Slice[E] {
	if sl.IsNested() {
		return sl.Flatten().FlattenAll()
	} else {
		var result Slice[E]
		for _, v := range sl {
			result.Push(v)
		}
		return result
	}
}

func (sl *Slice[T]) Reverse() {
	for i, j := 0, len(*sl)-1; i < j; i, j = i+1, j-1 {
		(*sl)[i], (*sl)[j] = (*sl)[j], (*sl)[i]
	}
}
