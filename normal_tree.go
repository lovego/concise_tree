package concise_tree

import (
	"reflect"
)

// NormalTree是一个正规形式的树结构。它经常被用在很多常见情况下。
//
// NormalTree is a tree structure in normal form. It's commonly used in most common cases.
type NormalTree struct {
	Path     string            `json:"path"`
	Tags     map[string]string `json:"tags,omitempty"`
	Children []NormalTree      `json:"children,omitempty"`
}

func (t NormalTree) CodesMap() map[string]struct{} {
	m := make(map[string]struct{})
	t.setupCodesMap(m)
	return m
}

func (t NormalTree) setupCodesMap(m map[string]struct{}) {
	m[t.Path] = struct{}{}
	for _, child := range t.Children {
		child.setupCodesMap(m)
	}
}

// ToNormal converts a ConciseTree to a NormalTree.
func ToNormal(node ConciseTree) NormalTree {
	var tree NormalTree
	convert(&tree, reflect.ValueOf(node))
	return tree
}

func convert(tree *NormalTree, node reflect.Value) {
	if node.Kind() == reflect.Ptr {
		node = node.Elem()
	}

	typ := node.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type {
		case nodeType:
			convertLeafNode(tree, field, node.Field(i).Addr())
		case ptr2nodeType:
			convertLeafNode(tree, field, node.Field(i))
		default:
			convertNonLeafNode(tree, field, node.Field(i))
		}
	}
}

func convertLeafNode(tree *NormalTree, field reflect.StructField, value reflect.Value) {
	node := value.Interface().(*Node)
	if field.Anonymous {
		tree.Path, tree.Tags = node.Path(), node.Tags()
	} else {
		tree.Children = append(tree.Children, NormalTree{Path: node.Path(), Tags: node.Tags()})
	}
}

func convertNonLeafNode(tree *NormalTree, field reflect.StructField, value reflect.Value) {
	if field.Anonymous && field.Tag == `` {
		// 匿名嵌入且节点名称为空，只用来做类型共享
		convert(tree, value)
	} else if exported(field.Name) {
		// 其余的导出字段都应该是树节点
		child := NormalTree{}
		convert(&child, value)
		tree.Children = append(tree.Children, child)
	}
}
