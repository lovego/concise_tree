package concise_tree

import (
	"log"
	"reflect"

	"github.com/lovego/struct_tag"
)

// 设置树的所有节点的name、code、desc
func Setup(tree ConciseTree, name, code, desc string) {
	treeValue := reflect.ValueOf(tree)
	if treeValue.Kind() != reflect.Ptr {
		log.Panicf(`tree should be a pointer, not %v`, treeValue.Kind())
	}
	setup(treeValue, &NodeInfo{name: name, code: code, desc: desc})
}

func setup(node reflect.Value, info *NodeInfo) {
	if node.Kind() == reflect.Ptr {
		if node.IsNil() {
			if node.CanSet() {
				node.Set(reflect.New(node.Type().Elem()))
			} else {
				log.Panicf(`field "%s" CanSet: false`, info.path())
			}
		}
		node = node.Elem()
	}
	if node.Kind() != reflect.Struct {
		log.Panicf(`node "%s" should be struct, not %v`, info.path(), node.Kind())
	}
	isNode := setupChildrenFields(node, info)
	if info.mustBeNode() && !isNode {
		log.Panicf(`node "%s" should anonymously embed concise_tree.Node`, info.path())
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
		log.Panicf(`node "%s" should be exported`, info.getCode(field)) // 非导出的设置不了name、code
	}
	name, code, desc := info.name, info.code, info.desc
	if field.Anonymous {
		*isNode = true
	} else {
		name, code, desc = info.getNameCodeDesc(field)
	}
	value.Interface().(*Node).Set(name, code, desc)
}

func setupField(field reflect.StructField, value reflect.Value, info *NodeInfo) {
	if field.Anonymous && struct_tag.Get(string(field.Tag), "name") == "" {
		// 匿名嵌入且节点名称为空，只用来做类型共享，所以继续使用当前的NodeInfo
		info.anonymousField = field.Name
		setup(value, info)
	} else if exported(field.Name) {
		// 其余的导出字段都应该是树节点
		name, code, desc := info.getNameCodeDesc(field)
		setup(value, &NodeInfo{name: name, code: code, desc: desc})
	} else {
		log.Panicf(`tree node "%s" must be exported`, info.getCode(field))
	}
}

type NodeInfo struct {
	name, code, desc             string
	anonymousField               string
	childrenNames, childrenCodes map[string]struct{}
}

func (ni NodeInfo) mustBeNode() bool {
	return ni.anonymousField == `` // non anonymous field must be a tree node
}

func (ni NodeInfo) path() string {
	if ni.anonymousField == `` {
		return ni.code
	} else {
		return ni.code + `.` + ni.anonymousField
	}
}

func (ni NodeInfo) getCode(field reflect.StructField) (code string) {
	if code = struct_tag.Get(string(field.Tag), "code"); code == "" {
		code = lowerFirstByte(field.Name)
	}
	if ni.code != `` {
		code = ni.code + `.` + code
	}
	return
}

func (ni *NodeInfo) getNameCodeDesc(field reflect.StructField) (name, code, desc string) {
	if name = struct_tag.Get(string(field.Tag), "name"); name == `` {
		log.Panicf(`name of node "%s" is empty`, code)
	}
	if code = struct_tag.Get(string(field.Tag), "code"); code == "" {
		code = lowerFirstByte(field.Name)
	}
	ni.validateChild(name, code)
	if ni.code != `` {
		code = ni.code + `.` + code
	}

	desc = struct_tag.Get(string(field.Tag), "desc")
	return
}

func (ni *NodeInfo) validateChild(name, code string) {
	if ni.childrenNames == nil {
		ni.childrenNames = make(map[string]struct{})
	}
	if _, ok := ni.childrenNames[name]; ok {
		log.Panicf(`node "%s" has children node of same name "%s"`, ni.code, name)
	} else {
		ni.childrenNames[name] = struct{}{}
	}
	if ni.childrenCodes == nil {
		ni.childrenCodes = make(map[string]struct{})
	}
	if _, ok := ni.childrenCodes[code]; ok {
		log.Panicf(`node "%s" has children node of same code "%s"`, ni.code, code)
	} else {
		ni.childrenCodes[code] = struct{}{}
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
