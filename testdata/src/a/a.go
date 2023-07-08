package a

import (
	"fmt"
)

func f() { // want "pattern"
	var gopher int
	fmt.Println(gopher)
}
