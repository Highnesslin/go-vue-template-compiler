package main

import (
	"errors"
	"fmt"
	"strings"
)

func traverseChildren(node ASTNode) {
	for i := 0; i < len(node.children); i++ {
		child := node.children[i]
		Transform(child)
	}
}

func Transform(node interface{}) error {
	switch node := node.(type) {
	case ASTNode:
		{
			fmt.Println("ast", node)
			// 1. 处理 props
			for i, prop := range node.props {
				parts := strings.SplitN(prop.value, "=", 2)
				if len(parts) != 2 {
					return errors.New("文本解析出错: " + prop.value)
				}

				valueLen := len(parts[1])
				// 在Go中字符串是不可变数据类型, 一旦创建不可修改, 修改时实际上是创建了新的字符串副本, 所以修改时无需用指针
				node.props[i].field = parts[0]
				node.props[i].value = parts[1][1 : valueLen-1]
			}
			// 2. 处理子节点
			traverseChildren(node)
		}
	}
	// TODO 实际上这里还可以做很多事情, 比如特殊处理v-if、v-for、v-modal、v-once等等
	// ...

	return nil
}

// func main() {
// 	ast := ASTNode{
// 		ELEMENT,
// 		"div",
// 		[]Prop{Prop{"", `id="count"`, true}},
// 		[]interface{}{TextNode{
// 			INTERPOLATION,
// 			"count",
// 		}},
// 	}

// 	Transform(ast)

// 	fmt.Println("ast", ast)
// }
