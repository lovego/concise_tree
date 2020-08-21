package concise_tree_test

import (
	"fmt"

	tree "github.com/lovego/concise_tree"
)

type Modules struct {
	*tree.Node

	Bill *struct {
		*tree.Node
		List   *tree.Node `name:"列表" desc:"单据列表"`
		Detail *tree.Node `name:"详情" desc:"单据详情"`
	} `name:"单据" desc:"各种单据"`

	Goods *struct {
		*tree.Node
		Create *tree.Node `name:"创建" desc:"商品创建"`
		Update *tree.Node `name:"更新" desc:"商品更新"`
		Delete *tree.Node `name:"删除" desc:"商品删除"`
	} `name:"商品" desc:"商品（含库存）"`
}

func ExampleSetup() {
	modules := &Modules{}
	tree.Setup(modules, "权限树", "", "根节点")

	var n tree.ConciseTree
	n = modules
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Bill
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Bill.List
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Bill.Detail
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Goods
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Goods.Create
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Goods.Update
	fmt.Println(n.Name(), n.Code(), n.Desc())

	n = modules.Goods.Delete
	fmt.Println(n.Name(), n.Code(), n.Desc())

	// Output:
	// 权限树  根节点
	// 单据 bill 各种单据
	// 列表 bill.list 单据列表
	// 详情 bill.detail 单据详情
	// 商品 goods 商品（含库存）
	// 创建 goods.create 商品创建
	// 更新 goods.update 商品更新
	// 删除 goods.delete 商品删除
}
