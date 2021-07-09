package concise_tree

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

// NormalTree是一个正规形式的树结构。它经常被用在很多常见情况下。
//
// NormalTree is a tree structure in normal form. It's commonly used in most common cases.
type NormalTree struct {
	pathsMap              map[string]*NormalTreeNode
	excludingPaths        []string
	childrenPaths         []string
	expandedChildrenPaths []string
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
	kept, removedPaths := t.keep(fn)
	tree := NormalTree{NormalTreeNode: kept, excludingPaths: removedPaths}
	tree.init()
	return tree
}

func (t *NormalTree) PathsMap() map[string]*NormalTreeNode {
	return t.pathsMap
}

// CheckPaths return an error if has unknown path
func (t *NormalTree) CheckPaths(paths []string) error {
	for _, path := range paths {
		if _, ok := t.pathsMap[path]; !ok {
			return errors.New("unknown path: " + path)
		}
	}
	return nil
}

// CleanPaths remove paths that's not in the normal tree.
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

func (t *NormalTree) ExcludingPaths() []string {
	return t.excludingPaths
}

// ExpandPaths expand paths so that they are not ancestor of any excluding path.
func (t *NormalTree) ExpandPaths(paths []string) (result []string) {
	for _, path := range paths {
		if node := t.pathsMap[path]; node != nil {
			result = append(result, node.ExpandPath(t.excludingPaths)...)
		}
	}
	return
}

func (t *NormalTree) ChildrenPaths() []string {
	return t.childrenPaths
}

func (t *NormalTree) ExpandedChildrenPaths() []string {
	return t.expandedChildrenPaths
}

func (t *NormalTree) init() {
	t.pathsMap = make(map[string]*NormalTreeNode)
	t.setupPathsMap(t.pathsMap)
	t.childrenPaths = t.NormalTreeNode.ChildrenPaths()
	t.expandedChildrenPaths = t.ExpandPaths(t.childrenPaths)
}

// Belongs return true if path or any ancestor of path is included in pathsMap.
func Belongs(path string, pathsMap map[string]struct{}) bool {
	for {
		if _, ok := pathsMap[path]; ok {
			return true
		}
		if i := strings.LastIndexByte(path, '.'); i > 0 {
			path = path[:i]
		} else {
			return false
		}
	}
}

// remove duplicate paths that belongs to another one in paths.
func RemoveDuplicatePaths(paths []string) []string {
	m := make(map[string]struct{})
	sort.Strings(paths)
	for _, path := range paths {
		if !Belongs(path, m) {
			m[path] = struct{}{}
		}
	}
	if len(paths) == len(m) {
		return paths
	}
	var result = make([]string, len(m))
	i := 0
	for path := range m {
		result[i] = path
		i++
	}
	sort.Strings(result)
	return result
}
