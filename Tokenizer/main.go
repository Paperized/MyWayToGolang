package main

import (
	"fmt"
)

func main() {
	jsonInput := `{
		"str": "string{}text",
		"true": true,
		"false": false,
		"null": null,
		"int": 1004,
		"float": 10.90
	}`

	tokens, err := TokenizeJsonString(jsonInput)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Result: ", FormatTrim(tokens))
}
