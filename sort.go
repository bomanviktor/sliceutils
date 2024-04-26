package sliceutils

func (sl Slice[T]) quickSort(low, high int) {
	if low < high {
		pivot := sl.partition(low, high)
		sl.quickSort(low, pivot-1)
		sl.quickSort(pivot+1, high)
	}
}

func (sl Slice[T]) partition(low, high int) int {
	pivot := sl[high]
	i := low - 1
	for j := low; j < high; j++ {
		if sl[j].Lt(pivot) {
			i++
			sl[i], sl[j] = sl[j], sl[i]
		}
	}
	sl[i+1], sl[high] = sl[high], sl[i+1]
	return i + 1
}

// # Sort
//
// Sorts the slice in place using quicksort.
//
//	[1,4,3,5,2]Sort() return [1,2,3,4,5]
func (sl Slice[T]) Sort() {
	sl.quickSort(0, sl.Len()-1)
}

func (sl Slice[T]) mergeSort(f func(v1 T, v2 T) bool, left, right int) {
	if left < right {
		mid := (left + right) / 2
		sl.mergeSort(f, left, mid)
		sl.mergeSort(f, mid+1, right)
		sl.merge(f, left, mid, right)
	}
}

func (sl Slice[T]) merge(f func(v1 T, v2 T) bool, left, mid, right int) {
	temp := make(Slice[T], right-left+1)
	i, j, k := left, mid+1, 0

	for i <= mid && j <= right {
		if f(sl[i], sl[j]) {
			temp[k] = sl[i]
			i++
		} else {
			temp[k] = sl[j]
			j++
		}
		k++
	}

	for i <= mid {
		temp[k] = sl[i]
		i++
		k++
	}

	for j <= right {
		temp[k] = sl[j]
		j++
		k++
	}

	for i, val := range temp {
		sl[left+i] = val
	}
}

// # SortBy
//
// Sorts the slice in place by the result of the given function.
func (sl Slice[T]) SortBy(f func(v1 T, v2 T) bool) {
	if sl.Len() <= 1 {
		return
	}
	sl.mergeSort(f, 0, sl.Len()-1)
}

// # IsSorted
//
// Returns true if the slice is sorted.
//
//	[1,2,3]IsSorted() return true
//	[3,2,1]IsSorted() return false
func (sl Slice[T]) IsSorted() bool {
	for i := 0; i < sl.Len()-1; i++ {
		if sl.Get(i).Lt(sl.Get(i + 1)) {
			return false
		}
	}
	return true
}

// # IsSortedBy
//
// # Returns true if the slice is sorted according to the function f
//
//	[1,2,3]IsSortedBy(func(v1 == v2 - 1)) return true
//	[1,2,3]IsSortedBy(func(v1 == v2 - 2)) return false
func (sl Slice[T]) IsSortedBy(f func(v1, v2 T) bool) bool {
	for i := 0; i < sl.Len()-1; i++ {
		if !f(sl.Get(i), sl.Get(i+1)) {
			return false
		}
	}
	return true
}
