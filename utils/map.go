package utils

// CopyMap 【浅拷贝】会复制一个 map（值类型是可直接赋值的）
func CopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// DeepCopyMap 深拷贝
func DeepCopyMap[K comparable, V any](src map[K]V, clone func(V) V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = clone(v) // 调用外部提供的复制方法
	}
	return dst
}
