[English](Readme.md) | 简体中文

# concise\_tree
用一种简洁的方式定义一棵简洁的树。
只支持静态树，这意味着树结构在编译时就固化了。

[![Build Status](https://github.com/lovego/concise_tree/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/concise_tree/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/concise_tree/badge.svg?branch=master&1)](https://coveralls.io/github/lovego/concise_tree)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/concise_tree)](https://goreportcard.com/report/github.com/lovego/concise_tree)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/concise_tree)](https://pkg.go.dev/github.com/lovego/concise_tree@v0.0.3)

## 特征
- 树结构从一个结构体类型构建而来，所有需要做的就是定义一个结构体。
  无需对结构体进行赋值，因此不用把每个字段名写两次。
- 所有的树节点都通过结构体字段定义和引用。因此在编译时就保证了节点引用的正确。
- 没有中间的"Children"字段。因此树定义和节点引用都更简洁。
- `concise_tree.Setup()` 自动设置好所有树节点的代码、名称和描述。
  代码是将祖先节点和本节点的字段名用"."连接起来的字符串；
  名称从结构体字段的`name`标签取; 描述从结构体字段的`desc`标签取。
- 将简洁树转换为含"Children"数组的普通树。

## 安装
`$ go get github.com/lovego/concise_tree`

## 示例
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

## 文档
[https://godoc.org/github.com/lovego/concise\_tree](https://godoc.org/github.com/lovego/concise_tree)
