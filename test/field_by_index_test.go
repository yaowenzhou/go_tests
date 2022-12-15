package test

import (
	"fmt"
	"reflect"
	"testing"
)

type SimpleStruct1 struct {
	A int
	B *SimpleStruct
}

func fieldByIndex(in interface{}) (index []int, err error) {
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
		if f.Type.Kind() != reflect.Ptr {
			return index, fmt.Errorf("the field must be a pointer")
		}
		g := f.Type.Elem()
		if g.Kind() != reflect.Struct {
			return index, fmt.Errorf("the field must be a pointer of struct")
		}
		f, ok = g.FieldByName("B")
		if !ok {
			return index, fmt.Errorf("struct's field B does not have field B")
		}
		index = f.Index
		cache[typ] = index
	}
	reflect.ValueOf(in).Elem().FieldByIndex([]int{1, 1}).SetInt(12)
	return
}

func TestFieldByIndex(t *testing.T) {
	ss1 := &SimpleStruct1{}
	ss1.B = &SimpleStruct{}
	index, err := fieldByIndex(ss1)
	if err == nil {
		fmt.Println(index)
	} else {
		fmt.Println(err)
	}
	fmt.Println("ss1.B.B =", ss1.B.B)
}
