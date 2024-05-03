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

// Represents a slice value. It needs to implement Eq and Ord.
type Value[T any] interface {
	Eq[T]
	Ord[T]
}

// The generic Slice type. Needs to implement the Value interface.
type Slice[T Value[any]] []T

// Interface for implementing your own default value.
type Default[T Value[any]] interface {
	Default() T
}

// Default implementation for Slice[T]. Returns the default value of type T.
func (sl Slice[T]) Default() T {
	var out T
	return out
}

// # New
//
// Create a new Slice of type T
func New[T Value[any]](v ...T) Slice[T] {
	return v
}

// # Pop
//
// Remove the last element of the slice and return the value
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

// # PopFront
//
// Remove the first element of the slice and return the value
func (sl *Slice[T]) PopFront() (T, error) {
	if sl.IsEmpty() {
		return sl.Default(), ErrIsEmpty
	}
	if sl == nil {
		return sl.Default(), ErrIsNil
	}
	firstElement := (*sl)[0]
	if sl.Len() == 1 {
		*sl = New[T]()
		return firstElement, nil
	}

	*sl = (*sl)[1:]
	return firstElement, nil
}

// # Remove
//
// Remove the element at index n
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

// # Push
//
// Add value(s) at the end of the slice
func (sl *Slice[T]) Push(v ...T) error {
	if sl == nil {
		return ErrIsNil
	}
	*sl = append(*sl, v...)
	return nil
}

// # PushFront
//
// Add value(s) at the start of the slice
func (sl *Slice[T]) PushFront(v ...T) error {
	if sl == nil {
		return ErrIsNil
	}
	*sl = append(v, *sl...)
	return nil
}

// # Insert
//
// Add value(s) at index n, shifting all of the values after n to the right
func (sl *Slice[T]) Insert(n int, v ...T) error {
	if sl.Len() < n {
		return ErrOutOfRange
	}
	if sl == nil {
		return ErrIsNil
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
		result = append(copy, v...)
	} else {
		for i, val := range copy {
			if i == n {
				result.Push(v...)
				result.Push(val)
			} else {
				result.Push(val)
			}
		}
	}
	*sl = result
	return nil
}

// # Count
//
// Return the total amount of occurances of v
func (sl Slice[T]) Count(v any) int {
	count := 0
	for _, val := range sl {
		if val.Eq(v) {
			count++
		}
	}
	return count
}

// # DeepCount
//
// Flatten all nested structure, and return the total amount of occurances of v
func (sl Slice[T]) DeepCount(v any) int {
	return sl.FlattenAll().Count(v)
}

// # CountBy
//
// Return the total amount of occurances where the provided function returns true
func (sl Slice[T]) CountBy(f func(T) bool) int {
	count := 0
	for _, val := range sl {
		if f(val) {
			count++
		}
	}
	return count
}

// # Contains
//
// Returns true if slice contains v
func (sl Slice[T]) Contains(v any) bool {
	for _, val := range sl {
		if val.Eq(v) {
			return true
		}
	}
	return false
}

// # ContainsMany
//
// Returns true if slice contains all values v
func (sl Slice[T]) ContainsMany(v ...any) bool {
	for _, val := range v {
		if !sl.Contains(val) {
			return false
		}
	}
	return true
}

// # ContainsOnly
//
// Returns true if slice only contains v
func (sl Slice[T]) ContainsOnly(v any) bool {
	for _, val := range sl {
		if !val.Eq(v) {
			return false
		}
	}
	return true
}

// # IsPrefixOf
//
// Returns true if slice sl is prefix of other
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

// # ForEach
//
// Loop through all elements in the slice and apply a provided function to the value
func (sl Slice[T]) ForEach(f func(T)) {
	for _, v := range sl {
		f(v)
	}
}

// # Map
//
// Loop through all elements in the slice and apply a provided function to the value. Return the result.
func (sl Slice[T]) Map(f func(v T) T) Slice[T] {
	var mappedSlice Slice[T]
	for _, v := range sl {
		mappedSlice.Push(f(v))
	}
	return mappedSlice
}

// # MapWhile
//
// Same as Map, but only change the value if the returned value of the provided function is not nil.
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

// # StepBy
//
// Return a Slice that starts at index 0, and only contains elements with n steps in between
func (sl Slice[T]) StepBy(n uint) Slice[T] {
	step := New[T]()
	for i := 0; i < sl.Len(); i += int(n) {
		step.Push(sl[i])
	}
	return step
}

// # Filter
//
// Apply a filter function on all elements and return a slice of all elements that returned true
func (sl Slice[T]) Filter(f func(v T) bool) Slice[T] {
	var filteredSlice Slice[T]
	for _, v := range sl {
		if f(v) {
			filteredSlice = append(filteredSlice, v)
		}
	}
	return filteredSlice
}

// # FilterMap
//
// First apply Filter, and then Map. Return the result.
func (sl Slice[T]) FilterMap(f func(v T) bool, f2 func(v T) T) Slice[T] {
	return sl.Filter(f).Map(f2)
}

// # IsNested
//
// Returns true if the slice contains any type of nested structure
func (sl Slice[T]) IsNested() bool {
	if sl.IsEmpty() {
		return false
	}
	return reflect.TypeOf(sl[0]).Kind() == reflect.Slice
}

// # Get
//
// Get the value at index n without modifying the slice
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

// # GetRange
//
// Same as Get, but specify a range to get from.
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

// # Repeat
//
// Create a new slice that is the provided slice repeated n amount of times. Return the new slice.
func (sl Slice[T]) Repeat(n uint) Slice[T] {
	copy := sl
	for ; n > 1; n-- {
		sl = append(sl, copy...)
	}
	return sl
}

// # Reverse
//
// Create a copy of the slice and reverse it. Return the reversed slice.
func (sl Slice[T]) Reverse() Slice[T] {
	var rev Slice[T]
	for i := sl.Len() - 1; i >= 0; i-- {
		rev.Push(sl[i])
	}
	return rev
}

// # Concat
//
// Concatenate two slices into one slice.
func (sl Slice[T]) Concat(sl2 Slice[T]) Slice[T] {
	return append(sl, sl2...)
}

// # Join
//
// Join two slices together with the separator(s) sep.
func (sl Slice[T]) Join(sl2 Slice[T], sep ...T) Slice[T] {
	sl = append(sl, sep...)
	return append(sl, sl2...)
}

// # Min
//
// Return the minimum value of the slice.
func (sl Slice[T]) Min() (T, error) {
	min, err := sl.First()
	for _, v := range sl {
		if v.Lt(min) {
			min = v
		}
	}
	return min, err
}

// # Max
//
// Return the maximum value of the slice.
func (sl Slice[T]) Max() (T, error) {
	min, err := sl.First()
	for _, v := range sl {
		if v.Gt(min) {
			min = v
		}
	}
	return min, err
}

// # MaxBy
//
// Return the maximum value of the slice based on the function f.
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

// # First
//
// Return the value at index 0.
func (sl Slice[T]) First() (T, error) {
	return sl.Get(0)
}

// # Last
//
// Return the value at the last index
func (sl Slice[T]) Last() (T, error) {
	return sl.Get(sl.Len() - 1)
}

// # IndexIS
//
// Returns true if the element at index n is the same as value
func (sl Slice[T]) IndexIs(n int, value T) bool {
	if n < 0 {
		for n < 0 {
			n += sl.Len()
		}
		n++
	}
	return sl[n].Eq(value)
}

// # StartsWith
//
// Returns true if the slice starts with value
func (sl Slice[T]) StartsWith(value T) (bool, error) {
	if firstValue, err := sl.First(); err != nil {
		return false, err
	} else {
		return firstValue.Eq(value), nil
	}
}

// # EndsWith
//
// Returns true if the slice ends with v
func (sl Slice[T]) EndsWith(v T) (bool, error) {
	if lastValue, err := sl.Last(); err != nil {
		return false, err
	} else {
		return lastValue.Eq(v), nil
	}
}

// # Find
//
// Return the first element where the provided function f returns true
func (sl Slice[T]) Find(f func(v T) bool) (T, error) {
	for _, v := range sl {
		if f(v) {
			return v, nil
		}
	}
	return sl.Default(), ErrDoesNotExist
}

// # FindMap
//
// Return the first element where the provided function f returns not nil
func (sl Slice[T]) FindMap(f func(v T) *T) (T, error) {
	for _, v := range sl {
		if f(v) != nil {
			return v, nil
		}
	}
	return sl.Default(), ErrDoesNotExist
}

// # FirstIndexOf
//
// Return the index of first instance of v
func (sl Slice[T]) FirstIndexOf(v T) (int, error) {
	for i, value := range sl {
		if value.Eq(v) {
			return i, nil
		}
	}
	return -1, ErrDoesNotExist
}

// # LastIndexOf
//
// Return the index of last instance of v
func (sl Slice[T]) LastIndexOf(v T) (int, error) {
	for i := sl.Len() - 1; i >= 0; i-- {
		if sl[i].Eq(v) {
			return i, nil
		}
	}
	return -1, ErrDoesNotExist
}

// # AllIndexesOf
//
// Return the index of all instances of v
func (sl Slice[T]) AllIndexesOf(v T) (Slice[Int], error) {
	indexes := New[Int]()
	if !sl.Contains(v) {
		return indexes, ErrDoesNotExist
	}
	for i, v := range sl {
		if v.Eq(v) {
			indexes.Push(Int(i))
		}
	}
	return indexes, nil
}

type V Value[any]

// # Fold
//
// Apply function f on all elements of the slice and accumulate them into one value
func (sl Slice[T]) Fold(init V, f func(V, T) V) V {
	for _, v := range sl {
		init = f(init, v)
	}
	return init
}

// # Reduce
//
// Works the same as Fold but starts accumulating at the first element of the slice
func (sl Slice[T]) Reduce(f func(acc, v T) T) (T, error) {
	acc, err := sl.First()
	for _, v := range sl.Skip(1) {
		acc = f(acc, v)
	}
	return acc, err
}

// # Skip
//
// Return a new slice where n amount of elements are skipped
func (sl Slice[T]) Skip(n uint) Slice[T] {
	if sl.Len() < int(n) {
		return New[T]()
	}
	return sl[n:]
}

// # SkipWhile
//
// Return a new slice where all elements until the provided function f returns false are skipped
func (sl Slice[T]) SkipWhile(f func(T) bool) Slice[T] {
	for i, v := range sl {
		if !f(v) {
			return sl[i:]
		}
	}
	return sl
}

// # Take
//
// Return a new slice where only n amount of elements are retained
func (sl Slice[T]) Take(n uint) Slice[T] {
	if sl.Len() < int(n) {
		return sl
	}
	return sl[:n]
}

// # TakeWhile
//
// Return a new slice where all elements until f returns false are retained
func (sl Slice[T]) TakeWhile(f func(T) bool) Slice[T] {
	for i, v := range sl {
		if !f(v) {
			return sl[:i]
		}
	}
	return sl
}

// # Zip
//
// Return a new slice where two slices are zipped together into one slice
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

// # All
//
// Return true if function f returns true on all elements of the slice
func (sl Slice[T]) All(f func(v T) bool) bool {
	for _, v := range sl {
		if !f(v) {
			return false
		}
	}
	return true
}

// # Any
//
// Return true if function f returns true on any element of the slice
func (sl Slice[T]) Any(f func(v T) bool) bool {
	for _, v := range sl {
		if f(v) {
			return true
		}
	}
	return false
}

// # Enumerate
//
// Works like ForEach, but also provides the index
func (sl Slice[T]) Enumerate(f func(int, T)) {
	for i, v := range sl {
		f(i, v)
	}
}

// # Copy
//
// Returns a copy of the slice.
func (sl Slice[T]) Copy() Slice[T] {
	return New(sl...)
}

// # Len
//
// Returns the length of the slice
func (sl Slice[T]) Len() int {
	return len(sl)
}

// # IsEmpty
//
// Returns true if the slice is empty
func (sl Slice[T]) IsEmpty() bool {
	return sl.Len() == 0
}
