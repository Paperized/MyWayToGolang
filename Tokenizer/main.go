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

	tokens, err := TokenizeJsonString(jsonInput)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Result: ", Format(tokens))

	obj, err := JsonStringToMap(jsonInput)
	if err != nil {
		fmt.Println(err)
	}

	objAsMap := obj.(map[string]any)
	fmt.Println("Result: ", obj)
	fmt.Println((objAsMap["list"].([]any)[3]))
}
