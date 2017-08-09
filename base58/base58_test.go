package utils

import "testing"
import "fmt"

func TestBase58001(t *testing.T) {
	b58 := NewBitcoinBase58()
	ss, err := b58.EncodeToString([]byte("我们是接班人"))
	fmt.Printf("%s\n%v\n", ss, err)
}
