English | [简体中文](Readme_cn.md)

# concise\_tree
A concise way to define a concise tree.
Only used for static tree, that means tree structure is frozed at compile time.

[![Build Status](https://github.com/lovego/concise_tree/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/concise_tree/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/concise_tree/badge.svg?branch=master&1)](https://coveralls.io/github/lovego/concise_tree)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/concise_tree)](https://goreportcard.com/report/github.com/lovego/concise_tree)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/concise_tree)](https://pkg.go.dev/github.com/lovego/concise_tree@v0.0.6)

## Feature
- A tree is constructed from a struct type, all need to do is define a struct.
  No struct value assignment is required, so you don't need to write every field name twice.
- All tree nodes are defined and referenced by struct fields. So node reference is ensured right at compile time.
- No intermediate "Children" fields. So tree definition and node reference are more concise.
- `concise_tree.Setup()` set up path and tags of all tree nodes automatically.
  Path is generated by concating field names of `ancestor nodes and self node` by "." ;
  Tags is parsed from struct field's tag;
- Convert a concise tree to a normal tree with `Children` slice.

## Install
`$ go get github.com/lovego/concise_tree`

## Example
```go
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
```



## Documentation
[https://godoc.org/github.com/lovego/concise\_tree](https://godoc.org/github.com/lovego/concise_tree)
