// 反射基本版
package test

import (
	"fmt"
	"reflect"
	"testing"
)

func populateStructReflect(in interface{}) error {
	val := reflect.ValueOf(in)
	if val.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("you must pass in a pointer")
	}
	elmv := val.Elem()
	if elmv.Type().Kind() != reflect.Struct {
		return fmt.Errorf("you must pass in a pointer to a struct")
	}
	fval := elmv.FieldByName("B")
	fval.SetInt(42)
	return nil
}

func BenchmarkReflectBase(b *testing.B) { // 反射赋值基础测试
	b.ReportAllocs()
	var m SimpleStruct
	for i := 0; i < b.N; i++ {
		if err := populateStructReflect(&m); err != nil {
			b.Fatal(err)
		}
		if m.B != 42 {
			b.Fatalf("unexpected value %d for B", m.B)
		}
	}
}

// Running tool: E:\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkReflect$ go_tests/test

// goos: windows
// goarch: amd64
// pkg: go_tests/test
// cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
// BenchmarkReflect
// BenchmarkReflect-12
// 16117268                75.55 ns/op            8 B/op          1 allocs/op
// PASS
// ok      go_tests/test   1.489s

// > Test run finished at 2022/11/12 09:06:09 <
