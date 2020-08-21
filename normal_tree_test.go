package concise_tree_test

import (
	"encoding/json"
	"fmt"

	tree "github.com/lovego/concise_tree"
)

func ExampleToNormal() {
	modules := &Modules{}
	tree.Setup(modules, "权限树", "", "根节点")
	normalTree := tree.ToNormal(modules)

	if b, err := json.MarshalIndent(normalTree, "", "  "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}

	if b, err := json.MarshalIndent(normalTree.CodesMap(), "", "  "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// {
	//   "name": "权限树",
	//   "code": "",
	//   "desc": "根节点",
	//   "children": [
	//     {
	//       "name": "单据",
	//       "code": "bill",
	//       "desc": "各种单据",
	//       "children": [
	//         {
	//           "name": "列表",
	//           "code": "bill.list",
	//           "desc": "单据列表"
	//         },
	//         {
	//           "name": "详情",
	//           "code": "bill.detail",
	//           "desc": "单据详情"
	//         }
	//       ]
	//     },
	//     {
	//       "name": "商品",
	//       "code": "goods",
	//       "desc": "商品（含库存）",
	//       "children": [
	//         {
	//           "name": "创建",
	//           "code": "goods.create",
	//           "desc": "商品创建"
	//         },
	//         {
	//           "name": "更新",
	//           "code": "goods.update",
	//           "desc": "商品更新"
	//         },
	//         {
	//           "name": "删除",
	//           "code": "goods.delete",
	//           "desc": "商品删除"
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
