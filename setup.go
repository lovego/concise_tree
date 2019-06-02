package struct_tree

import (
	"log"
	"reflect"
	"strings"
)

type NodeInfo struct {
	name, code, anonymousField   string
	childrenNames, childrenCodes map[string]struct{}
}

func (ni NodeInfo) path() string {
	if ni.anonymousField == `` {
		return ni.code
	} else {
		return ni.code + `.` + ni.anonymousField
	}
}

func (ni NodeInfo) mustBeNode() bool {
	return ni.anonymousField == ``
}

func (ni *NodeInfo) getNameCode(field reflect.StructField) (name, code string) {
	tagParts := strings.Split(string(field.Tag), ",")
	name = tagParts[0]
	if len(tagParts) > 1 && tagParts[1] != "" {
		code = tagParts[1]
	} else {
		code = lowerFirstByte(field.Name)
	}
	ni.addChild(name, code)

	if ni.code != `` {
		code = ni.code + `.` + code
	}
	if name == `` {
		log.Panicf(`节点"%s"名称不能为空`, code)
	}
	return
}

func (ni *NodeInfo) addChild(name, code string) {
	if ni.childrenNames == nil {
		ni.childrenNames = make(map[string]struct{})
	}
	if _, ok := ni.childrenNames[name]; ok {
		log.Panicf(`节点"%s"有同名的子节点"%s"`, ni.code, name)
	} else {
		ni.childrenNames[name] = struct{}{}
	}
	if ni.childrenCodes == nil {
		ni.childrenCodes = make(map[string]struct{})
	}
	if _, ok := ni.childrenCodes[code]; ok {
		log.Panicf(`节点"%s"有同Code的子节点"%s"`, ni.code, code)
	} else {
		ni.childrenCodes[code] = struct{}{}
	}
}

func (ni NodeInfo) getCode(field reflect.StructField) (code string) {
	tagParts := strings.Split(string(field.Tag), ",")
	if len(tagParts) > 1 && tagParts[1] != "" {
		code = tagParts[1]
	} else {
		code = field.Name
	}
	if ni.code != `` {
		code = ni.code + `.` + code
	}
	return
}

// 设置树的所有节点的name、code
func Setup(tree NodeIfc, name, code string) {
	treeValue := reflect.ValueOf(tree)
	if treeValue.Kind() != reflect.Ptr {
		log.Panicf(`根节点应该是一个指针，而不是%v`, treeValue.Kind())
	}
	setup(treeValue, &NodeInfo{name: name, code: code})
}

func setup(node reflect.Value, info *NodeInfo) {
	if node.Kind() == reflect.Ptr {
		if node.IsNil() {
			if node.CanSet() {
				node.Set(reflect.New(node.Type().Elem()))
			} else {
				log.Panicf(`字段"%s"不能Set（使用非指针类型，或者导出该字段）`, info.path())
			}
		}
		node = node.Elem()
	}
	if node.Kind() != reflect.Struct {
		log.Panicf(`节点"%s"应该是结构体，而不是%v`, info.path(), node.Kind())
	}
	isNode := setupChildrenFields(node, info)
	if info.mustBeNode() && !isNode {
		log.Panicf(`节点"%s"应该匿名嵌入tree.Node结构体`, info.path())
	}
}

func setupChildrenFields(stuct reflect.Value, info *NodeInfo) (isNode bool) {
	for i, typ := 0, stuct.Type(); i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type {
		case nodeType:
			setupNameCode(field, stuct.Field(i).Addr(), info, &isNode)
		case ptr2nodeType:
			value := stuct.Field(i)
			if value.IsNil() {
				value.Set(reflect.New(nodeType))
			}
			setupNameCode(field, value, info, &isNode)
		default:
			setupField(field, stuct.Field(i), info)
		}
	}
	return
}

func setupNameCode(field reflect.StructField, value reflect.Value, info *NodeInfo, isNode *bool) {
	if !exported(field.Name) {
		log.Panicf(`节点"%s"不能是非导出的`, info.getCode(field)) // 非导出的设置不了name、code
	}
	name, code := info.name, info.code
	if field.Anonymous {
		*isNode = true
	} else {
		name, code = info.getNameCode(field)
	}
	value.Interface().(*Node).SetNameCode(name, code)
}

func setupField(field reflect.StructField, value reflect.Value, info *NodeInfo) {
	if field.Anonymous && field.Tag == `` {
		// 匿名嵌入且节点名称为空，只用来做类型共享，所以继续使用当前的NodeInfo
		info.anonymousField = field.Name
		setup(value, info)
	} else if exported(field.Name) {
		// 其余的导出字段都应该是树节点
		name, code := info.getNameCode(field)
		setup(value, &NodeInfo{name: name, code: code})
	} else {
		log.Panicf(`非导出字段"%s"不能作为树节点`, info.getCode(field))
	}
}

func exported(name string) bool {
	return name[0] >= 'A' && name[0] <= 'Z'
}

func lowerFirstByte(s string) string {
	if s[0] >= 'A' && s[0] <= 'Z' {
		b := []byte(s)
		b[0] += ('a' - 'A')
		s = string(b)
	}
	return s
}
