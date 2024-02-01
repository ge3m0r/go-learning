package main

import "fmt"

func Array() {
	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1:%v, len:%d, cap:%d", a1, len(a1), cap(a1))
}

func Slice() {
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("s1:%v, len:%d, cap:%d", s1, len(s1), cap(s1))
	s2 := make([]int, 3, 4)
	fmt.Printf("s2:%v, len:%d, cap:%d", s2, len(s2), cap(s2))
	s2 = append(s2, 9)
	fmt.Printf("s2:%v, len:%d, cap:%d", s2, len(s2), cap(s2))
}
