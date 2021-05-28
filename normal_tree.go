package concise_tree

import (
	"reflect"
)

// NormalTree是一个正规形式的树结构。它经常被用在很多常见情况下。
//
// NormalTree is a tree structure in normal form. It's commonly used in most common cases.
type NormalTree struct {
	paths map[string]struct{}
	NormalTreeNode
}

// ToNormal converts a ConciseTree to a NormalTree.
func ToNormal(node ConciseTree) NormalTree {
	var tree NormalTree
	convert(&tree.NormalTreeNode, reflect.ValueOf(node))
	tree.paths = make(map[string]struct{})
	tree.setupPathsMap(tree.paths)
	return tree
}

// Keep return a new tree, keep only nodes that fn returns true.
// If fn returns false, the node and its decendants are all removed from the new tree.
func (t NormalTree) Keep(fn func(NormalTreeNode) bool) NormalTree {
	if !fn(t.NormalTreeNode) {
		return NormalTree{}
	}
	tree := NormalTree{NormalTreeNode: t.keep(fn)}
	tree.paths = make(map[string]struct{})
	tree.setupPathsMap(tree.paths)
	return tree
}

func (t NormalTree) PathsMap() map[string]struct{} {
	return t.paths
}
