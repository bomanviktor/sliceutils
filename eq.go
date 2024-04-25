package sliceutils

// This is the interface you need to implement on your type to be able to put it inside the Slice
type Eq[T any] interface {
	Eq(v T) bool
}

// Implementations of Eq on all the builtin types
func (i Int) Eq(v any) bool {
	return i == v
}

func (i I8) Eq(v any) bool {
	return i == v
}

func (i I16) Eq(v any) bool {
	return i == v
}

func (i I32) Eq(v any) bool {
	return i == v
}

func (i I64) Eq(v any) bool {
	return i == v
}

func (u Uint) Eq(v any) bool {
	return u == v
}

func (u U8) Eq(v any) bool {
	return u == v
}

func (u U16) Eq(v any) bool {
	return u == v
}

func (u U32) Eq(v any) bool {
	return u == v
}

func (u U64) Eq(v any) bool {
	return u == v
}

func (b Byte) Eq(v any) bool {
	return b == v
}

func (f F32) Eq(v any) bool {
	return f == v
}
func (f F64) Eq(v any) bool {
	return f == v
}
func (b Bool) Eq(v any) bool {
	return b == v
}

func (c C64) Eq(v any) bool {
	return c == v
}

func (c C128) Eq(v any) bool {
	return c == v
}

func (s Str) Eq(v any) bool {
	return s == v
}

func (r Rune) Eq(v any) bool {
	return r == v
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
