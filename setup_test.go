package struct_tree_test

import (
	"fmt"

	"github.com/lovego/struct_tree"
)

func ExampleSetup() {
	type Write struct {
		struct_tree.Node
		Create struct_tree.Node  `新增`
		Update *struct_tree.Node `编辑`
		Delete *struct_tree.Node `删除`
	}
	type other struct {
		Audit struct_tree.Node `审核`
	}
	type ReadWrite struct {
		*struct_tree.Node
		Read   *struct_tree.Node `读取`
		*Write `写入`
		other
	}
	var root = new(ReadWrite)

	struct_tree.Setup(root, "根节点", "")
	for _, n := range []interface {
		Name() string
		Code() string
	}{
		root,
		root.Read,
		root.Write,
		&root.Create,
		root.Update,
		root.Delete,
		&root.Audit,
	} {
		fmt.Printf("`%s`: %s\n", n.Code(), n.Name())
	}

	// Output:
	// ``: 根节点
	// `Read`: 读取
	// `Write`: 写入
	// `Write.Create`: 新增
	// `Write.Update`: 编辑
	// `Write.Delete`: 删除
	// `Audit`: 审核
}
