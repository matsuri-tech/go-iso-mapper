package main

import (
	"fmt"
	"matsuri/golang-tools/codeGen"
)

type BBB string

type Huga struct {
	A BBB
	D int
}
type Hoge struct {
	B Huga
}

func main() {
	result := codeGen.Generate(Hoge{})
	fmt.Println(result)
}
