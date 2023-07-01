package main

import (
	"errors"
	"fmt"
	"strings"
)

func traverseChildren(node *ASTNode) {
	for i := 0; i < len(node.children); i++ {
		child := node.children[i]
		Transform(child)
	}
}

func Transform(node interface{}) error {
	switch draft := node.(type) {
	case *ASTNode:
		{
			// 1. 处理 props
			for i, prop := range draft.props {
				parts := strings.SplitN(prop.value, "=", 2)
				if len(parts) != 2 {
					return errors.New("文本解析出错: " + prop.value)
				}

				field := parts[0]
				value := parts[1]

				// 在Go中字符串是不可变数据类型, 一旦创建不可修改, 修改时实际上是创建了新的字符串副本, 所以修改时无需用指针
				if draft.props[i].static == true {
					draft.props[i].field = field
					draft.props[i].value = value
				} else {
					// 动态字符串
					draft.props[i].field = field[1:(len(field))]
					// 掐头去尾
					draft.props[i].value = fmt.Sprintf("vm.%v", value[1:len(value)-1])
				}
			}
			// 2. 处理子节点
			traverseChildren(draft)
		}
	case *TextNode:
		{
			if draft.tagType == INTERPOLATION {
				// 如果是"插值表达式", 则从 vm 读取
				draft.content = fmt.Sprintf("vm.%v", draft.content)
			} else {
				// 如果是"纯文本", 则添加引号
				draft.content = fmt.Sprintf("'%v'", draft.content)
			}
		}
	}
	// TODO 实际上这里还可以做很多事情, 比如特殊处理v-if、v-for、v-modal、v-once等等
	// ...

	return nil
}
