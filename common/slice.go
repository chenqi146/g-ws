package common

type BaseType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// DeleteSlice 删除指定元素。
func DeleteSlice[T BaseType](s []T, elem T) []T {
	j := 0
	for _, v := range s {

		if elem != v {
			s[j] = v
			j++
		}
	}
	return s[:j]
}
