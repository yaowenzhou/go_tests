package test

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

// 测试cast.ToXxx类型的效果
// 可以得出其只对go自定义的数据类型如 int/uint/float32等等有效

func TestCastToXxx(t *testing.T) {
	var x wzyaoInt = 1
	fmt.Println(cast.ToInt64(x))
}

// Running tool: E:\Go\bin\go.exe test -timeout 30s -run ^TestCastToXxx$ go_tests/test
// === RUN   TestCastToXxx
// 0
// --- PASS: TestCastToXxx (0.00s)
// PASS
// ok      go_tests/test   0.507s

// > Test run finished at 2022/11/13 16:03:08 <
