package main

import "fmt"

func main() {
	template := `
		<template>
			<div id="app" className="flex">
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

	code, error := Generate(&ast)

	if error != nil {
		fmt.Println("error", error)
	}

	fmt.Println("code", code)
}
