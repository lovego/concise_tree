package concise_tree

import (
	"errors"
	"reflect"
)

// NormalTree是一个正规形式的树结构。它经常被用在很多常见情况下。
//
// NormalTree is a tree structure in normal form. It's commonly used in most common cases.
type NormalTree struct {
	pathsMap      map[string]struct{}
	childrenPaths []string
	NormalTreeNode
}

// ToNormal converts a ConciseTree to a NormalTree.
func ToNormal(node ConciseTree) NormalTree {
	var tree NormalTree
	convert(&tree.NormalTreeNode, reflect.ValueOf(node))
	tree.init()
	return tree
}

// Keep return a new tree, keep only nodes that fn returns true.
// If fn returns false, the node and its decendants are all removed from the new tree.
func (t *NormalTree) Keep(fn func(NormalTreeNode) bool) NormalTree {
	if !fn(t.NormalTreeNode) {
		return NormalTree{}
	}
	tree := NormalTree{NormalTreeNode: t.keep(fn)}
	tree.init()
	return tree
}

func (t *NormalTree) PathsMap() map[string]struct{} {
	return t.pathsMap
}

func (t *NormalTree) CheckPaths(paths []string) error {
	for _, path := range paths {
		if _, ok := t.pathsMap[path]; !ok {
			return errors.New("unknown path: " + path)
		}
	}
	return nil
}

func (t *NormalTree) CleanPaths(paths []string) []string {
	j := 0
	for _, path := range paths {
		if _, ok := t.pathsMap[path]; ok {
			paths[j] = path
			j++
		}
	}
	return paths[:j]
}

func (t *NormalTree) ChildrenPaths() []string {
	return t.childrenPaths
}

func (t *NormalTree) init() {
	t.pathsMap = make(map[string]struct{})
	t.setupPathsMap(t.pathsMap)
	t.childrenPaths = t.NormalTreeNode.ChildrenPaths()
}
