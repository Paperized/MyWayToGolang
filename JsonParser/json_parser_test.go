package main

import (
	"encoding/json"
	"json-parser/parser"
	"testing"
)

func BenchmarkCustomJson(*testing.B) {
	parser.JsonStringToMap(jsonInput, false)
}

func BenchmarkCustomJsonAlloc(*testing.B) {
	parser.JsonStringToMap(jsonInput, true)
}

func BenchmarkGoJson(b *testing.B) {
	var jsonInputBytes = []byte(jsonInput)
	var r map[string]any
	b.ResetTimer()

	json.Unmarshal(jsonInputBytes, &r)
}
