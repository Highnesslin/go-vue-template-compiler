package main

// func nodeToCode(node interface{}) (string, error) {
// 	switch node := node.(type) {
// 	case ASTNode:
// 		{
// 			children := ""

// 			for _, child := range node.children {
// 				cur, error := nodeToCode(child)
// 				if error != nil {
// 					return "", error
// 				}
// 				children += cur
// 			}

// 			props := ""
// 			for _, prop := range node.props {
// 				props += fmt.Sprintf("%s='%v'", prop.field, prop.value)
// 			}

// 			return fmt.Sprintf("h('%s', %s, %s)", node.tagName, props, children), nil
// 		}
// 	case TextNode:
// 		{
// 			return node.content, nil
// 		}
// 	}

// 	return "", errors.New("不识别的标签")
// }

// func Generate(node ASTNode) (string, error) {
// 	return nodeToCode(node)
// }
