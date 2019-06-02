package struct_tree_test

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/lovego/struct_tree"
)

type ReadWrite struct {
	*struct_tree.Node
	Read   *struct_tree.Node `读取`
	*Write `写入`
	Other
}
type Write struct {
	struct_tree.Node
	Create struct_tree.Node  `新增`
	Update *struct_tree.Node `编辑`
	Delete *struct_tree.Node `删除`
}
type Other struct {
	Audit struct_tree.Node `审核`
}

var root = new(ReadWrite)

func init() {
	struct_tree.Setup(root, "根节点", "")
}

func ExampleSetup() {
	for _, n := range []struct_tree.NodeIfc{
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
	// `read`: 读取
	// `write`: 写入
	// `write.create`: 新增
	// `write.update`: 编辑
	// `write.delete`: 删除
	// `audit`: 审核
}

func ExampleConvert() {
	tree := struct_tree.Convert(root)
	if b, err := json.MarshalIndent(tree, ``, ` `); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}

	var codes []string
	for k := range tree.CodesMap() {
		codes = append(codes, k)
	}
	sort.Strings(codes)
	fmt.Printf("%#v\n", codes)

	// Output:
	// {
	//  "name": "根节点",
	//  "code": "",
	//  "children": [
	//   {
	//    "name": "读取",
	//    "code": "read"
	//   },
	//   {
	//    "name": "写入",
	//    "code": "write",
	//    "children": [
	//     {
	//      "name": "新增",
	//      "code": "write.create"
	//     },
	//     {
	//      "name": "编辑",
	//      "code": "write.update"
	//     },
	//     {
	//      "name": "删除",
	//      "code": "write.delete"
	//     }
	//    ]
	//   },
	//   {
	//    "name": "审核",
	//    "code": "audit"
	//   }
	//  ]
	// }
	// []string{"", "audit", "read", "write", "write.create", "write.delete", "write.update"}
}
