package string

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestString001(t *testing.T) {
	header := make(map[string]string, 3)
	fmt.Printf("===\n")
	for k, v := range header {
		fmt.Printf("%s=%s\n", k, v)
	}
	fmt.Printf("===\n")
	fmt.Printf("=%s=\n", strings.TrimSpace(" \n  "))
	var x string
	var b map[string]string
	d := make(map[string]string, 0)
	var e map[string]string
	var c []string
	if x = "A"; len(x) > 0 {

	}
	fmt.Printf("|%s|\n", x)
	fmt.Printf("|%v|\n", b)
	fmt.Printf("|%v|\n", c)
	fmt.Printf("|%v|\n", d)
	fmt.Printf("%t\n", b == nil)
	fmt.Printf("%t\n", c == nil)
	fmt.Printf("%t\n", reflect.DeepEqual(b, d))
	fmt.Printf("%t\n", reflect.DeepEqual(b, e))
}

func TestString002(t *testing.T) {
	s1 := "a;b,c"
	fmt.Printf("%v\n", Split(s1))
	s2 := "a;中文,"
	sa2 := Split(s2)
	fmt.Printf("%d,%v\n", len(sa2), sa2)
	s3 := ""
	sa3 := Split(s3)
	fmt.Printf("%d,%v\n", len(sa3), sa3)

}
