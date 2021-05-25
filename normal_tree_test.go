package concise_tree_test

import (
	"encoding/json"
	"fmt"

	tree "github.com/lovego/concise_tree"
)

func ExampleToNormal() {
	modules := &Modules{}
	tree.Setup(modules, "", map[string]string{"name": "根节点"})
	normalTree := tree.ToNormal(modules)

	if b, err := json.MarshalIndent(normalTree, "", "  "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}

	if b, err := json.MarshalIndent(normalTree.PathsMap(), "", "  "); err != nil {
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
	// {
	//   "": {},
	//   "bill": {},
	//   "bill.detail": {},
	//   "bill.list": {},
	//   "goods": {},
	//   "goods.create": {},
	//   "goods.delete": {},
	//   "goods.update": {}
	// }
}
