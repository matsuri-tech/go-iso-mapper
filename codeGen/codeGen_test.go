package codeGen

import (
	"testing"
	"time"
)

func TestGetPkgPrefix(t *testing.T) {
	s := "aaa/cccc/bbbb"
	result := getPackagePrefix(s)
	if result != "bbbb." {
		t.Error(result)
	}
}

/*func TestGenInitializer(t *testing.T) {
	type Hoge struct {
		a time.Time
		b string
		c time.Duration
	}
	GenInitializer(Hoge{})
}*/

/*func TestGenFlattenStruct(t *testing.T) {
	type Huga struct {
		a string
		c uint
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	result := GenFlattenStruct(Hoge{})
	if result != "type Hoge struct {  a string  c uint  d int  t time.Time}" {
		t.Error()
	}
}

func TestGenFlattenStruct2(t *testing.T) {
	type Huga struct {
		a string
		c uint
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
		e string
	}
	result := GenFlattenStruct(Hoge{})
	if result != "type Hoge struct {  a string  c uint  d int  t time.Time}" {
		t.Error(result)
	}
}

func TestGenPrimitiveStruct(t *testing.T) {
	type Huga struct {
		a string
		c uint
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	result := GenPrimitiveStruct(Hoge{})
	if result != "type Hoge struct {  b struct {}}" {
		t.Error(result)
	}
}*/

func TestGenPrimitiveStructMap(t *testing.T) {
	type Huga struct {
		a string
		c uint
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	result := GenPrimitiveStructMap(Hoge{})

	expectedMap := StructMap{
		"b": StructMap{
			"c": "uint",
			"d": "int",
			"t": "string",
			"a": "string",
		},
	}

	if !isMapEqual(result, expectedMap) {
		t.Error(result)
	}

}

func TestGenPrimitiveStructMap2(t *testing.T) {
	type BBB string

	type CCC uint

	type Huga struct {
		a BBB
		c CCC
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	result := GenPrimitiveStructMap(Hoge{})

	expectedMap := StructMap{
		"b": StructMap{
			"c": "uint",
			"d": "int",
			"t": "string",
			"a": "string",
		},
	}

	if !isMapEqual(result, expectedMap) {
		t.Error(result)
	}

}

func TestGenPrimitiveStructMap3(t *testing.T) {
	type BBB string

	type CCC uint

	type Huga struct {
		a BBB
		c CCC `coarseString:"true"`
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	result := GenPrimitiveStructMap(Hoge{})

	expectedMap := StructMap{
		"b": StructMap{
			"c": "string",
			"d": "int",
			"t": "string",
			"a": "string",
		},
	}

	if !isMapEqual(result, expectedMap) {
		t.Error(result)
	}

}

func TestGenerateMapper(t *testing.T) {
	type BBB string

	type Huga struct {
		A BBB
		D int
	}
	type Hoge struct {
		B Huga
	}
	result := generateMapper(Hoge{})

	if result != "" {
		t.Error(result)
	}

}
