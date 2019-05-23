package struct_tree

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
