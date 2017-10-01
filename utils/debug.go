package utils

import (
	"fmt"
)

func EchoStrSlice(strs ...[]string) {
	for i, val := range strs {
		fmt.Print(`[`, i, `] `)
		for j, v := range val {
			fmt.Println(j, `: `, v)
		}
		fmt.Print("\n")
	}
}

func EchoBytes(args interface{}) {
	switch v := args.(type) {
	case []byte:
		fmt.Println(string(v))
	case []rune:
		fmt.Println(string(v))
	case [][]byte:
		for i, val := range v {
			fmt.Print(i, `: `, string(val), ` `)
		}
		fmt.Print("\n")
	case [][]rune:
		for i, val := range v {
			fmt.Print(i, `: `, string(val), ` `)
		}
		fmt.Print("\n")
	case [][][]byte:
		for i, val := range v {
			fmt.Print(i, ` `)
			EchoBytes(val)
		}
	case [][][]rune:
		for i, val := range v {
			fmt.Print(i, ` `)
			EchoBytes(val)
		}
	default:
		fmt.Println(v)
	}

}
