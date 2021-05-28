package concise_tree

import (
	"fmt"
	"reflect"
)

func ExampleSetup_panic() {
	defer func() {
		fmt.Println(recover())
	}()
	Setup(struct{ *Node }{}, "", nil)

	// Output:
	// tree should be a pointer, not struct
}

func ExampleNodeInfo_setup_panic1() {
	defer func() {
		fmt.Println(recover())
	}()
	(&nodeInfo{
		value: reflect.ValueOf((*int)(nil)),
		path:  "abc",
	}).setup(true, "")

	// Output:
	// field "abc" CanSet: false
}

func ExampleNodeInfo_setup_panic2() {
	defer func() {
		fmt.Println(recover())
	}()
	(&nodeInfo{path: "abc"}).setup(true, "")

	// Output:
	// node "abc" should be struct, not invalid
}

func ExampleNodeInfo_setup_panic3() {
	defer func() {
		fmt.Println(recover())
	}()
	(&nodeInfo{
		value: reflect.ValueOf(&struct{}{}),
		path:  "abc",
	}).setup(true, "")

	// Output:
	// node "abc" should anonymously embed concise_tree.Node
}

func ExampleNodeInfo_setPath_panic() {
	defer func() {
		fmt.Println(recover())
	}()
	(&nodeInfo{
		field: reflect.StructField{Name: "Test"},
		tags:  make(map[string]string),
	}).setPath(&nodeInfo{path: "abc"}, true)

	// Output:
	// name of node "abc.test" is empty
}

func ExampleNodeInfo_setupAsLeaf() {
	(&nodeInfo{
		value: reflect.ValueOf(&Node{}).Elem(),
		field: reflect.StructField{Name: "Test"},
		tags:  map[string]string{"name": "test"},
	}).setupAsLeaf(&nodeInfo{path: "abc"})

	// Output:
}

func ExampleNodeInfo_setupAsLeaf_panic() {
	defer func() {
		fmt.Println(recover())
	}()
	(&nodeInfo{
		field: reflect.StructField{Name: "test"},
		tags:  map[string]string{"name": "test"},
	}).setupAsLeaf(&nodeInfo{path: "abc"})

	// Output:
	// node "abc.test" should be exported
}

func ExampleNodeInfo_setupAsNonleaf() {
	tree := struct{ Node }{}
	(&nodeInfo{
		value: reflect.ValueOf(&tree).Elem(),
		field: reflect.StructField{Name: "Test", Anonymous: true},
		tags:  map[string]string{},
	}).setupAsNonleaf(&nodeInfo{path: "abc"}, "")

	fmt.Println(tree.Path())
	// Output:
	// abc
}

func ExampleNodeInfo_setupAsNonleaf_panic() {
	defer func() {
		fmt.Println(recover())
	}()

	(&nodeInfo{
		field: reflect.StructField{Name: "test"},
		tags:  map[string]string{},
	}).setupAsNonleaf(&nodeInfo{path: "abc"}, "")

	// Output:
	// tree node "abc.test" must be exported
}

func ExampleNodeInfo_validateChild_sameCode() {
	defer func() {
		fmt.Println(recover())
	}()
	type anonymous struct {
		A *Node `name:"testB"`
	}
	Setup(&struct {
		*Node
		A *Node `name:"testA"`
		anonymous
	}{}, "root", nil)

	// Output:
	// node "root" has children node of same code "a"
}

func ExampleNodeInfo_validateChild_sameName() {
	defer func() {
		fmt.Println(recover())
	}()
	Setup(&struct {
		*Node
		A *Node `name:"test"`
		B *Node `name:"test"`
	}{}, "root", nil)

	// Output:
	// node "root" has children node of same name "test"
}
