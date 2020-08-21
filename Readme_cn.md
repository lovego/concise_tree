[English](Readme.md) | 简体中文

# concise\_tree
用一个简洁的方式定义一棵简洁的树。
只支持静态树，意味着树结构在编译时就固化了。

[![Build Status](https://travis-ci.org/lovego/concise_tree.svg?branch=master)](https://travis-ci.org/lovego/concise_tree)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/concise_tree/master.svg)](https://coveralls.io/github/lovego/concise_tree?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/concise_tree)](https://goreportcard.com/report/github.com/lovego/concise_tree)
[![GoDoc](https://godoc.org/github.com/lovego/concise_tree?status.svg)](https://godoc.org/github.com/lovego/concise_tree)

## 特征
- 树结构从一个结构体类型构建而来，所有需要做的就是定义一个结构体。
  无需对结构体的进行赋值，因此不用把每个字段名写两次。
- 所有的树节点都通过结构体字段定义。因此在编译时就保证了节点引用的正确。
- 没有中间的"Children"字段。因此树定义和节点引用都更简洁。
- `concise_tree.Setup()` 自动设置好所有树节点的代码、名称和描述。
  代码是将祖先节点和本节点的字段名用"."连接起来的字符串；
  名称从结构体字段的`name`标签取; 描述从结构体字段的`desc`标签取。
- 将简洁树转换为含数组形式的"Children"字段的普通树。

## 安装
`$ go get github.com/lovego/concise_tree`

## 示例
```
```

## 文档
[https://godoc.org/github.com/lovego/concise\_tree](https://godoc.org/github.com/lovego/concise_tree)
