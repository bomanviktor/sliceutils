package sliceutils

// # Ordering
//
// int8 type alias used in Compare, Gt, Lt, Ge and Le
//
//	Less = -1
//	Equal = 0
//	Greater = 1
type Ordering int8

const (
	Less    Ordering = -1
	Equal   Ordering = 0
	Greater Ordering = 1
)

// # Ord
//
// Interface for ordering of slices. Used for sorting.
type Ord[T any] interface {
	Gt(v2 T) bool
	Lt(v2 T) bool
}

// # Compare
//
// Lexicographically compare one slice with another.
//
// https://en.wikipedia.org/wiki/Lexicographic_order
//
// Returns Ordering.
func (sl Slice[T]) Compare(other Slice[T]) Ordering {
	if sl.IsEmpty() && other.IsEmpty() {
		return Equal
	}

	if sl.IsPrefixOf(other) {
		return Less
	}

	if other.IsPrefixOf(sl) {
		return Greater
	}

	for i := 0; i < sl.Len(); i++ {
		v1, v2 := sl[i], other[i]
		if v1.Gt(v2) {
			return Greater
		}
		if v1.Lt(v2) {
			return Less
		}
	}

	return Equal
}

// # Gt
//
// Returns true if slice is lexographically greater than other
func (sl Slice[T]) Gt(other any) bool {
	switch other := other.(type) {
	case Slice[T]:
		return sl.Compare(other) == Greater
	default:
		return false
	}
}

// # Lt
//
// Returns true if slice is lexographically lesser than other
func (sl Slice[T]) Lt(other any) bool {
	switch other := other.(type) {
	case Slice[T]:
		return sl.Compare(other) == Less
	default:
		return false
	}
}

// # Ge
//
// Returns true if slice is not lexographically lesser than other
func (sl Slice[T]) Ge(other any) bool {
	switch other := other.(type) {
	case Slice[T]:
		return sl.Compare(other) != Less
	default:
		return false
	}
}

// # Le
//
// Returns true if slice is not lexographically greater than other
func (sl Slice[T]) Le(other any) bool {
	switch other := other.(type) {
	case Slice[T]:
		return sl.Compare(other) != Greater
	default:
		return false
	}
}

// Gt and Lt implementations for all builtin types

func (v Int) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case int:
		return int(v) > vt
	case Int:
		return int(v) > int(vt)
	default:
		return false
	}
}

func (v I8) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case int8:
		return int8(v) > vt
	case I8:
		return int8(v) > int8(vt)
	default:
		return false
	}
}

func (v I16) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case int16:
		return int16(v) > vt
	case I16:
		return int16(v) > int16(vt)
	default:
		return false
	}
}

func (v I32) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case int32:
		return int32(v) > vt
	case I32:
		return int32(v) > int32(vt)
	default:
		return false
	}
}

func (v I64) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case int64:
		return int64(v) > vt
	case I64:
		return int64(v) > int64(vt)
	default:
		return false
	}
}

func (v Uint) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case uint:
		return uint(v) > vt
	case Uint:
		return uint(v) > uint(vt)
	default:
		return false
	}
}

func (v U8) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case uint8:
		return uint8(v) > vt
	case U8:
		return uint8(v) > uint8(vt)
	default:
		return false
	}
}

func (v U16) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case uint16:
		return uint16(v) > vt
	case U16:
		return uint16(v) > uint16(vt)
	default:
		return false
	}
}

func (v U32) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case uint32:
		return uint32(v) > vt
	case U32:
		return uint32(v) > uint32(vt)
	default:
		return false
	}
}

func (v U64) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case uint64:
		return uint64(v) > vt
	case U64:
		return uint64(v) > uint64(vt)
	default:
		return false
	}
}

func (v Str) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case string:
		return string(v) > vt
	case Str:
		return string(v) > string(vt)
	default:
		return false
	}
}

func (v Rune) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case int32:
		return rune(v) > vt
	case Rune:
		return rune(v) > rune(vt)
	default:
		return false
	}
}

func (v Bool) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case bool:
		return bool(v) && !vt
	case Bool:
		return bool(v) && !bool(vt)
	default:
		return false
	}
}

func (v C64) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case complex64:
		return real(complex128(v)) > real(complex128(vt)) && imag(complex128(v)) > imag(complex128(vt))
	case C64:
		return real(complex128(v)) > real(complex128(vt)) && imag(complex128(v)) > imag(complex128(vt))
	default:
		return false
	}
}

func (v C128) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case complex128:
		return real(v) > real(vt) && imag(v) > imag(vt)
	case C128:
		return real(v) > real(vt) && imag(v) > imag(vt)
	default:
		return false
	}
}

func (v F32) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case float32:
		return float32(v) > vt
	case F32:
		return float32(v) > float32(vt)
	default:
		return false
	}
}

func (v F64) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case float64:
		return float64(v) > vt
	case F64:
		return float64(v) > float64(vt)
	default:
		return false
	}
}

func (v Byte) Gt(v2 any) bool {
	switch vt := v2.(type) {
	case byte:
		return byte(v) > vt
	case Byte:
		return byte(v) > byte(vt)
	default:
		return false
	}
}

func (v Int) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case int:
		return int(v) < vt
	case Int:
		return int(v) < int(vt)
	default:
		return false
	}
}

func (v I8) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case int8:
		return int8(v) < vt
	case I8:
		return int8(v) < int8(vt)
	default:
		return false
	}
}

func (v I16) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case int16:
		return int16(v) < vt
	case I16:
		return int16(v) < int16(vt)
	default:
		return false
	}
}

func (v I32) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case int32:
		return int32(v) < vt
	case I32:
		return int32(v) < int32(vt)
	default:
		return false
	}
}

func (v I64) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case int64:
		return int64(v) < vt
	case I64:
		return int64(v) < int64(vt)
	default:
		return false
	}
}

func (v Uint) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case uint:
		return uint(v) < vt
	case Uint:
		return uint(v) < uint(vt)
	default:
		return false
	}
}

func (v U8) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case uint8:
		return uint8(v) < vt
	case U8:
		return uint8(v) < uint8(vt)
	default:
		return false
	}
}

func (v U16) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case uint16:
		return uint16(v) < vt
	case U16:
		return uint16(v) < uint16(vt)
	default:
		return false
	}
}

func (v U32) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case uint32:
		return uint32(v) < vt
	case U32:
		return uint32(v) < uint32(vt)
	default:
		return false
	}
}

func (v U64) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case uint64:
		return uint64(v) < vt
	case U64:
		return uint64(v) < uint64(vt)
	default:
		return false
	}
}

func (v Str) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case string:
		return string(v) < vt
	case Str:
		return string(v) < string(vt)
	default:
		return false
	}
}

func (v Rune) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case int32:
		return rune(v) < vt
	case Rune:
		return rune(v) < rune(vt)
	default:
		return false
	}
}

func (v Bool) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case bool:
		return !bool(v) && vt
	case Bool:
		return !bool(v) && bool(vt)
	default:
		return false
	}
}

func (v C64) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case complex64:
		return real(complex128(v)) < real(complex128(vt)) && imag(complex128(v)) < imag(complex128(vt))
	case C64:
		return real(complex128(v)) < real(complex128(vt)) && imag(complex128(v)) < imag(complex128(vt))
	default:
		return false
	}
}

func (v C128) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case complex128:
		return real(v) < real(vt) && imag(v) < imag(vt)
	case C128:
		return real(v) < real(vt) && imag(v) < imag(vt)
	default:
		return false
	}
}

func (v F32) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case float32:
		return float32(v) < vt
	case F32:
		return float32(v) < float32(vt)
	default:
		return false
	}
}

func (v F64) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case float64:
		return float64(v) < vt
	case F64:
		return float64(v) < float64(vt)
	default:
		return false
	}
}

func (v Byte) Lt(v2 any) bool {
	switch vt := v2.(type) {
	case byte:
		return byte(v) < vt
	case Byte:
		return byte(v) < byte(vt)
	default:
		return false
	}
}
