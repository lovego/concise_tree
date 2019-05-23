package struct_tree

import (
	"log"
	"reflect"
	"strings"
)

var nodeType = reflect.TypeOf(Node{})
var ptr2nodeType = reflect.PtrTo(nodeType)

// 设置树的所有节点的name、code
func Setup(tree interface{}, name, code string) {
	treeValue := reflect.ValueOf(tree)
	if treeValue.Kind() != reflect.Ptr {
		log.Panicf("根节点应该是一个指针，而不是%v\n", treeValue.Kind())
	}
	setup(treeValue, name, code, true)
}

func setup(node reflect.Value, name, code string, mustBeNode bool) {
	if node.Kind() == reflect.Ptr {
		if node.IsNil() {
			node.Set(reflect.New(node.Type().Elem()))
		}
		node = node.Elem()
	}
	if node.Kind() != reflect.Struct {
		log.Panicf("节点`%s`应该是结构体，而不是%v\n", code, node.Kind())
	}
	isNode := setupChildrenFields(node, name, code)
	if mustBeNode && !isNode {
		log.Panicf("节点`%s`应该匿名嵌入tree.Node结构体\n", code)
	}
}

func setupChildrenFields(stuct reflect.Value, name, code string) (isNode bool) {
	for i, typ := 0, stuct.Type(); i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type {
		case nodeType:
			setNameCode(field, stuct.Field(i).Addr(), name, code, &isNode)
		case ptr2nodeType:
			value := stuct.Field(i)
			if value.IsNil() {
				value.Set(reflect.New(nodeType))
			}
			setNameCode(field, value, name, code, &isNode)
		default:
			setupField(field, stuct.Field(i), name, code)
		}
	}
	return
}

func setNameCode(field reflect.StructField, value reflect.Value, name, code string, isNode *bool) {
	if !exported(field.Name) {
		_, code = getNameCode(field, code)
		log.Panicf("节点`%s`不能是非导出的\n", code) // 非导出的设置不了name、code
	}
	if field.Anonymous {
		*isNode = true
	} else {
		name, code = getNameCode(field, code)
	}
	value.Interface().(*Node).SetNameCode(name, code)
}

func getNameCode(field reflect.StructField, base string) (name, code string) {
	tagParts := strings.Split(string(field.Tag), ",")
	name = tagParts[0]
	if len(tagParts) > 1 {
		code = tagParts[1]
	} else {
		code = field.Name
	}
	if base != `` {
		code = base + `.` + code
	}
	return
}

func setupField(field reflect.StructField, value reflect.Value, name, code string) {
	if exported(field.Name) {
		name, code = getNameCode(field, code)
		// 导出的字段都应该是树节点
		setup(value, name, code, true)
	} else if field.Anonymous {
		// 非导出的匿名字段不应该是树节点，只用来做类型共享，所以继续使用当前的name、code
		setup(value, name, code, false)
	} else {
		_, code := getNameCode(field, code)
		log.Panicf("`%s`不能既是非导出的，又是非匿名的\n", code)
	}
}

func exported(name string) bool {
	return name[0] >= 'A' && name[0] <= 'Z'
}
