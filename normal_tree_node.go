package concise_tree

import (
	"reflect"

	"github.com/lovego/struct_tag"
)

type NormalTreeNode struct {
	Path     string            `json:"path"`
	Tags     map[string]string `json:"tags,omitempty"`
	Children []NormalTreeNode  `json:"children,omitempty"`
}

func (t NormalTreeNode) setupPathsMap(m map[string]struct{}) {
	m[t.Path] = struct{}{}
	for _, child := range t.Children {
		child.setupPathsMap(m)
	}
}

func (t NormalTreeNode) keep(fn func(NormalTreeNode) bool) NormalTreeNode {
	var tree NormalTreeNode
	tree.Path, tree.Tags = t.Path, t.Tags
	for _, child := range t.Children {
		if fn(child) {
			tree.Children = append(tree.Children, child.keep(fn))
		}
	}
	return tree
}

func convert(tree *NormalTreeNode, node reflect.Value) {
	if node.Kind() == reflect.Ptr {
		node = node.Elem()
	}

	typ := node.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type {
		case nodeType, ptr2nodeType:
			convertLeafNode(tree, field, node.Field(i))
		default:
			convertNonleafNode(tree, field, node.Field(i))
		}
	}
}

func convertLeafNode(tree *NormalTreeNode, field reflect.StructField, value reflect.Value) {
	if value.Kind() != reflect.Ptr {
		value = value.Addr()
	}
	node := value.Interface().(*Node)
	if field.Anonymous {
		tree.Path, tree.Tags = node.Path(), node.Tags()
	} else {
		tree.Children = append(tree.Children, NormalTreeNode{Path: node.Path(), Tags: node.Tags()})
	}
}

func convertNonleafNode(tree *NormalTreeNode, field reflect.StructField, value reflect.Value) {
	if field.Anonymous && struct_tag.Get(string(field.Tag), "name") == `` {
		// 匿名嵌入且节点名称为空，只用来做类型共享
		convert(tree, value)
	} else if exported(field.Name) {
		// 其余的导出字段都应该是树节点
		child := NormalTreeNode{}
		convert(&child, value)
		tree.Children = append(tree.Children, child)
	}
}
