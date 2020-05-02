package main

import (
	"testing"
)

func TestNodeYaml(t *testing.T) {
	yml := `
%YAML 1.2
---
a: 
  aa:
    bb:
      - true
      - 0
      - 1.1
      - !!int 2
      - 3s
      - 2015-01-01
      - !!timestamp "2015-01-01"
      - ""
      - plip
      - !!str "plop"
      - !!string "plop"
      - ~
      - null
      - !expr "*.*"
b: c
`
	var ast ast
	ast.Read([]byte(yml))
//	yamlToNode(ast.V)
}

func TestRead(t *testing.T) {
		yml := `
%YAML 1.2
---
a: 1
b: c
`
	var a ast
	a.Read([]byte(yml))
	var v map[string]interface{}
	v = a.V.(map[string]interface{})
	va := v["a"].(uint64)
	if (va != 1) {
		t.Errorf("Unmarshal v.A = %v; want 1", va)
	}
	vc := v["b"].(string)
	if ( vc != "c") {
		t.Errorf("Unmarshal v.B = %s; want 'c'", vc)
	}
}
