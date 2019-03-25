package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type StructPropInfo struct {
	Name string
	Type reflect.Type
}

type StructPropInfos []StructPropInfo

func (si StructPropInfo) getGrpcTypeStr() string {
	switch si.Type {
	default:
		return "string"
	}
}

func strHeadLower(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

// /を含まない文字列の末尾を取得する
var  r = regexp.MustCompile(`([^/]+)$`)

func getPackagePrefix(pkgPath string) string {
	s := r.FindStringSubmatch(pkgPath)
	if len(s) != 0 {
		return s[0] + "."
	}
	return ""
}

func GenInitializer(st interface{}) string {
	infos := StructPropInfos{}
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		info := StructPropInfo{
			Name: f.Name,
			Type: f.Type,
		}
		infos = append(infos, info)
	}
	fmt.Println("func New" + s.Name() + "(")
	for _, r := range infos {
		fmt.Println("  " + strHeadLower(r.Name) + " " + getPackagePrefix(r.Type.PkgPath()) + r.Type.Name() + ",")
	}
	fmt.Println(") *" + s.Name() + "{")
	fmt.Println("return " + "&" + s.Name() + "{")
	for _, r := range infos {
		fmt.Println("  " + r.Name + ":" + strHeadLower(r.Name) + ",")
	}
	fmt.Println(" }")
	fmt.Println("}")
	return ""
}

func GenFlatStruct(st interface{}) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	name := s.Name()
	props := genFlatStructSub(st)
	fmt.Println("type " + name + " struct {")
	for _,prop := range props {
		fmt.Println("  " + prop.Name + " " + getPackagePrefix(prop.Type.PkgPath()) + prop.Type.Name())
	}
	fmt.Println("}")
	return ""
}

func genFlatStructSub(st interface{}) StructPropInfos {
	infos := StructPropInfos{}
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	numField := s.NumField()
	for i := 0; i < numField; i++ {
		f := s.Field(i)
		if f.Type.Kind().String() == "struct" && f.Type.Name() != "Time" {
			v := reflect.New(f.Type).Elem().Interface()
			fmt.Println(v)
			infos = append(infos, genFlatStructSub(v)...)
		} else {
			info := StructPropInfo{
				Name: f.Name,
				Type: f.Type,
			}
			infos = append(infos, info)
		}
	}
	return infos
}

func ProtoGen(st interface{}) string {
	infos := StructPropInfos{}
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Println(f.Name + ":" + f.Type.Name() + ":" + f.Type.Kind().String())
		info := StructPropInfo{
			Name: f.Name,
			Type: f.Type,
		}
		infos = append(infos, info)
	}
	fmt.Println("message " + s.Name() + "{")
	for _, r := range infos {
		fmt.Println("	" + r.getGrpcTypeStr() + " " + r.Name)
	}
	fmt.Println("}")
	return ""
}
