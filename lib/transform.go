package main

func traverseChildren(parent ASTNode) {
	for i := 0; i < len(parent.children); i++ {
		child := parent.children[i]
		if child.(type) == string {
			continue
		}
		Transform(child)
	}
}

// node: RootNode | TemplateChildNode,
func Transform(node ASTNode) {
	switch node.tagType {
	case INTERPOLATION:
		break
	case ELEMENT:
	case ROOT:
		traverseChildren(node)
		break
	}
}

// func main() {
// 	template := `
// 		<template>
// 			<div id="count">
// 				<span>{{ add }}</span>
// 				<span>
// 					number:{{ count }}
// 				</span>
// 				<img src="asdsad.png" />
// 			</div>
// 		</template>
// 	`

// 	tokens := Parse(template)

// 	fmt.Println("tokens", tokens)
// }
