package a

import (
	"fmt"
)

type Hoge struct { // want "feature envy"
	hoge int
}

func (hoge *Hoge) f() {
	huga := Huga{}
	fmt.Println(hoge.hoge, huga.huga, huga.huga, huga.huga, huga.huga, huga.huga)
}

type Huga struct {
	huga int
}
