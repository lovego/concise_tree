package concise_tree

import (
	"reflect"

	"github.com/lovego/struct_tag"
)

// NormalTree是一个正规形式的树结构。它经常被用在很多常见情况下。
//
// NormalTree is a tree structure in normal form. It's commonly used in most common cases.
type NormalTree struct {
	Path     string            `json:"path"`
	Tags     map[string]string `json:"tags,omitempty"`
	Children []NormalTree      `json:"children,omitempty"`
}

func (t NormalTree) PathsMap() map[string]struct{} {
	m := make(map[string]struct{})
	t.setupPathsMap(m)
	return m
}

func (t NormalTree) setupPathsMap(m map[string]struct{}) {
	m[t.Path] = struct{}{}
	for _, child := range t.Children {
		child.setupPathsMap(m)
	}
}

// Keep return a new tree, keep only nodes that fn returns true.
// If fn returns false, the node and its decendants are all removed from the new tree.
func (t NormalTree) Keep(fn func(NormalTree) bool) NormalTree {
	if !fn(t) {
		return NormalTree{}
	}
	return t.keep(fn)
}

func (t NormalTree) keep(fn func(NormalTree) bool) NormalTree {
	var tree NormalTree
	tree.Path, tree.Tags = t.Path, t.Tags
	for _, child := range t.Children {
		if fn(child) {
			tree.Children = append(tree.Children, child.Keep(fn))
		}
	}
	return tree
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
		case nodeType, ptr2nodeType:
			convertLeafNode(tree, field, node.Field(i))
		default:
			convertNonleafNode(tree, field, node.Field(i))
		}
	}
}

func convertLeafNode(tree *NormalTree, field reflect.StructField, value reflect.Value) {
	if value.Kind() != reflect.Ptr {
		value = value.Addr()
	}
	node := value.Interface().(*Node)
	if field.Anonymous {
		tree.Path, tree.Tags = node.Path(), node.Tags()
	} else {
		tree.Children = append(tree.Children, NormalTree{Path: node.Path(), Tags: node.Tags()})
	}
}

func convertNonleafNode(tree *NormalTree, field reflect.StructField, value reflect.Value) {
	if field.Anonymous && struct_tag.Get(string(field.Tag), "name") == `` {
		// 匿名嵌入且节点名称为空，只用来做类型共享
		convert(tree, value)
	} else if exported(field.Name) {
		// 其余的导出字段都应该是树节点
		child := NormalTree{}
		convert(&child, value)
		tree.Children = append(tree.Children, child)
	}
}
