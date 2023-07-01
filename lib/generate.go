package main

import (
	"errors"
	"fmt"
	"strings"
)

func nodeToCode(node interface{}) (string, error) {
	switch node := node.(type) {
	case *ASTNode:
		{
			var children []string

			for _, child := range node.children {
				cur, error := nodeToCode(child)
				if error != nil {
					return "", error
				}
				children = append(children, cur)
			}

			var props []string
			for _, prop := range node.props {
				props = append(props, fmt.Sprintf("%s:%v", prop.field, prop.value))
			}

			var result string
			if len(props) > 0 {
				result = fmt.Sprintf("%s, { %s }", result, strings.Join(props, ", "))
			}

			if len(children) > 0 {
				result = fmt.Sprintf("%s, [ %s ]", result, strings.Join(children, ", "))
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
