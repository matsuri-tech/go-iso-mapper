package main

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

func TestGenFlatStruct(t *testing.T) {
	type Huga struct {
		a string
		c uint
		d int
		t time.Time
	}
	type Hoge struct {
		b Huga
	}
	GenFlatStruct(Hoge{})
	t.Error("")
}