package log

import (
	"fmt"
	"testing"
)

func TestLog_1(t *testing.T) {
	fmt.Printf("sssssssssssssss\n")
	x := LogConfig{}
	InitConfig(x)
	for index := 0; index < 10; index++ {
		Error("aaa=%v", "eee")
		Debug("aaa=%v", "eee")
		Warn("aaa=%v", "eee")
	}
}
func TestLog_2(t *testing.T) {
	InitConfig(LogConfig{})
	fmt.Printf("sssssssssssssss\n")
	xxx()
}

func xxx() {
	yyy()
}
func yyy() {
	zzz()
}
func zzz() {
	//fmt.Printf("%s", Stack())
	Error("aaa=%v", "eee")
	fmt.Printf("bbbbbb\n")
}
