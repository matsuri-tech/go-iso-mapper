package codeGen

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type structPropInfo struct {
	Name string
	Type reflect.Type
}

type structMap = map[string]interface{}

type structPropInfos []structPropInfo

func Generate(st interface{}, stringMapperMethodName string, timeFormat string) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	structMap := genPrimitiveStructMap(st)
	structDef := "type " + s.Name() + "Map " + showStructDef(structMap)
	mapper := generateMapper(st, structMap, stringMapperMethodName, timeFormat)
	return structDef + "\n" + mapper
}

func showStructDef(stMap structMap) string {
	var result = "struct {"

	// 型定義は順序依存なので,keyをsortする必要がある
	var keys []string
	for k := range stMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		st, isStruct := stMap[k].(structMap)
		if isStruct {
			result = result + "\n" + k + " " + showStructDef(st)
		} else {
			result = result + "\n" + k + " " + stMap[k].(string)
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
	infos := structPropInfos{}
	str := ""
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		info := structPropInfo{
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

func genPrimitiveStructMap(st interface{}) structMap {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	numField := s.NumField()
	var result = structMap{}
	for i := 0; i < numField; i++ {
		f := s.Field(i)
		// 再帰的にstructを探索、time.Timeはstringに潰す
		if f.Tag.Get("goIsoMapper") == "coarseString" {
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
		} else if f.Type.Kind().String() == "slice" {
			result[f.Name] = f.Type.String()
		} else {
			result[f.Name] = f.Type.Kind().String()
		}
	}
	return result
}

func generateMapper(st interface{}, structMap structMap, stringMapperMethodName string, timeFormat string) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	arg := strings.ToLower(s.Name())
	var result = "func MapFrom" + s.Name() + "(" + arg + " " + s.Name() + ") " + s.Name() + "Map" + " {\nreturn " + s.Name() + "Map"
	result = result + generateMapperSub(arg, stringMapperMethodName, timeFormat, st, structMap)
	result = result + "\n}"
	return result
}

func generateMapperSub(prefix string, stringMapperMethodName string, timeFormat string, st interface{}, stMap structMap) string {
	s := reflect.New(reflect.TypeOf(st)).Elem().Type()
	var result = "{"
	numField := s.NumField()
	for i := 0; i < numField; i++ {
		f := s.Field(i)
		if f.Tag.Get("goIsoMapper") == "coarseString" {
			result = result + "\n" + f.Name + ":" + " " + prefix + "." + f.Name + "." + stringMapperMethodName + "()" + ","
			continue
		}

		if f.Type.Kind().String() == "struct" {
			if f.Type.Name() == "Time" {
				result = result + "\n" + f.Name + ":" + " " + prefix + "." + f.Name + ".Format(\"" + timeFormat + "\")" + ","
			} else {
				v := reflect.New(f.Type).Elem().Interface()
				vv, ok := stMap[f.Name]
				if !ok {
					panic(v)
				}
				vvv, ok := vv.(structMap)
				if !ok {
					panic(vv)
				}
				result = result + "\n" + f.Name + ":" + " " + showStructDef(vvv) + generateMapperSub(prefix+"."+f.Name, stringMapperMethodName, timeFormat, v, vvv) + ","
			}
		} else {
			typ, ok := stMap[f.Name]
			if !ok {
				panic(f.Name)
			}
			if typ == f.Type.Name() {
				result = result + "\n" + f.Name + ":" + " " + prefix + "." + f.Name + ","
			} else {
				result = result + "\n" + f.Name + ":" + " " + typ.(string) + "(" + prefix + "." + f.Name + ")" + ","
			}
		}
	}
	result = result + "\n}"
	return result
}

func isMapEqual(m structMap, another structMap) bool {
	for k, v := range m {
		switch v.(type) {
		case structMap:
			mm, _ := v.(structMap)
			mmm, ok := another[k]
			if !ok {
				return false
			}
			mmmm, _ := mmm.(structMap)
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
