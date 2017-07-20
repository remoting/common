package common

import (
	"fmt"
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
	fmt.Printf("=%s=", strings.TrimSpace(" \n  "))

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
