package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestTrimSpace(t *testing.T) {
	str := "    sss   "
	fmt.Printf("'%s'\n", strings.TrimSpace(str))
	var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}
	fmt.Println(asciiSpace['c'])
}
