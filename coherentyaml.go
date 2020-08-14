package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"log"
)

// plusieurs fichiers peuvent lu, ils deviennent un seul node avec un Coherent comme racine
// de base, coherentyaml, ne retourne rien.
// en cas d'erreur, le programme retourne la contradiction et sa position

func main() {
	flag.Parse()
	filename := flag.Args()[0]

	ymlContent, err := ioutil.ReadFile(filename)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Read file error : %s", err)
	}
	
	var ast Ast
	ast.Read(ymlContent)
	node := BigUglySwitch(ast.Interface())
	err = node.IsCoherent() 
	if nil != err { 
		log.Fatal(err)
	}
}
