package concise_tree

import (
	"log"
	"reflect"

	"github.com/lovego/struct_tag"
)

type nodeInfo struct {
	value          reflect.Value
	field          reflect.StructField
	path           string
	tags           map[string]string
	anonymousField string
	isNode         bool
	childrenCodes  map[string]struct{}
	childrenNames  map[string]struct{}
}

func (ni *nodeInfo) setup() {
	if ni.value.Kind() == reflect.Ptr {
		if ni.value.IsNil() {
			if ni.value.CanSet() {
				ni.value.Set(reflect.New(ni.value.Type().Elem()))
			} else {
				log.Panicf(`field "%s" CanSet: false`, ni.getPath())
			}
		}
		ni.value = ni.value.Elem()
	}
	if ni.value.Kind() != reflect.Struct {
		log.Panicf(`node "%s" should be struct, not %v`, ni.getPath(), ni.value.Kind())
	}
	ni.setupChildren()
	if ni.mustBeNode() && !ni.isNode {
		log.Panicf(`node "%s" should anonymously embed concise_tree.Node`, ni.getPath())
	}
}

func (ni *nodeInfo) setupChildren() {
	for i := 0; i < ni.value.NumField(); i++ {
		child := &nodeInfo{
			value: ni.value.Field(i),
			field: ni.value.Type().Field(i),
		}
		child.tags = struct_tag.Parse(string(child.field.Tag))

		switch child.value.Type() {
		case nodeType:
			child.value = child.value.Addr()
			child.setupAsLeaf(ni)
		case ptr2nodeType:
			if child.value.IsNil() {
				child.value.Set(reflect.New(nodeType))
			}
			child.setupAsLeaf(ni)
		default:
			child.setupAsNonLeaf(ni)
		}
	}
}

func (ni *nodeInfo) setPath(parent *nodeInfo) {
	code := ni.tags["code"]
	if code == "" {
		code = lowerFirstByte(ni.field.Name)
	}

	if parent.path != `` {
		ni.path = parent.path + `.` + code
	} else {
		ni.path = code
	}

	if ni.tags["name"] == `` {
		log.Panicf(`name of node "%s" is empty`, ni.path)
	}
	ni.validateChild(code, ni.tags["name"])
}

func (ni *nodeInfo) setupAsLeaf(parent *nodeInfo) {
	if !exported(ni.field.Name) {
		log.Panicf(`node "%s" should be exported`, ni.path) // 非导出的设置不了path、tags
	}
	if ni.field.Anonymous {
		parent.isNode = true
		ni.value.Interface().(*Node).Set(parent.path, parent.tags)
	} else {
		ni.setPath(parent)
		ni.value.Interface().(*Node).Set(ni.path, ni.tags)
	}
}

func (ni *nodeInfo) setupAsNonLeaf(parent *nodeInfo) {
	if ni.field.Anonymous && ni.tags["name"] == "" {
		// 匿名嵌入且节点名称为空，只用来做类型共享，所以继续使用parent的path
		parent.anonymousField = ni.field.Name
		ni.path = parent.path
		ni.setup()
	} else if exported(ni.field.Name) {
		// 其余的导出字段都应该是树节点
		ni.setPath(parent)
		ni.setup()
	} else {
		log.Panicf(`tree node "%s" must be exported`, ni.path)
	}
}

func (ni nodeInfo) mustBeNode() bool {
	return ni.anonymousField == `` // non anonymous field must be a tree node
}

func (ni nodeInfo) getPath() string {
	if ni.anonymousField == `` {
		return ni.path
	} else {
		return ni.path + `.` + ni.anonymousField
	}
}

func (ni *nodeInfo) validateChild(code, name string) {
	if ni.childrenCodes == nil {
		ni.childrenCodes = make(map[string]struct{})
	}
	if _, ok := ni.childrenCodes[code]; ok {
		log.Panicf(`node "%s" has children node of same code "%s"`, ni.path, code)
	} else {
		ni.childrenCodes[code] = struct{}{}
	}

	if ni.childrenNames == nil {
		ni.childrenNames = make(map[string]struct{})
	}
	if _, ok := ni.childrenNames[name]; ok {
		log.Panicf(`node "%s" has children node of same name "%s"`, ni.path, name)
	} else {
		ni.childrenNames[name] = struct{}{}
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
