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

func main() {
	var help bool
	var verbose bool
	
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	if(help) {
		usage()
	}

	SetVerbose(verbose)
	filename := flag.Arg(0)
	if (filename == "") {
		usage()
	}
	
	node1 := makeNodeFromFile(filename)
	err := node1.IsCoherent() 
	if nil != err {
		fatalError(err,node1)
		//log.Fatal(err)
	}

	filename2 := flag.Arg(1)
	if (filename2 == "") {
		os.Exit(0)
	}
	
	node2 := makeNodeFromFile(filename2)

	err = node2.IsCoherentWith(node1) 
	if nil != err { 
		fmt.Print(node.ToYAMLString(node1) + "vs\n")
		fatalError(err,node2)
	}

	
}

func fatalError(err error, n node.Node) {
	fmt.Print(node.ToYAMLString(n))
	log.Fatal(err)
}

func makeNode(s string) node.Node {
	var ast Ast
	ast.Read([]byte(s))
	VerbosePrintfIn("Making ast : \n%v\n", ast.Interface())
	n:= node.BigUglySwitch(ast.Interface())
	VerbosePrintfIn("Making node : \n%v", node.ToYAMLString(n))
	return n
}

func makeNodeFromFile(filename string) node.Node {
	ymlContent, err := ioutil.ReadFile(filename)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Read file error : %s", err)
	}
	VerbosePrintfIn("Parsing file : %v\n%s", filename, string(ymlContent))
	return makeNode(string(ymlContent))
}

func usage() {
		fmt.Print("Usage: coherentyaml fichier1 [fichier2 ...]\n");
		flag.PrintDefaults()
		os.Exit(0)
}

