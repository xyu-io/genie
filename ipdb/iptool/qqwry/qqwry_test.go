package qqwry

import (
	"fmt"
	"testing"
)

func TestQqwry(t *testing.T) {
	qQwry, err := NewQQwry("./qqwry_utf.dat")
	if err != nil {
		return
	}
	ipStr := "117.187.194.117"

	find, err := qQwry.Find(ipStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(find)
}
