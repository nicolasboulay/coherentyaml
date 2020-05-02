package main

import (
	"github.com/goccy/go-yaml"
	"log"
)

type ast struct {
	V interface{}
}

func (ast *ast) Read(yml string) {
	 err := yaml.Unmarshal([]byte(yml), &ast.V);
	if err != nil {
		log.Fatal(err)
	}
}

