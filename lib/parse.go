/**
 * Parse
 */
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	// "github.com/oleiade/lane"
)

type TYPE string

type Prop struct {
	field  string
	value  string
	static bool
}

type ASTNode struct {
	tagType  TYPE          // 类型: ROOT标签 元素 纯文本 插值表达式 IF标签 FOR标签
	tagName  string        // 标签名
	props    []*Prop       // props
	children []interface{} // []*ASTNode | []*TextNode
}

type TextNode struct {
	tagType TYPE // 纯文本 插值表达式
	content string
}

/**
 * 枚举
 */
const (
	ROOT          TYPE = "ROOT"          // 根节点
	ELEMENT       TYPE = "ELEMENT"       // 元素
	PLAIN_TEXT    TYPE = "PLAIN_TEXT"    // 纯文本
	INTERPOLATION TYPE = "INTERPOLATION" // 插值表达式
	IF_NODE       TYPE = "INTERPOLATION" // if标签
	FOR_NODE      TYPE = "INTERPOLATION" // for标签
)

/**
 * 转Token用到的正则表达式
 */
var (
	START_TAG_PATTERN       = regexp.MustCompile(`^<([a-zA-Z]+)`)
	START_TAG_CLOSE_PATTERN = regexp.MustCompile(`^\s*(\/?)>`)
	END_TAG_PATTERN         = regexp.MustCompile(`^<\/([a-zA-Z]+)>`)
	// DYNAMIC_PROPS_PATTERN   = regexp.MustCompile(`^\s*((?:v-[\w-]+:|@|:|#)\[[^=]+?\][^\s"'<>\/=]*)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>` + "`" + `]+)))?`)
	DYNAMIC_PROPS_PATTERN = regexp.MustCompile(`^\s*([^\s"'<>\/=]+)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>` + "`" + `]+)))?`)
	STATIC_PROPS_PATTERN  = regexp.MustCompile(`^\s*([^\s"':<>\/=]+)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>` + "`" + `]+)))?`)
	INTERPOLATION_PATTERN = regexp.MustCompile(`^{{\s*(.*?)\s*}}`)
)

/**
 * 将代码转成AST
 */
func Parse(template string) ASTNode {
	stack := []ASTNode{}

	/**
	 * 去除空白符
	 */
	shakingSpacing := func() {
		pattern := regexp.MustCompile(`^\s+`)
		for {
			match := pattern.FindStringSubmatch(template)

			if len(match) == 0 {
				break
			}
			template = template[len(match[0]):]
		}
	}

	/**
	 * 截断字符串
	 */
	advanceBy := func(pos int) {
		template = template[pos:]
	}

	/**
	 * 入栈
	 */
	pushStack := func(node ASTNode) {
		stack = append(stack, node)
	}

	/**
	 * 出栈
	 */
	popStack := func() (ASTNode, error) {
		stackLen := len(stack)
		if stackLen > 1 {
			last := stack[stackLen-1]
			stack = stack[:stackLen-1]
			return last, nil
		}
		return ASTNode{}, errors.New("Token栈已空")
	}

	/**
	 * 向栈尾元素追加子节点
	 */
	pushChildrenToStackTail := func(arg interface{}) {
		stackLen := len(stack)
		if stackLen > 0 {
			stack[stackLen-1].children = append(stack[stackLen-1].children, arg)
		}
	}

	/**
	 * 解析开始标签
	 */
	parseStartTag := func() (ASTNode, error) {
		node := ASTNode{}

		match := START_TAG_PATTERN.FindStringSubmatch(template)
		if len(match) == 0 {
			return node, errors.New("开始标签匹配失败")
		}

		// 1. Tag
		// token = match[1]
		if len(stack) == 0 {
			node.tagType = ROOT
		} else {
			node.tagType = ELEMENT
		}

		node.tagName = match[1]
		advanceBy(len(match[0]))

		for {
			shakingSpacing()

			// 2. 关闭标签 >
			match = START_TAG_CLOSE_PATTERN.FindStringSubmatch(template)
			if len(match) > 0 {
				advanceBy(len(match[0]))

				if len(match[0]) == 2 {
					// 自闭标签 />
					// token = node.tagName
					// 收集 ASTNode 的信息
					pushChildrenToStackTail(&node)
				} else {
					pushStack(node)
				}
				break
			}

			// 3. props [动态 和 静态]
			// 	  3.1. 静态 props
			match = STATIC_PROPS_PATTERN.FindStringSubmatch(template)

			if len(match) > 0 {
				// token = match[0]
				node.props = append(node.props, &Prop{"", match[0], true})
				advanceBy(len(match[0]))
			}

			// 	  3.2. 动态 props
			match = DYNAMIC_PROPS_PATTERN.FindStringSubmatch(template)

			if len(match) > 0 {
				// token = match[0]
				node.props = append(node.props, &Prop{"", match[0], false})
				advanceBy(len(match[0]))
			}
		}

		return node, nil
	}

	/**
	 * 解析文本
	 */
	parseText := func() {
		endTokens := []string{"<", "{{"}
		endIndex := len(template)

		for i := 0; i < len(endTokens); i++ {
			index := strings.Index((template)[1:], endTokens[i])
			if index != -1 && endIndex > index {
				endIndex = index
			}
		}

		text := template[0:endIndex]

		if len(text) > 0 {
			// TODO: 临时+1
			// token = text
			advanceBy(len(text) + 1)

			pushChildrenToStackTail(&TextNode{PLAIN_TEXT, text})
		}
	}

	/**
	 * 解析插值表达式
	 */
	parseInterpolation := func() {
		match := INTERPOLATION_PATTERN.FindStringSubmatch(template)
		if len(match) > 0 {
			// token = match[1]
			advanceBy(len(match[0]))

			pushChildrenToStackTail(&TextNode{INTERPOLATION, match[1]})
		}
	}

	/**
	 * Main
	 */
	for {
		shakingSpacing()

		if len(template) == 0 {
			break
		}

		if template[0] == '<' {
			// 1. 开始标签 <div
			if template[1] != '/' {
				parseStartTag()
			}

			// 2. 结束标签 </div>
			match := END_TAG_PATTERN.FindStringSubmatch(template)
			if len(match) > 0 {
				// token = match[1]
				advanceBy(len(match[0]))
				last, err := popStack()
				if err != nil {
					break
				} else {
					pushChildrenToStackTail(&last)
				}
			}

		} else if template[0] == '{' {
			// 3. 差值表达式 {{ xxx }}
			parseInterpolation()
		} else {
			// 4. 文本 plain
			parseText()
		}
	}

	return stack[len(stack)-1]
}

/**
 * transform
 */
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

/**
 * generate
 */
func nodeToCode(node interface{}) (string, error) {
	switch node := node.(type) {
	case *ASTNode:
		{
			children := ""

			for _, child := range node.children {
				cur, error := nodeToCode(child)
				if error != nil {
					return "", error
				}
				children += cur + ","
			}

			props := ""
			for _, prop := range node.props {
				props += fmt.Sprintf(" %s:%v,", prop.field, prop.value)
			}

			var result string
			if props != "" {
				result = fmt.Sprintf("%s, {%s}", result, props)
			}

			if children != "" {
				result = fmt.Sprintf("%s, [%s]", result, children)
			}

			return fmt.Sprintf("h('%s'%s)", node.tagName, result), nil
		}
	case *TextNode:
		{
			return node.content, nil
		}
	}

	return "", errors.New("不识别的标签")
}

func Generate(node *ASTNode) (string, error) {
	vNode, err := nodeToCode(node)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`
		import { h } from 'vue'
		export function render () {
			const vm = this
			return %v
		}
	`, vNode), nil
}

func main() {
	template := `
		<template>
			<div id="count">
				<span :className="position">{{ add }}</span>
				<span>
					number:{{ count }}
				</span>
				<img src="asdsad.png" />
			</div>
		</template>
	`

	ast := Parse(template)

	Transform(&ast)

	fmt.Println("ast", ast)

	code, error := Generate(&ast)

	if error != nil {
		fmt.Println("error", error)
	}

	fmt.Println("code", code)
}
