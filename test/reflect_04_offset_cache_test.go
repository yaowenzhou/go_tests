package test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

var unsafeCache2 = make(map[uintptr]uintptr)

func populateStructUnsafe2(in interface{}) error {
	inf := (*intface)(unsafe.Pointer(&in))
	offset, ok := unsafeCache2[uintptr(inf.typ)]
	if !ok {
		typ := reflect.TypeOf(in)
		if typ.Kind() != reflect.Ptr {
			return fmt.Errorf("you must pass in a pointer")
		}
		if typ.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("you must pass in a pointer to a struct")
		}
		f, ok := typ.Elem().FieldByName("B")
		if !ok {
			return fmt.Errorf("struct does not have field B")
		}
		if f.Type.Kind() != reflect.Int {
			return fmt.Errorf("field B should be an int")
		}
		offset = f.Offset
		unsafeCache2[uintptr(inf.typ)] = offset
	}
	*(*int)(unsafe.Pointer(uintptr(inf.value) + offset)) = 42
	return nil
}

// map[reflect.Type]uintptr的方案，map需要调用hash接口检查二者是否相等
// 可以直接记录接口类型信息的地址即可免去hash计算

func BenchmarkReflectWithUnsafe1(b *testing.B) {
	b.ReportAllocs()
	var m SimpleStruct
	for i := 0; i < b.N; i++ {
		if err := populateStructUnsafe2(&m); err != nil {
			b.Fatal(err)
		}
		if m.B != 42 {
			b.Fatalf("unexpected value %d for B", m.B)
		}
	}
}
