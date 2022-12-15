package test

import (
	"fmt"
	"testing"

	"golang.org/x/exp/constraints"
)

type wzyaoInt int

func compare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left < right
}

func BenchmarkMap2Slice(b *testing.B) {
	// lo.MapToSlice()
	fmt.Println(compare(wzyaoInt(1), wzyaoInt(2)))
}
