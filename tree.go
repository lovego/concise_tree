package struct_tree

import "reflect"

var nodeType = reflect.TypeOf(Node{})
var ptr2nodeType = reflect.PtrTo(nodeType)

// 树节点
// 树的所有非叶子节点都应该嵌入此类型，所有叶子节点都应该是此类型。
type Node struct {
	name, code string
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Code() string {
	return n.code
}

func (n *Node) SetNameCode(name, code string) {
	n.name, n.code = name, code
}

type NodeIfc interface {
	Name() string
	Code() string
}

type Tree struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Children []Tree `json:"children,omitempty"`
}

func (t Tree) CodesMap() map[string]struct{} {
	m := make(map[string]struct{})
	t.setupCodesMap(m)
	return m
}

func (t Tree) setupCodesMap(m map[string]struct{}) {
	m[t.Code] = struct{}{}
	for _, child := range t.Children {
		child.setupCodesMap(m)
	}
}
