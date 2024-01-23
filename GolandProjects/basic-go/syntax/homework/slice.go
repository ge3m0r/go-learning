package homework

// 切片的删除操作
func Delete(idx int, vals []any) []any {
	//对idx进行验证
	if idx < 0 || idx >= len(vals) {
		panic("idx超过切片范围")
	}

	//循环替代
	for i := idx; i < len(vals)-1; i++ {
		vals[i] = vals[i+1]
	}
	vals = vals[:len(vals)-1]
	return vals
}

// 改造为泛型方法
func Delete1[T any](idx int, vals []T) []T {
	if idx < 0 || idx >= len(vals) {
		panic("idx超过切片范围")
	}

	//循环替代
	for i := idx; i < len(vals)-1; i++ {
		vals[i] = vals[i+1]
	}
	vals = vals[:len(vals)-1]
	return vals
}

// 高性能和缩容机制
// 容积大的时候可以快速缩容，节省空间，可以每次 0.8 倍缩容，而容积小每次到了 0.5 原来的内容才进行第二次缩容
// 缩容函数
func Shrink(val []int) []int {
	c, l := cap(val), len(val)
	newc := c
	if c <= 256 && (c/l >= 2) {
		newc = int(float32(c) * 0.8)
	} else if c >= 256 && c/l >= 4 {
		newc = int(float32(c) / 2)
	}
	if newc != c {
		s := make([]int, 0, newc)
		s = append(s, val...)
		return s
	} else {
		return nil
	}

}
