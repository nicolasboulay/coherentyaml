package main

import (
	"fmt"
	"github.com/goccy/go-yaml"
	//goyaml3 "gopkg.in/yaml.v3"
	"log"
)

func main() {
	yml := `---
foo: 1
bar: c
A: 2
B: d
`
	var v struct {
		A int    `yaml:"foo" json:"A"`
		B string `yaml:"bar" json:"B"`
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		log.Fatal(err)
	}
	fmt.Println(v.A)
	fmt.Println(v.B)

   fmt.Println("Hello, World!")
}
