package main

import (
	"fmt"
	"matsuri/golang-tools/codeGen"
	"time"
)

type BBB string

type Huga struct {
	A BBB
	C int
	D int `goMapper:"coarseString"`
	T time.Time
}
type Hoge struct {
	B Huga
}

func main() {
	result := codeGen.Generate(Hoge{}, "toString", "2006-01-02")
	fmt.Println(result)
}
