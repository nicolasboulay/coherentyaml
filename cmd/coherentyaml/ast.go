package main

import (
	"github.com/goccy/go-yaml"
	"log"
)

type Ast struct {
	V interface{}
}

func (ast *Ast) Read(yml []byte) {
	 err := yaml.Unmarshal([]byte(yml), &ast.V);
	if err != nil {
		log.Fatal(err)
	}
}

func (ast *Ast) Interface() interface{} {
	return ast.V
}
