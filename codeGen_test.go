package goIsoMapper

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
	result := genPrimitiveStructMap(Hoge{})

	expectedMap := structMap{
		"b": structMap{
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
	result := genPrimitiveStructMap(Hoge{})

	expectedMap := structMap{
		"b": structMap{
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
		c CCC `goIsoMapper:"coarseString"`
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	result := genPrimitiveStructMap(Hoge{})

	expectedMap := structMap{
		"b": structMap{
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

func TestGenerate(t *testing.T) {
	type BBB string

	type Huga struct {
		A BBB
		C int
		D int `goIsoMapper:"coarseString"`
		T time.Time
	}
	type Hoge struct {
		B Huga
	}
	result := Generate(Hoge{}, "toString", "2006-01-02")

	// テストコードがやばすぎる
	expectedStr := `type HogeMap struct {
B struct {
A string
C int
D string
T string
}
}
func MapFromHoge(hoge Hoge) HogeMap {
return HogeMap{
B: struct {
A string
C int
D string
T string
}{
A: string(hoge.B.A),
C: hoge.B.C,
D: hoge.B.D.toString(),
T: hoge.B.T.Format("2006-01-02"),
},
}
}`

	if result != expectedStr {
		t.Error(result)
	}

}

func TestGenerate2(t *testing.T) {

	type Hoge struct {
		B []int
	}
	result := Generate(Hoge{}, "toString", "2006-01-02")

	// テストコードがやばすぎる
	expectedStr := `type HogeMap struct {
B []int
}
func MapFromHoge(hoge Hoge) HogeMap {
return HogeMap{
B: []int(hoge.B),
}
}`

	if result != expectedStr {
		t.Error(result)
	}

}
