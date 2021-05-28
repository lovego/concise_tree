package concise_tree

import (
	"log"
	"reflect"

	"github.com/lovego/struct_tag"
)

type nodeInfo struct {
	value         reflect.Value
	field         reflect.StructField
	path          string
	tags          map[string]string
	isNode        bool
	childrenCodes map[string]struct{}
	childrenNames map[string]struct{}
}

func (ni *nodeInfo) setup(mustBeNode bool, anonymousPath string) {
	if ni.value.Kind() == reflect.Ptr {
		if ni.value.IsNil() {
			if ni.value.CanSet() {
				ni.value.Set(reflect.New(ni.value.Type().Elem()))
			} else {
				log.Panicf(`field "%s" CanSet: false`, ni.path+anonymousPath)
			}
		}
		ni.value = ni.value.Elem()
	}
	if ni.value.Kind() != reflect.Struct {
		log.Panicf(`node "%s" should be struct, not %v`, ni.path+anonymousPath, ni.value.Kind())
	}
	ni.setupChildren(anonymousPath)
	if mustBeNode && !ni.isNode {
		log.Panicf(`node "%s" should anonymously embed concise_tree.Node`, ni.path+anonymousPath)
	}
}

func (ni *nodeInfo) setupChildren(anonymousPath string) {
	for i := 0; i < ni.value.NumField(); i++ {
		child := &nodeInfo{
			value: ni.value.Field(i),
			field: ni.value.Type().Field(i),
		}
		child.tags = struct_tag.Parse(string(child.field.Tag))

		switch child.value.Type() {
		case nodeType, ptr2nodeType:
			child.setupAsLeaf(ni)
		default:
			child.setupAsNonleaf(ni, anonymousPath)
		}
	}
}

func (ni *nodeInfo) setPath(parent *nodeInfo, validate bool) {
	code := ni.tags["code"]
	if code == "" {
		code = lowerFirstByte(ni.field.Name)
	}

	if parent.path != `` {
		ni.path = parent.path + `.` + code
	} else {
		ni.path = code
	}
	if !validate {
		return
	}

	if ni.tags["name"] == `` {
		log.Panicf(`name of node "%s" is empty`, ni.path)
	}
	parent.validateChild(code, ni.tags["name"])
}

func (ni *nodeInfo) setupAsLeaf(parent *nodeInfo) {
	if !exported(ni.field.Name) {
		ni.setPath(parent, false)
		log.Panicf(`node "%s" should be exported`, ni.path) // 非导出的设置不了path、tags
	}
	if ni.value.Kind() != reflect.Ptr {
		ni.value = ni.value.Addr()
	}
	if ni.value.IsNil() {
		ni.value.Set(reflect.New(nodeType))
	}
	if ni.field.Anonymous {
		parent.isNode = true // only parent need this
		ni.value.Interface().(*Node).Set(parent.path, parent.tags)
	} else {
		ni.setPath(parent, true)
		ni.value.Interface().(*Node).Set(ni.path, ni.tags)
	}
}

func (ni *nodeInfo) setupAsNonleaf(parent *nodeInfo, anonymousPath string) {
	if ni.field.Anonymous && ni.tags["name"] == "" {
		// 匿名嵌入且节点名称为空，只用来做类型共享，所以继续使用parent的信息：
		// 除了value更新为child的value，parent的field字段不会被使用外，其他字段都必须用parent的值。
		parentValue := parent.value
		parent.value = ni.value
		parent.setup(false, anonymousPath+"."+ni.field.Name)
		// 恢复parent的value，因为在上层setupChildren循环中还在继续使用。
		parent.value = parentValue
	} else if exported(ni.field.Name) {
		// 其余的导出字段都应该是树节点
		ni.setPath(parent, true)
		ni.setup(true, "")
	} else {
		ni.setPath(parent, false)
		log.Panicf(`tree node "%s" must be exported`, ni.path)
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
