package test

import (
	"fmt"
	"reflect"
	"testing"
)

// TestReflectCustomType 自定义数据类型使用reflect查看其Kind
func TestReflectCustomType(t *testing.T) {
	var i wzyaoInt = 1
	fmt.Println(reflect.TypeOf(&i).Elem().Kind())
	fmt.Println(reflect.ValueOf(&i).Elem().Type())
	fmt.Println(reflect.ValueOf(&SimpleStruct{}).Elem().Field(1))
	fmt.Println(reflect.ValueOf(&SimpleStruct{}).Elem().Field(1).Type().Kind())
}
