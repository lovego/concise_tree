package concise_tree

import (
	"log"
	"reflect"
)

var nodeType = reflect.TypeOf(Node{})
var ptr2nodeType = reflect.PtrTo(nodeType)

// ConciseTree use an interface that describe `Node` to represent a concise tree.
type ConciseTree interface {
	Path() string
	Tags() map[string]string
}

// Node 代表ConciseTree的一个节点，用来构造结一棵ConciseTree。
// 所有非叶子节点都应该嵌入此类型，所有叶子节点都应该直接使用此类型。
//
// Node represents a node of ConciseTree，and is used to construct a ConciseTree.
// All nonleaf nodes should embed this type, all leaf nodes should use this type directly.
type Node struct {
	path string
	tags map[string]string
}

func (n *Node) Path() string {
	return n.path
}

func (n *Node) Tags() map[string]string {
	return n.tags
}

func (n *Node) Set(path string, tags map[string]string) {
	n.path, n.tags = path, tags
}

// 设置树的所有节点的path、tags
func Setup(tree ConciseTree, path string, tags map[string]string) {
	value := reflect.ValueOf(tree)
	if value.Kind() != reflect.Ptr {
		log.Panicf(`tree should be a pointer, not %v`, value.Kind())
	}
	(&nodeInfo{
		value: value,
		path:  path,
		tags:  tags,
	}).setup()
}
