package nolock

import (
	"testing"
	"fmt"
)

type Aa struct {
}

func TestHt(t *testing.T) {
	ht := NewHashTable(5)
	ht.Put(`abc`, `opajfowe`, 0)
	ht.Put(`bcde`, `aaaaaaaaaaaaa`, 0)
	a := ht.Get(`abc`)
	b := ht.Get(`bcde`)
	fmt.Println(a, b)
	fmt.Println(ht)
}
