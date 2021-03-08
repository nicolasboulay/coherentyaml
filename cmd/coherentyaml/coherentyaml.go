package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"log"
	"github.com/nicolasboulay/coherentyaml/cmd/node"
)

// plusieurs fichiers peuvent lu, ils deviennent un seul node avec un Coherent comme racine
// de base, coherentyaml, ne retourne rien.
// en cas d'erreur, le programme retourne la contradiction et sa position

func usage() {
		fmt.Print("Usage: coherentyaml fichier1 [fichier2 ...]\n");
		flag.PrintDefaults()
		os.Exit(0)
}

func main() {
	var help bool
	
	flag.BoolVar(&help, "h", false, "help")
	flag.Parse()

	if(help) {
		usage()
	}
	
	filename := flag.Arg(0)

	if (filename == "") {
		usage()
	}
	
	ymlContent, err := ioutil.ReadFile(filename)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Read file error : %s", err)
	}
	
	var ast Ast
	ast.Read(ymlContent)
	node1 := node.BigUglySwitch(ast.Interface())
	err = node1.IsCoherent() 
	if nil != err { 
		log.Fatal(err)
	}

	filename2 := flag.Arg(0)

	if (filename2 == "") {
		os.Exit(0)
	}
	
	ymlContent2, err := ioutil.ReadFile(filename)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Read file error : %s", err)
	}
	
	var ast2 Ast
	ast2.Read(ymlContent2)
	node2 := node.BigUglySwitch(ast2.Interface())
	err = node2.IsCoherentWith(node1) 
	if nil != err { 
		log.Fatal(err)
	}

	
}
