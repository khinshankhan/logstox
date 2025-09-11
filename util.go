package logstox

func First[T any](vals []T, keep func(T) bool) (*T, bool) {
	for _, v := range vals {
		if keep(v) {
			return &v, true
		}
	}
	return nil, false
}

func getZero[T any]() T {
	var zero T
	return zero
}

func FirstNonZero[T comparable](zero T, vals ...T) (T, bool) {
	z := getZero[T]()

	ptr, ok := First(vals, func(v T) bool { return v != z })
	if ok {
		return *ptr, true
	}
	return z, false
}
