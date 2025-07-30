package val

import "time"

func Ptr[T comparable](t T) *T {
	return &t
}

func UnPtr[T comparable](t *T) T {
	if t == nil {
		var tc T
		return tc
	}
	return *t
}

func Time(t time.Time) *time.Time {
	return &t
}

func String(a string) *string {
	return &a
}

func StringValue(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

func Int(a int) *int {
	return &a
}

func IntValue(a *int) int {
	if a == nil {
		return 0
	}
	return *a
}

func Int8(a int8) *int8 {
	return &a
}

func Int8Value(a *int8) int8 {
	if a == nil {
		return 0
	}
	return *a
}

func Int16(a int16) *int16 {
	return &a
}

func Int16Value(a *int16) int16 {
	if a == nil {
		return 0
	}
	return *a
}

func Int32(a int32) *int32 {
	return &a
}

func Int32Value(a *int32) int32 {
	if a == nil {
		return 0
	}
	return *a
}

func Int64(a int64) *int64 {
	return &a
}

func Int64Value(a *int64) int64 {
	if a == nil {
		return 0
	}
	return *a
}

func Bool(a bool) *bool {
	return &a
}

func BoolValue(a *bool) bool {
	if a == nil {
		return false
	}
	return *a
}

func Uint(a uint) *uint {
	return &a
}

func UintValue(a *uint) uint {
	if a == nil {
		return 0
	}
	return *a
}

func Uint8(a uint8) *uint8 {
	return &a
}

func Uint8Value(a *uint8) uint8 {
	if a == nil {
		return 0
	}
	return *a
}

func Uint16(a uint16) *uint16 {
	return &a
}

func Uint16Value(a *uint16) uint16 {
	if a == nil {
		return 0
	}
	return *a
}

func Uint32(a uint32) *uint32 {
	return &a
}

func Uint32Value(a *uint32) uint32 {
	if a == nil {
		return 0
	}
	return *a
}

func Uint64(a uint64) *uint64 {
	return &a
}

func Uint64Value(a *uint64) uint64 {
	if a == nil {
		return 0
	}
	return *a
}

func Float32(a float32) *float32 {
	return &a
}

func Float32Value(a *float32) float32 {
	if a == nil {
		return 0
	}
	return *a
}

func Float64(a float64) *float64 {
	return &a
}

func Float64Value(a *float64) float64 {
	if a == nil {
		return 0
	}
	return *a
}

func ArrayAny(val ...any) []any {
	return val
}
