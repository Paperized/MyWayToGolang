package main

import (
	"fmt"
	"json-parser/parser"
)

const jsonInput = `{
	"str": "string{}text",
	"true": true,
	"false": false,
	"null": null,
	"int": 1004,
	"float": 10.90,
	"obj": {
		"he": "ha"
	},
	"list": [
		10,
		20,
		30.5,
		{
			"isBool": true
		}
	]
}`

func main() {
	fmt.Println(parser.JsonStringToMap(jsonInput))
}
