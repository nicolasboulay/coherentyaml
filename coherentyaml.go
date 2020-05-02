package main

import (
	"fmt"
//	"github.com/goccy/go-yaml"
	//goyaml3 "gopkg.in/yaml.v3"
//	"log"
	"flag"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()
	filename := flag.Args()[0]

	ymlContent, err := ioutil.ReadFile(filename)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Read file error : %s", err)
	}
	
	var ast Ast
	
	ast.Read(ymlContent)
	
	fmt.Println(ast)
}
