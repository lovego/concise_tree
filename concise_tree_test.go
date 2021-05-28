package concise_tree_test

import (
	"fmt"

	tree "github.com/lovego/concise_tree"
)

type cud struct {
	Create *tree.Node `name:"创建" desc:"商品创建"`
	Update *tree.Node `name:"更新" desc:"商品更新"`
	Delete *tree.Node `name:"删除" desc:"商品删除"`
}

type Modules struct {
	*tree.Node

	Bill *struct {
		*tree.Node
		List   *tree.Node `name:"列表" desc:"单据列表"`
		Detail *tree.Node `name:"详情" desc:"单据详情"`
	} `name:"单据" desc:"各种单据"`

	Goods *struct {
		tree.Node
		cud
	} `name:"商品" desc:"商品（含库存）"`
}

func ExampleSetup() {
	modules := &Modules{}
	tree.Setup(modules, "", map[string]string{"name": "根节点"})

	var n tree.ConciseTree
	n = modules
	fmt.Println(n.Path(), n.Tags())

	n = modules.Bill
	fmt.Println(n.Path(), n.Tags())

	n = modules.Bill.List
	fmt.Println(n.Path(), n.Tags())

	n = modules.Bill.Detail
	fmt.Println(n.Path(), n.Tags())

	n = modules.Goods
	fmt.Println(n.Path(), n.Tags())

	n = modules.Goods.Create
	fmt.Println(n.Path(), n.Tags())

	n = modules.Goods.Update
	fmt.Println(n.Path(), n.Tags())

	n = modules.Goods.Delete
	fmt.Println(n.Path(), n.Tags())

	// Output:
	// map[name:根节点]
	// bill map[desc:各种单据 name:单据]
	// bill.list map[desc:单据列表 name:列表]
	// bill.detail map[desc:单据详情 name:详情]
	// goods map[desc:商品（含库存） name:商品]
	// goods.create map[desc:商品创建 name:创建]
	// goods.update map[desc:商品更新 name:更新]
	// goods.delete map[desc:商品删除 name:删除]
}
