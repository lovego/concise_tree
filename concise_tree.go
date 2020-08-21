package concise_tree

import "reflect"

var nodeType = reflect.TypeOf(Node{})
var ptr2nodeType = reflect.PtrTo(nodeType)

// ConciseTree use an interface that describe `Node` to represent a concise tree.
type ConciseTree interface {
	Name() string
	Code() string
	Desc() string
}

// Node 代表ConciseTree的一个节点，用来构造结一棵ConciseTree。
// 所有非叶子节点都应该嵌入此类型，所有叶子节点都应该直接使用此类型。
//
// Node represents a node of ConciseTree，and is used to construct a ConciseTree.
// All nonleaf nodes should embed this type, all leaf nodes should use this type directly.
type Node struct {
	name, code, desc string
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Code() string {
	return n.code
}

func (n *Node) Desc() string {
	return n.desc
}

func (n *Node) Set(name, code, desc string) {
	n.name, n.code, n.desc = name, code, desc
}
