package concise_tree_test

import (
	"encoding/json"
	"fmt"
	"sort"

	tree "github.com/lovego/concise_tree"
)

var normalTree tree.NormalTree

func init() {
	modules := &Modules{}
	tree.Setup(modules, "", map[string]string{"name": "根节点"})
	normalTree = tree.ToNormal(modules)
}

func ExampleToNormal() {
	if b, err := json.MarshalIndent(normalTree, "", "  "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}
	// Output:
	// {
	//   "path": "",
	//   "tags": {
	//     "name": "根节点"
	//   },
	//   "children": [
	//     {
	//       "path": "bill",
	//       "tags": {
	//         "desc": "各种单据",
	//         "name": "单据"
	//       },
	//       "children": [
	//         {
	//           "path": "bill.list",
	//           "tags": {
	//             "desc": "单据列表",
	//             "name": "列表"
	//           }
	//         },
	//         {
	//           "path": "bill.detail",
	//           "tags": {
	//             "desc": "单据详情",
	//             "name": "详情"
	//           }
	//         }
	//       ]
	//     },
	//     {
	//       "path": "goods",
	//       "tags": {
	//         "desc": "商品（含库存）",
	//         "name": "商品"
	//       },
	//       "children": [
	//         {
	//           "path": "goods.create",
	//           "tags": {
	//             "desc": "商品创建",
	//             "name": "创建"
	//           }
	//         },
	//         {
	//           "path": "goods.update",
	//           "tags": {
	//             "desc": "商品更新",
	//             "name": "更新"
	//           }
	//         },
	//         {
	//           "path": "goods.delete",
	//           "tags": {
	//             "desc": "商品删除",
	//             "name": "删除"
	//           }
	//         }
	//       ]
	//     }
	//   ]
	// }
}

func ExampleNormalTree_PathsMap() {
	var pathNames = []string{}
	for path, node := range normalTree.PathsMap() {
		pathNames = append(pathNames, path+":"+node.Tags["name"])
	}
	sort.Strings(pathNames)
	fmt.Println(pathNames)

	// Output:
	// [:根节点 bill.detail:详情 bill.list:列表 bill:单据 goods.create:创建 goods.delete:删除 goods.update:更新 goods:商品]
}

func keepNameNotEqual(s string) {
	kept := normalTree.Keep(func(node tree.NormalTreeNode) bool {
		return node.Tags["name"] != s
	})
	fmt.Println(kept.ExcludingPaths())
	fmt.Println(kept.ExpandedChildrenPaths())
	if b, err := json.MarshalIndent(kept, "", "  "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}
}

func ExampleNormalTree_Keep_1() {
	keepNameNotEqual("删除")
	// Output:
	// [goods.delete]
	// [bill goods.create goods.update]
	// {
	//   "path": "",
	//   "tags": {
	//     "name": "根节点"
	//   },
	//   "children": [
	//     {
	//       "path": "bill",
	//       "tags": {
	//         "desc": "各种单据",
	//         "name": "单据"
	//       },
	//       "children": [
	//         {
	//           "path": "bill.list",
	//           "tags": {
	//             "desc": "单据列表",
	//             "name": "列表"
	//           }
	//         },
	//         {
	//           "path": "bill.detail",
	//           "tags": {
	//             "desc": "单据详情",
	//             "name": "详情"
	//           }
	//         }
	//       ]
	//     },
	//     {
	//       "path": "goods",
	//       "tags": {
	//         "desc": "商品（含库存）",
	//         "name": "商品"
	//       },
	//       "children": [
	//         {
	//           "path": "goods.create",
	//           "tags": {
	//             "desc": "商品创建",
	//             "name": "创建"
	//           }
	//         },
	//         {
	//           "path": "goods.update",
	//           "tags": {
	//             "desc": "商品更新",
	//             "name": "更新"
	//           }
	//         }
	//       ]
	//     }
	//   ]
	// }
}

func ExampleNormalTree_Keep_2() {
	keepNameNotEqual("商品")

	// Output:
	// [goods]
	// [bill]
	// {
	//   "path": "",
	//   "tags": {
	//     "name": "根节点"
	//   },
	//   "children": [
	//     {
	//       "path": "bill",
	//       "tags": {
	//         "desc": "各种单据",
	//         "name": "单据"
	//       },
	//       "children": [
	//         {
	//           "path": "bill.list",
	//           "tags": {
	//             "desc": "单据列表",
	//             "name": "列表"
	//           }
	//         },
	//         {
	//           "path": "bill.detail",
	//           "tags": {
	//             "desc": "单据详情",
	//             "name": "详情"
	//           }
	//         }
	//       ]
	//     }
	//   ]
	// }
}

func ExampleNormalTree_Keep_remove_root() {
	r := (&tree.NormalTree{
		NormalTreeNode: tree.NormalTreeNode{Path: "root"},
	}).Keep(func(node tree.NormalTreeNode) bool {
		return node.Path != "root"
	})
	fmt.Printf("%+v\n", r)

	// Output:
	// {pathsMap:map[] excludingPaths:[] childrenPaths:[] expandedChildrenPaths:[] NormalTreeNode:{Path: Tags:map[] Children:[]}}
}

func ExampleNormalTree_CheckPaths() {
	fmt.Println(normalTree.CheckPaths([]string{"bill", "goods", "goods.create"}))
	fmt.Println(normalTree.CheckPaths([]string{"goods", "goods.create", "goods.insert"}))

	// Output:
	// <nil>
	// unknown path: goods.insert
}

func ExampleNormalTree_CleanPaths() {
	fmt.Println(normalTree.CleanPaths([]string{"goods", "goods.create", "goods.insert"}))

	// Output:
	// [goods goods.create]
}

func ExampleNormalTree_ChildrenPaths() {
	fmt.Println(normalTree.ChildrenPaths())

	// Output:
	// [bill goods]
}

func ExampleBelongs() {
	var m = map[string]struct{}{
		"A":     struct{}{},
		"B.1":   struct{}{},
		"C.1.1": struct{}{},
	}
	fmt.Println(tree.Belongs("A", m), tree.Belongs("A.1", m), tree.Belongs("A.1.1", m), tree.Belongs("A1", m))
	fmt.Println(tree.Belongs("B", m), tree.Belongs("B.1", m), tree.Belongs("B.1.1", m), tree.Belongs("B.2", m))
	fmt.Println(tree.Belongs("C", m), tree.Belongs("C.1", m), tree.Belongs("C.1.1", m), tree.Belongs("C.1.2", m))
	fmt.Println(tree.Belongs("D", m), tree.Belongs("D.1", m), tree.Belongs("D.1.1", m), tree.Belongs("D.1.2", m))

	// Output:
	// true true true false
	// false true true false
	// false false true false
	// false false false false
}

func ExampleRemoveDuplicatePaths() {
	fmt.Println(tree.RemoveDuplicatePaths([]string{
		"B.1", "A", "A.1", "B.2", "A.2", "B.1.2", "C",
	}))
	fmt.Println(tree.RemoveDuplicatePaths([]string{
		"B.1", "A", "B.2", "C",
	}))
	// Output:
	// [A B.1 B.2 C]
	// [A B.1 B.2 C]
}
