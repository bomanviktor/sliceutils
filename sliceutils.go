package sliceutils

import (
	"errors"
	"reflect"
)

// Type aliases to enable the implementation of Value for the builtin types
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

var (
	ErrDoesNotExist = errors.New("value does not exist in slice")
	ErrIsEmpty      = errors.New("slice is empty")
	ErrIsNil        = errors.New("slice is nil")
	ErrOutOfRange   = errors.New("index is out of range")
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

func (sl *Slice[T]) Pop() (T, error) {
	if sl.IsEmpty() {
		return sl.Default(), ErrIsEmpty
	}
	if sl == nil {
		return sl.Default(), ErrIsNil
	}
	lastIndex := sl.Len() - 1
	lastElement := (*sl)[lastIndex]

	*sl = (*sl)[:lastIndex]
	return lastElement, nil
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

func (sl *Slice[T]) Remove(n uint) (T, error) {
	if sl.IsEmpty() {
		return sl.Default(), ErrIsEmpty
	}

	if sl.Len() < int(n) {
		return sl.Default(), ErrOutOfRange
	}
	if sl == nil {
		return sl.Default(), ErrIsNil
	}
	copy := *sl
	value := copy[n]
	if int(n) == copy.Len()-1 {
		copy = copy[:n]
	} else {
		copy = append(copy[:n], copy[n+1:]...)
	}
	*sl = copy
	return value, nil
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

func (sl *Slice[T]) Insert(n int, values ...T) error {
	if sl == nil || sl.Len() < n {
		return ErrOutOfRange
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
	return nil
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

func (sl Slice[T]) ContainsMany(v ...any) bool {
	for _, val := range v {
		if !sl.Contains(val) {
			return false
		}
	}
	return true
}

func (sl Slice[T]) ContainsAll(v any) bool {
	for _, val := range sl {
		if !val.Eq(v) {
			return false
		}
	}
	return true
}

func (sl Slice[T]) IsPrefixOf(other Slice[T]) bool {
	if sl.Len() >= other.Len() {
		return false
	}
	for i := 0; i < sl.Len(); i++ {
		if !sl[i].Eq(other[i]) {
			return false
		}
	}
	return true
}

func (sl Slice[T]) ForEach(f func(T)) {
	for _, v := range sl {
		f(v)
	}
}

func (sl Slice[T]) Map(f func(v T) T) Slice[T] {
	var mappedSlice Slice[T]
	for _, v := range sl {
		mappedSlice.Push(f(v))
	}
	return mappedSlice
}

func (sl Slice[T]) MapWhile(f func(v T) *T) Slice[T] {
	var mappedSlice Slice[T]
	for _, v := range sl {
		if f(v) != nil {
			mappedSlice.Push(*f(v))
		} else {
			mappedSlice.Push(v)
		}
	}
	return mappedSlice
}

func (sl Slice[T]) StepBy(n uint) Slice[T] {
	step := New[T]()
	for i := 0; i < sl.Len(); i += int(n) {
		step.Push(sl[i])
	}
	return step
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

func (sl Slice[T]) Get(n int) (T, error) {
	if sl.IsEmpty() {
		return sl.Default(), ErrIsEmpty
	}
	if sl.Len() < n {
		return sl.Default(), ErrOutOfRange
	}
	if n < 0 {
		for n < 0 {
			n += sl.Len()
		}
		n++
	}
	return sl[n], nil
}

func (sl Slice[T]) GetRange(from, to int) (Slice[T], error) {
	r := New[T]()
	if sl.IsEmpty() {
		return r, ErrIsEmpty
	}
	if sl.Len() < from || sl.Len() < to {
		return r, ErrOutOfRange
	}
	for ; from < to; from++ {
		r.Push(sl[from])
	}
	return r, nil
}

func (sl Slice[T]) Repeat(n uint) Slice[T] {
	copy := sl
	for ; n > 1; n-- {
		sl = append(sl, copy...)
	}
	return sl
}

func (sl Slice[T]) Reverse() Slice[T] {
	var rev Slice[T]
	for i := sl.Len() - 1; i >= 0; i-- {
		rev.Push(sl[i])
	}
	return rev
}

func (sl Slice[T]) Concat(sl2 Slice[T]) Slice[T] {
	return append(sl, sl2...)
}

func (sl Slice[T]) Join(sl2 Slice[T], sep ...T) Slice[T] {
	sl = append(sl, sep...)
	return append(sl, sl2...)
}

func (sl Slice[T]) Min() (T, error) {
	min, err := sl.First()
	for _, v := range sl {
		if v.Lt(min) {
			min = v
		}
	}
	return min, err
}

func (sl Slice[T]) Max() (T, error) {
	min, err := sl.First()
	for _, v := range sl {
		if v.Gt(min) {
			min = v
		}
	}
	return min, err
}

func (sl Slice[T]) MaxBy(f func(T) T) (T, error) {
	if sl.IsEmpty() {
		return sl.Default(), ErrIsEmpty
	}
	max := sl.Default()
	for _, v := range sl {
		if f(v).Gt(max) {
			max = v
		}
	}
	return max, nil
}

func (sl Slice[T]) First() (T, error) {
	return sl.Get(0)
}

func (sl Slice[T]) Last() (T, error) {
	return sl.Get(sl.Len() - 1)
}

func (sl Slice[T]) IndexIs(n int, value T) bool {
	if n < 0 {
		for n < 0 {
			n += sl.Len()
		}
		n++
	}
	return sl[n].Eq(value)
}

func (sl Slice[T]) StartsWith(value T) (bool, error) {
	if firstValue, err := sl.First(); err != nil {
		return false, err
	} else {
		return firstValue.Eq(value), nil
	}
}

func (sl Slice[T]) EndsWith(value T) (bool, error) {
	if lastValue, err := sl.Last(); err != nil {
		return false, err
	} else {
		return lastValue.Eq(value), nil
	}
}

func (sl Slice[T]) Find(f func(v T) bool) (T, error) {
	for _, v := range sl {
		if f(v) {
			return v, nil
		}
	}
	return sl.Default(), ErrDoesNotExist
}

func (sl Slice[T]) FindMap(f func(v T) *T) (T, error) {
	for _, v := range sl {
		if f(v) != nil {
			return v, nil
		}
	}
	return sl.Default(), ErrDoesNotExist
}

func (sl Slice[T]) FirstIndexOf(v T) (int, error) {
	for i, value := range sl {
		if value.Eq(v) {
			return i, nil
		}
	}
	return -1, ErrDoesNotExist
}

func (sl Slice[T]) LastIndexOf(v T) (int, error) {
	for i := sl.Len() - 1; i >= 0; i-- {
		if sl[i].Eq(v) {
			return i, nil
		}
	}
	return -1, ErrDoesNotExist
}

func (sl Slice[T]) AllIndexesOf(value T) (Slice[Int], error) {
	indexes := New[Int]()
	if !sl.Contains(value) {
		return indexes, ErrDoesNotExist
	}
	for i, v := range sl {
		if v.Eq(value) {
			indexes.Push(Int(i))
		}
	}
	return indexes, nil
}

type V Value[any]

func (sl Slice[T]) Fold(init V, f func(V, T) V) V {
	for _, v := range sl {
		init = f(init, v)
	}
	return init
}

func (sl Slice[T]) Reduce(f func(acc, v T) T) (T, error) {
	acc, err := sl.First()
	for _, v := range sl.Skip(1) {
		acc = f(acc, v)
	}
	return acc, err
}

func (sl Slice[T]) Skip(n uint) Slice[T] {
	if sl.Len() < int(n) {
		return New[T]()
	}
	return sl[n:]
}

func (sl Slice[T]) SkipWhile(f func(T) bool) Slice[T] {
	for i, v := range sl {
		if !f(v) {
			return sl[i:]
		}
	}
	return sl
}

func (sl Slice[T]) Take(n uint) Slice[T] {
	if sl.Len() < int(n) {
		return sl
	}
	return sl[:n]
}

func (sl Slice[T]) TakeWhile(f func(T) bool) Slice[T] {
	for i, v := range sl {
		if !f(v) {
			return sl[:i]
		}
	}
	return sl
}

func (sl Slice[T]) Zip(other Slice[T]) Slice[T] {
	zipped := New[T]()
	if sl.Len() > other.Len() {
		for i := 0; i < other.Len(); i++ {
			zipped.Push(sl[i], other[i])
		}
		zipped.Push(sl[other.Len():]...)
		return zipped
	} else if sl.Len() < other.Len() {
		for i := 0; i < sl.Len(); i++ {
			zipped.Push(sl[i], other[i])
		}
		zipped.Push(other[sl.Len():]...)
		return zipped
	} else {
		for i := 0; i < sl.Len(); i++ {
			zipped.Push(sl[i], other[i])
		}
		return zipped
	}
}

func (sl Slice[T]) All(f func(v T) bool) bool {
	for _, v := range sl {
		if !f(v) {
			return false
		}
	}
	return true
}

func (sl Slice[T]) Any(f func(v T) bool) bool {
	for _, v := range sl {
		if f(v) {
			return true
		}
	}
	return false
}

func (sl Slice[T]) Enumerate(f func(int, T)) {
	for i, v := range sl {
		f(i, v)
	}
}

func (sl Slice[T]) Copy() Slice[T] {
	return New(sl...)
}

func (sl Slice[T]) Len() int {
	return len(sl)
}

func (sl Slice[T]) IsEmpty() bool {
	return sl.Len() == 0
}
