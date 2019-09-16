package main

import (
	"fmt"
	"matsuri/golang-tools/codeGen"
	"time"
)

type BBB string

type Huga struct {
	A BBB
	D int
	T time.Time
}
type Hoge struct {
	B Huga
}

func main() {
	result := codeGen.Generate(Hoge{})
	fmt.Println(result)
}
