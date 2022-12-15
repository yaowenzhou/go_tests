package test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

// get struct field's offset
func getFieldOffset(in interface{}, fieldIndex int) {
	if in == nil {
		fmt.Println("paramter in is nil")
		return
	}
	if reflect.ValueOf(in).Pointer() == 0 {
		fmt.Println("paramter in is a nil pointer")
		return
	}
	if reflect.ValueOf(in).Pointer() == 0 {
		fmt.Println("paramter in is a nil pointer")
		return
	}
}

func interfaceNil(in []interface{}) {
	for _, v := range in {
		getFieldOffset(v, 0)
		fmt.Println((*intface)(unsafe.Pointer(&v)).typ)
		if v == nil {
			fmt.Println("v == nil(1) is ok")
			return
		}
		if reflect.ValueOf(v).Pointer() == 0 {
			fmt.Println("v == nil(2) is ok")
			return
		}
	}
	fmt.Println("v == nil is !ok")
}

func TestInterfaceNil(t *testing.T) {
	var ints []interface{}
	var simpleStructPtr *SimpleStruct
	ints = append(ints, nil, simpleStructPtr)
	interfaceNil(ints)
}
