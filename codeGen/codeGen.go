package codeGen

import (
	"reflect"
	"regexp"
	"strings"
)

type StructPropInfo struct {
	Name string
	Type reflect.Type
}

type StructMap = map[string]interface{}

type StructPropInfos []StructPropInfo

func Generate(st interface{}) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	structMap := genPrimitiveStructMap(st)
	structDef := "type " + s.Name() + "Map " + showStructDef(structMap)
	mapper := generateMapper(st, structMap)
	return structDef + "\n" + mapper
}

func showStructDef(stMap StructMap) string {
	var result = "struct {"
	for k, v := range stMap {
		st, isStruct := v.(StructMap)
		if isStruct {
			result = result + "\n" + k + " " + showStructDef(st)
		} else {
			result = result + "\n" + k + " " + v.(string)
		}
	}
	return result + "\n}"
}

func strHeadLower(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

// /を含まない文字列の末尾を取得する
var r = regexp.MustCompile(`([^/]+)$`)

func getPackagePrefix(pkgPath string) string {
	s := r.FindStringSubmatch(pkgPath)
	if len(s) != 0 {
		return s[0] + "."
	}
	return ""
}

func GenInitializer(st interface{}) string {
	infos := StructPropInfos{}
	str := ""
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		info := StructPropInfo{
			Name: f.Name,
			Type: f.Type,
		}
		infos = append(infos, info)
	}
	str = str + "func New" + s.Name() + "("
	for _, r := range infos {
		str = str + "  " + strHeadLower(r.Name) + " " + getPackagePrefix(r.Type.PkgPath()) + r.Type.Name() + ","
	}
	str = str + ") *" + s.Name() + "{"
	str = str + "return " + "&" + s.Name() + "{"
	for _, r := range infos {
		str = str + "  " + r.Name + ":" + strHeadLower(r.Name) + ","
	}
	str = str + "}"
	str = str + "\n}"
	return str
}

func genPrimitiveStructMap(st interface{}) StructMap {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	numField := s.NumField()
	var result = StructMap{}
	for i := 0; i < numField; i++ {
		f := s.Field(i)
		// 再帰的にstructを探索、time.Timeはstringに潰す
		if f.Tag.Get("coarseString") == "true" {
			result[f.Name] = "string"
			continue
		}
		if f.Type.Kind().String() == "struct" {
			switch {
			case f.Type.Name() == "Time":
				result[f.Name] = "string"
			default:
				v := reflect.New(f.Type).Elem().Interface()
				result[f.Name] = genPrimitiveStructMap(v)
			}
		} else {
			result[f.Name] = f.Type.Kind().String()
		}
	}
	return result
}

func generateMapper(st interface{}, structMap StructMap) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	arg := strings.ToLower(s.Name())
	var result = "func MapFrom" + s.Name() + "(" + arg + " " + s.Name() + ") " + s.Name() + "Map" + " {\nreturn " + s.Name() + "Map"
	result = result + generateMapperSub(arg, st, structMap)
	result = result + "\n}"
	return result
}

func generateMapperSub(prefix string, st interface{}, stMap StructMap) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	var result = "{"
	numField := s.NumField()
	for i := 0; i < numField; i++ {
		f := s.Field(i)
		if f.Type.Kind().String() == "struct" {
			if f.Type.Name() == "Time" {
				result = result + "\n" + f.Name + ":" + " " + prefix + "." + f.Name + ".Format(\"2006-01-02\")" + ","
			} else {
				v := reflect.New(f.Type).Elem().Interface()
				vv, ok := stMap[f.Name]
				if !ok {
					panic(v)
				}
				vvv, ok := vv.(StructMap)
				if !ok {
					panic(vv)
				}
				result = result + "\n" + f.Name + ":" + " " + showStructDef(vvv) + generateMapperSub(prefix+"."+f.Name, v, vvv) + ","
			}
		} else {
			typ, ok := stMap[f.Name]
			if !ok {
				panic(f.Name)
			}
			result = result + "\n" + f.Name + ":" + " " + typ.(string) + "(" + prefix + "." + f.Name + ")" + ","
		}
	}
	result = result + "\n}"
	return result
}

func isMapEqual(m StructMap, another StructMap) bool {
	for k, v := range m {
		switch v.(type) {
		case StructMap:
			mm, _ := v.(StructMap)
			mmm, ok := another[k]
			if !ok {
				return false
			}
			mmmm, _ := mmm.(StructMap)
			return isMapEqual(mm, mmmm)
		default:
			vv, ok := another[k]
			if !ok {
				return false
			}
			if vv != v {
				return false
			}
		}
	}
	return true
}
