package test

import (
	"fmt"
	"testing"
)

func typeConvert(in any) {
	if _, ok := in.(int); ok {
		fmt.Println("ok")
	} else {
		fmt.Println("!ok")
	}
}

func TestTypeConvert(t *testing.T) {
	typeConvert(wzyaoInt(1))
}
