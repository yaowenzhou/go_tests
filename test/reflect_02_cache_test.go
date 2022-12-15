package test

import (
	"fmt"
	"reflect"
	"testing"
)

var cache = make(map[reflect.Type][]int)

// getFieldBIndex 获取结构体的属性索引列表
func getFieldBIndex(in interface{}) (index []int, err error) {
	typ := reflect.TypeOf(in)
	index, ok := cache[typ]
	if !ok {
		if typ.Kind() != reflect.Ptr {
			return index, fmt.Errorf("you must pass in a pointer")
		}
		if typ.Elem().Kind() != reflect.Struct {
			return index, fmt.Errorf("you must pass in a pointer to a struct")
		}
		f, ok := typ.Elem().FieldByName("B")
		if !ok {
			return index, fmt.Errorf("struct does not have field B")
		}
		index = f.Index
		cache[typ] = index
	}
	return
}

func populateStructReflectCache(in interface{}) error {
	index, err := getFieldBIndex(in)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(in)
	elmv := val.Elem()
	fval := elmv.FieldByIndex(index)
	fval.SetInt(42)
	return nil
}

func BenchmarkReflectWithCache(b *testing.B) { // 反射赋值(使用缓存记录成员索引)
	b.ReportAllocs()
	var m SimpleStruct
	for i := 0; i < b.N; i++ {
		if err := populateStructReflectCache(&m); err != nil {
			b.Fatal(err)
		}
		if m.B != 42 {
			b.Fatalf("unexpected value %d for B", m.B)
		}
	}
}
