package utils

// 差集 A - B
func Diff[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	var diff []T
	for _, v := range a {
		if _, ok := set[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}

// 交集 A ∩ B
func Intersect[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	var res []T
	for _, v := range a {
		if _, ok := set[v]; ok {
			res = append(res, v)
		}
	}
	return res
}

// 并集 A ∪ B（去重）
func Union[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(a)+len(b))
	var res []T

	for _, v := range a {
		if _, ok := set[v]; !ok {
			set[v] = struct{}{}
			res = append(res, v)
		}
	}
	for _, v := range b {
		if _, ok := set[v]; !ok {
			set[v] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}
