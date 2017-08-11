package logger

import (
	"fmt"
	"testing"
)

func TestLog_1(t *testing.T) {
	fmt.Printf("sssssssssssssss\n")
	x := Config{}
	InitConfig(x)
	for index := 0; index < 10; index++ {
		Error("aaa=%v", "eee")
		Info("aaa=%v", "eee")
		Warn("aaa=%v", "eee")
	}
}
func TestLog_2(t *testing.T) {
	InitConfig(Config{})
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
	Info("aaa=%v", "eee")
	Warn("aaa=%v", "eee")
	Error("aaa=%v", "eee")
	fmt.Printf("bbbbbb\n")
}
