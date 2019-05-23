package struct_tree

import (
	"reflect"
)

func Convert(node NodeIfc) Tree {
	var tree Tree
	convert(&tree, reflect.ValueOf(node))
	return tree
}

func convert(tree *Tree, node reflect.Value) {
	if node.Kind() == reflect.Ptr {
		node = node.Elem()
	}

	typ := node.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type {
		case nodeType:
			convertCodeName(tree, field, node.Field(i).Addr())
		case ptr2nodeType:
			convertCodeName(tree, field, node.Field(i))
		default:
			convertField(tree, field, node.Field(i))
		}
	}
}

func convertCodeName(tree *Tree, field reflect.StructField, value reflect.Value) {
	node := value.Interface().(*Node)
	if field.Anonymous {
		tree.Name, tree.Code = node.Name(), node.Code()
	} else {
		tree.Children = append(tree.Children, Tree{Name: node.Name(), Code: node.Code()})
	}
}

func convertField(tree *Tree, field reflect.StructField, value reflect.Value) {
	if exported(field.Name) {
		// 导出的字段都应该是树节点
		child := Tree{}
		convert(&child, value)
		tree.Children = append(tree.Children, child)
	} else if field.Anonymous {
		// 非导出的匿名字段不应该是树节点，只用来做类型共享，所以继续使用当前的name、code
		convert(tree, value)
	}
}
