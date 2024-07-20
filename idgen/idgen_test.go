package idgen

import (
	"fmt"
	"testing"
)

func TestNewBase64IdGen(t *testing.T) {
	idgen, _ := NewIdGen(0)
	idgen1, _ := NewIdGen(1)
	for i := 0; i < 10; i++ {
		fmt.Println(idgen.NextId())
		fmt.Println(idgen.NextString())
		fmt.Println(idgen1.NextId())
		fmt.Println(idgen1.NextString())
		fmt.Println()
	}
}
