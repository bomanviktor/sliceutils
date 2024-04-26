package sliceutils

// This is the interface you need to implement on your type to be able to put it inside the Slice
type Eq[T any] interface {
	Eq(v T) bool
}

// Implementations of Eq on all the builtin types
func (i Int) Eq(v any) bool {
	switch vt := v.(type) {
	case Int:
		return i == vt
	case int:
		return int(i) == vt
	default:
		return false
	}
}
func (i I8) Eq(v any) bool {
	switch vt := v.(type) {
	case I8:
		return i == vt
	case int8:
		return int8(i) == vt
	default:
		return false
	}
}

func (i I16) Eq(v any) bool {
	switch vt := v.(type) {
	case I16:
		return i == vt
	case int16:
		return int16(i) == vt
	default:
		return false
	}
}

func (i I32) Eq(v any) bool {
	switch vt := v.(type) {
	case I32:
		return i == vt
	case int32:
		return int32(i) == vt
	default:
		return false
	}
}

func (i I64) Eq(v any) bool {
	switch vt := v.(type) {
	case I64:
		return i == vt
	case int64:
		return int64(i) == vt
	default:
		return false
	}
}

func (u Uint) Eq(v any) bool {
	switch vt := v.(type) {
	case Uint:
		return u == vt
	case uint:
		return uint(u) == vt
	default:
		return false
	}
}

func (u U8) Eq(v any) bool {
	switch vt := v.(type) {
	case U8:
		return u == vt
	case uint8:
		return uint8(u) == vt
	default:
		return false
	}
}

func (u U16) Eq(v any) bool {
	switch vt := v.(type) {
	case U16:
		return u == vt
	case uint16:
		return uint16(u) == vt
	default:
		return false
	}
}

func (u U32) Eq(v any) bool {
	switch vt := v.(type) {
	case U32:
		return u == vt
	case uint32:
		return uint32(u) == vt
	default:
		return false
	}
}

func (u U64) Eq(v any) bool {
	switch vt := v.(type) {
	case U64:
		return u == vt
	case uint64:
		return uint64(u) == vt
	default:
		return false
	}
}

func (b Byte) Eq(v any) bool {
	switch vt := v.(type) {
	case Byte:
		return b == vt
	case byte:
		return byte(b) == vt
	default:
		return false
	}
}

func (f F32) Eq(v any) bool {
	switch vt := v.(type) {
	case F32:
		return f == vt
	case float32:
		return float32(f) == vt
	default:
		return false
	}
}

func (f F64) Eq(v any) bool {
	switch vt := v.(type) {
	case F64:
		return f == vt
	case float64:
		return float64(f) == vt
	default:
		return false
	}
}

func (b Bool) Eq(v any) bool {
	switch vt := v.(type) {
	case Bool:
		return b == vt
	case bool:
		return bool(b) == vt
	default:
		return false
	}
}

func (c C64) Eq(v any) bool {
	switch vt := v.(type) {
	case C64:
		return c == vt
	case complex64:
		return complex64(c) == vt
	default:
		return false
	}
}

func (c C128) Eq(v any) bool {
	switch vt := v.(type) {
	case C128:
		return c == vt
	case complex128:
		return complex128(c) == vt
	default:
		return false
	}
}

func (s Str) Eq(v any) bool {
	switch vt := v.(type) {
	case Str:
		return s == vt
	case string:
		return string(s) == vt
	default:
		return false
	}
}

func (r Rune) Eq(v any) bool {
	switch vt := v.(type) {
	case Rune:
		return r == vt
	case rune:
		return rune(r) == vt
	default:
		return false
	}
}

// Implementation of Eq on the slice type

func (sl Slice[T]) Eq(v any) bool {
	switch v := v.(type) {
	case Slice[T]:
		if len(sl) != len(v) {
			return false
		}
		for i := 0; i < len(sl); i++ {
			if !sl[i].Eq(v[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (sl Slice[T]) Ne(v any) bool {
	return !sl.Eq(v)
}
