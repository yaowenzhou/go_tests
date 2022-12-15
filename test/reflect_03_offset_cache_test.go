package test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

var unsafeCache = make(map[reflect.Type]uintptr)

type intface struct {
	typ   unsafe.Pointer
	value unsafe.Pointer
}

func populateStructUnsafe(in interface{}) error {
	typ := reflect.TypeOf(in)
	offset, ok := unsafeCache[typ]
	if !ok {
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
		offset = f.Offset         // 获取偏移量
		unsafeCache[typ] = offset // 保存偏移量
	}
	// intface是一个与空接口相同的定义，使用它接收in的指针，然后方便获取结构体指针
	structPtr := (*intface)(unsafe.Pointer(&in)).value       // 获取结构体的指针
	fieldBPtr := unsafe.Pointer(uintptr(structPtr) + offset) // 得到结构体的成员B的地址
	*(*int)(fieldBPtr) = 42                                  // 表示将fieldBPtr转化为int指针类型然后解指针赋值为42
	return nil
}

func BenchmarkReflectWithUnSafe(b *testing.B) { // 缓存偏移量，然后直接使用偏移量来直接进行赋值
	b.ReportAllocs()
	var m SimpleStruct
	for i := 0; i < b.N; i++ {
		if err := populateStructUnsafe(&m); err != nil {
			b.Fatal(err)
		}
		if m.B != 42 {
			b.Fatalf("unexpected value %d for B", m.B)
		}
	}
}
