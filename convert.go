package sliceutils

import "reflect"

// E is used to convert from Slice[T] to Slice[E]
type E Value[any]

// All functions below make use of E

// # Flatten
//
// Removes one layer of nested structure.
//
//	Slice[Slice[T]]Flatten() return Slice[T]
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
			result.Push(v)
		}
	}
	return result
}

// # FlattenN
//
// Removes n layers of nested structure.
//
//	Slice[Slice[Slice[T]]]FlattenN(2) return Slice[T]
func (sl Slice[T]) FlattenN(n uint) Slice[E] {
	if sl.IsNested() && n > 0 {
		var result Slice[E]
		for _, v := range sl {
			result.Push(v)
		}
		return result
	} else {
		return sl.Flatten().FlattenN(n - 1)
	}
}

// # FlattenAll
//
//	Removes "all" layers of nested structure.
//
//	Slice[Slice[Slice[...]]]FlattenAll() return Slice[T]
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

func (sl Slice[T]) FlatMap(f func(v T) T) Slice[E] {
	mapped := New[T]()
	for _, v := range sl {
		mapped.Push(f(v))
	}
	return mapped.Flatten()
}

// # Split
//
//	Split the slice based on separator sep.
//
//	[1, 2, 3, 4, 5]Split(3) return [[1, 2], [4, 5]]
func (sl Slice[T]) Split(sep T) Slice[E] {
	var sp Slice[E]
	var buf Slice[T]
	if !sl.Contains(sep) {
		for _, val := range sl {
			sp.Push(val)
		}
		return sp
	}
	for _, value := range sl {
		if value.Eq(sep) {
			sp.Push(buf)
			buf.Clear()
		} else {
			buf.Push(value)
		}
	}
	sp = append(sp, buf)
	return sp
}

// # SplitN
//
//	Split the slice into n partitions based on separator sep.
//
//	[1,3,2,3,4,3,5]SplitN(3, 3) return [[1],[2],[4,3,5]]
func (sl Slice[T]) SplitN(n uint, sep T) Slice[E] {
	if n <= 1 || !sl.Contains(sep) {
		return sl.Split(sep)
	}
	var sp Slice[E]
	var buf Slice[T]
	for i, value := range sl {
		if value.Eq(sep) {
			sp.Push(buf)
			buf.Clear()
			n--
			if n <= 1 {
				return append(sp, sl[i:])
			}
		} else {
			buf.Push(value)
		}
	}
	sp = append(sp, buf)
	return sp
}

// # SplitOnce
//
//	Split the slice into two partitions on first occurance of sep.
//
//	[1, 3, 2, 3, 4, 5]SplitOnce(3) return [[1], [2, 3, 4, 5]]
func (sl Slice[T]) SplitOnce(sep T) Slice[E] {
	return sl.SplitN(2, sep)
}

// # SplitBy
//
//	Split the slice into two partitions on first occurance of sep.
//
//		[1,2,3,4,5]SplitBy(func(v T) bool {return v%2==0}) {
//	 	return [[1], [3], [5]]
//		}
func (sl Slice[T]) SplitBy(f func(v T) bool) Slice[E] {
	var sp Slice[E]
	var buf Slice[T]
	for _, value := range sl {
		if f(value) {
			sp.Push(buf)
			buf.Clear()
		} else {
			buf.Push(value)
		}
	}
	sp = append(sp, buf)
	if len(sp) == 1 {
		return sp.Flatten()
	}
	return sp
}

// # Chunk
//
// Create a new slice of non-overlapping chunks with the given size.
//
//	[1,2,3,4]Chunk(2) return [[1,2],[3,4]]
//	[1,2,3,4]Chunk(3) return [[1,2,3],[4]]
//	[1,2,3,4]Chunk(4) return [1,2,3,4]
//
// # Caution!
//
// Panics if size is 0
func (sl Slice[T]) Chunk(size uint) Slice[E] {
	if size == 0 {
		panic("chunk size cannot be 0")
	}
	if sl.IsEmpty() {
		return New[E]()
	}
	chunks := New[E]()
	var chunk Slice[T]

	for i, v := range sl {
		chunk.Push(v)

		// If the chunk size is reached or it's the last element, add the chunk to the result
		if uint(i+1)%size == 0 || i == sl.Len()-1 {
			chunks.Push(chunk)
			chunk.Clear()
		}
	}
	if len(chunks) == 1 {
		return chunks.Flatten()
	}

	return chunks
}

// # ChunkBy
//
// Create a new slice of non-overlapping chunks by the given function.
func (sl Slice[T]) ChunkBy(f func(T, T) bool) Slice[E] {
	if sl.IsEmpty() {
		return New[E]()
	}

	chunks := New[E]()
	var chunk Slice[T]

	for i, v := range sl {
		if i > 0 && !f(sl.Get(i-1), v) {
			chunks.Push(chunk)
			chunk.Clear()
		}
		chunk.Push(v)
	}
	chunks.Push(chunk)

	if len(chunks) == 1 {
		return chunks.Flatten()
	}

	return chunks
}

// # Windows
//
// Create overlapping windows of given size.
//
//	[1,2,3,4]Windows(2) return [[1,2], [2,3], [3,4]]
//	[1,2,3,4]Windows(3) return [[1,2,3], [2,3,4]]
//	[1,2,3,4]Windows(4) return [1,2,3,4]
//
// # Caution!
//
// Panics if size is 0
func (sl Slice[T]) Windows(size uint) Slice[E] {
	if size == 0 {
		panic("size of windows cannot be 0")
	}
	var windows Slice[E]
	if sl.Len() < int(size) {
		for _, value := range sl {
			windows.Push(value)
		}
		return windows
	}
	for i := 0; i < sl.Len()-int(size)+1; i++ {
		windows.Push(sl[i : i+int(size)])
	}
	return windows
}
