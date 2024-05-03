package sliceutils

// modify
//
// Functions for modifying the underlying slice.

// # Clear
//
// Remove all the values in the slice but keeps the same type.
//
//	Slice[Int]{1,2,3}Clear() return Slice[Int]{}
func (sl *Slice[T]) Clear() {
	*sl = New[T]()
}

// # Dedup
//
// Remove duplicate values in the slice.
//
//	[1,1,2,2,3,3]Dedup() return [1,2,3]
func (sl *Slice[T]) Dedup() {
	seen := Slice[T]{}
	for _, value := range *sl {
		if !seen.Contains(value) {
			seen.Push(value)
		}
	}
	*sl = seen
}

// # Fill
//
// Replace all the values with value of type T.
//
//	[1,2,3]Fill(1) return [1,1,1]
func (sl Slice[T]) Fill(value T) {
	for i := 0; i < sl.Len(); i++ {
		sl[i] = value
	}
}

// # FillWith
//
// Replace all values with the value returned by func f
//
//	[1,2,3]FillWith({return 1}) return [1,1,1]
func (sl Slice[T]) FillWith(f func() T) {
	for i := 0; i < sl.Len()-1; i++ {
		sl[i] = f()
	}
}

// # FillWithDefault
//
// Replace all values with the default value of T
//
//	[1,2,3]FillWithDefault() return [0,0,0]
func (sl Slice[T]) FillWithDefault(f func() T) {
	for i := 0; i < sl.Len()-1; i++ {
		sl[i] = sl.Default()
	}
}

// # Replace
//
// Replace the value at index n and return the replaced value
//
//	[1,2,3]Replace(1, 5) -> [1,5,3] and return 2
func (sl Slice[T]) Replace(n uint, value T) T {
	i := int(n)
	if sl.Len() < i {
		return sl.Default()
	}
	swappedValue := sl[i]
	sl[i] = value
	return swappedValue
}

// # ReverseMut
//
// Reverse the slice and return the result
//
//	[1,2,3]RevMut() -> [3,2,1] and return [3,2,1]
func (sl Slice[T]) ReverseMut() Slice[T] {
	for i, j := 0, sl.Len()-1; i < j; i, j = i+1, j-1 {
		sl[i], sl[j] = sl[j], sl[i]
	}
	return sl
}

// # RotateLeft
//
// Rotate all elements to the left n steps
//
//	[1,2,3]RotateLeft(1) -> [2,3,1]
//	[1,2,3]RotateLeft(2) -> [3,1,2]
//	[1,2,3]RotateLeft(3) -> [1,2,3]
func (sl *Slice[T]) RotateLeft(n uint) {
	length := uint(sl.Len())
	if length == 0 || n == 0 {
		return
	}

	copy := *sl
	n %= length
	*sl = append(copy[n:], copy[:n]...)
}

// # RotateRight
//
// Rotate all elements to the right n steps
//
//	[1,2,3]RotateRight(1) -> [3.1.2]
//	[1,2,3]RotateRight(2) -> [2,3,1]
//	[1,2,3]RotateRight(3) -> [1,2,3]
func (sl *Slice[T]) RotateRight(n uint) {
	length := uint(sl.Len())
	if length == 0 || n == 0 {
		return
	}
	copy := *sl
	n %= length
	*sl = append(copy[length-n:], copy[:length-n]...)
}

// # Set
//
// Set the value at index n with value T.
//
//	[1,2,3]Set(1,5) -> [1,5,3]
//	['a','b','c']Set(1,'z') -> ['a','z','c']
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

// # Swap
//
// Swap element n1 with element n2 in the slice
//
//	[1,2,3,4,5]Swap(1,3) -> [1,4,3,2,5]
func (sl Slice[T]) Swap(n1, n2 uint) {
	if sl.Len() < int(n1) || sl.Len() < int(n2) {
		return
	}
	sl[n1], sl[n2] = sl[n2], sl[n1]
}

// # SwapValues
//
// Swap all values v1 with v2 in the slice
//
//	[1,1,2,2,3,3]Swap(1,3) -> [3,3,2,2,1,1]
func (sl Slice[T]) SwapValues(v1, v2 T) {
	if !sl.Contains(v1) || sl.Contains(v2) {
		return
	}
	for i, value := range sl {
		if value.Eq(v1) {
			sl[i] = v2
		} else if value.Eq(v2) {
			sl[i] = v1
		}
	}
}

// # Purge
//
// Remove all elements of the slice and return the copy.
//
//	[1,2,3]Purge() -> [] and return [1,2,3]
func (sl *Slice[T]) Purge() Slice[T] {
	copy := sl.Copy()
	*sl = New[T]()
	return copy
}
