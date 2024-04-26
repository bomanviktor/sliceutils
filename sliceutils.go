package sliceutils

import (
	"reflect"
)

// Type aliases to enable the implementation of Val
type (
	Bool bool
	Str  string
	Rune rune
	Int  int
	I8   int8
	I16  int16
	I32  int32
	I64  int64
	Uint uint
	U8   uint8
	U16  uint16
	U32  uint32
	U64  uint64
	Byte byte
	F32  float32
	F64  float64
	C64  complex64
	C128 complex128
)

type Value[T any] interface {
	Eq[T]
	Ord[T]
}

type Slice[T Value[any]] []T

type Default[T Value[any]] interface {
	Default() T
}

func (sl Slice[T]) Default() T {
	var out T
	return out
}

func New[T Value[any]](values ...T) Slice[T] {
	return values
}

func (sl *Slice[T]) Pop() T {
	if sl.IsEmpty() || sl == nil {
		return sl.Default()
	}
	lastIndex := sl.Len() - 1
	lastElement := (*sl)[lastIndex]

	*sl = (*sl)[:lastIndex]
	return lastElement
}

func (sl *Slice[T]) PopFront() T {
	if sl.IsEmpty() || sl == nil {
		return sl.Default()
	}
	firstElement := (*sl)[0]
	if sl.Len() == 1 {
		*sl = New[T]()
		return firstElement
	}

	*sl = (*sl)[1:]
	return firstElement
}

func (sl *Slice[T]) PopN(n int) T {
	if sl.IsEmpty() || sl.Len() < int(n) || sl == nil {
		return sl.Default()
	}
	for n < 0 {
		n += sl.Len()
	}
	copy := *sl
	value := copy[n]
	if n == copy.Len()-1 {
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
	if sl == nil || sl.Len() < n {
		return
	}
	copy := *sl
	if n < 0 {
		for n < 0 {
			n += copy.Len()
		}
		n++
	}
	var result Slice[T]
	if n == copy.Len() {
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
	*sl = New[T]()
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

func (sl Slice[T]) DeepCount(v any) int {
	count := 0
	for _, val := range sl.FlattenAll() {
		if val.Eq(v) {
			count++
		}
	}
	return count
}

func (sl Slice[T]) CountFunc(f func(T) bool) int {
	count := 0
	for _, val := range sl {
		if f(val) {
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
	if sl.IsEmpty() {
		return false
	}
	return reflect.TypeOf(sl[0]).Kind() == reflect.Slice
}

func (sl Slice[T]) Get(n int) T {
	if sl.Len() < n {
		return sl.Default()
	}
	if n < 0 {
		for n < 0 {
			n += sl.Len()
		}
		n++
	}
	return sl[n]
}

func (sl Slice[T]) GetRange(from, to int) Slice[T] {
	chunk := New[T]()
	if sl.Len() < from || sl.Len() < to {
		return chunk
	}
	for ; from < to; from++ {
		chunk.Push(sl.Get(from))
	}
	return chunk
}

func (sl Slice[T]) Set(n int, value T) {
	if sl.Len() < n {
		return
	}
	if n < 0 {
		for n < 0 {
			n += sl.Len()
		}
		n++
	}
	sl[n] = value
}

func (sl Slice[T]) Replace(n uint, value T) T {
	i := int(n)
	if sl.Len() < i {
		return sl.Default()
	}
	swappedValue := sl[i]
	sl[i] = value
	return swappedValue
}

func (sl Slice[T]) Swap(x uint, y uint) {
	if sl.Len() < int(x) || sl.Len() < int(y) {
		return
	}
	sl[x], sl[y] = sl[y], sl[x]
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

func (sl Slice[T]) Repeat(n uint) Slice[T] {
	copy := sl
	for ; n > 1; n-- {
		sl = append(sl, copy...)
	}
	return sl
}

func (sl Slice[T]) Concat(sl2 Slice[T]) Slice[T] {
	return append(sl, sl2...)
}

func (sl Slice[T]) Join(sl2 Slice[T], sep ...T) Slice[T] {
	sl = append(sl, sep...)
	return append(sl, sl2...)
}

func (sl Slice[T]) First() T {
	return sl.Get(0)
}

func (sl Slice[T]) Last() T {
	return sl.Get(sl.Len() - 1)
}

func (sl Slice[T]) IndexIs(n int, value T) bool {
	if n < 0 {
		for n < 0 {
			n += sl.Len()
		}
		n++
	}
	return sl.Get(n).Eq(value)
}

func (sl Slice[T]) StartsWith(value T) bool {
	return sl.First().Eq(value)
}

func (sl Slice[T]) EndsWith(value T) bool {
	return sl.Last().Eq(value)
}

func (sl Slice[T]) Rev() Slice[T] {
	var rev Slice[T]
	for i := sl.Len() - 1; i >= 0; i-- {
		rev.Push(sl[i])
	}
	return rev
}

func (sl Slice[T]) RevMut() Slice[T] {
	for i, j := 0, sl.Len()-1; i < j; i, j = i+1, j-1 {
		sl[i], sl[j] = sl[j], sl[i]
	}
	return sl
}

func (sl Slice[T]) FirstIndexOf(v T) int {
	for i, value := range sl {
		if value.Eq(v) {
			return i
		}
	}
	return -1
}

func (sl Slice[T]) LastIndexOf(v T) int {
	for i := sl.Len() - 1; i >= 0; i-- {
		if sl.Get(i).Eq(v) {
			return i
		}
	}
	return -1
}

func (sl Slice[T]) Len() int {
	return len(sl)
}

func (sl Slice[T]) Fill(value T) {
	for i := 0; i < sl.Len()-1; i++ {
		sl[i] = value
	}
}

func (sl Slice[T]) FillWith(f func() T) {
	for i := 0; i < sl.Len()-1; i++ {
		sl[i] = f()
	}
}

func (sl Slice[T]) FillWithDefault(f func() T) {
	for i := 0; i < sl.Len()-1; i++ {
		sl[i] = sl.Default()
	}
}

func (sl Slice[T]) IsEmpty() bool {
	return sl.Len() == 0
}

func (sl *Slice[T]) RotateLeft(moves uint) {
	length := uint(sl.Len())
	if length == 0 || moves == 0 {
		return
	}

	copy := *sl
	moves %= length
	*sl = append(copy[moves:], copy[:moves]...)
}

func (sl *Slice[T]) RotateRight(moves uint) {
	length := uint(sl.Len())
	if length == 0 || moves == 0 {
		return
	}
	copy := *sl
	moves %= length
	*sl = append(copy[length-moves:], copy[:length-moves]...)
}
