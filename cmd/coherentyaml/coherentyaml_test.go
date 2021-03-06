package main

import (
	"testing"
	"fmt"
	"reflect"
	"github.com/nicolasboulay/coherentyaml/cmd/node"
)
func TestIso(t *testing.T) {
	tables := []string {
		"a: 1",
		`
a:
  a:
     1
`,`
a:
  - "toto"
  - "titi"
`,`
a:
  - ""
  - "titi"
`,`
- a: 2
- b: 3
`,`
- a: 1
  b: 2
`,
	}

	for _, yml := range tables {
		var ast Ast
		ast.Read([]byte(yml))
		node := node.BigUglySwitch(ast.Interface())
		err := node.IsCoherentWith(node) 
		if nil != err { 
			t.Errorf("Want coherency : %s\n%#v\n%v", err, node, node)
		}
	}
}
func TestIsoAll(t *testing.T) {
	
	yml := `
%YAML 1.2
---
a: 
  aa:
    bb:
      - true
      - 0
      - 1.1
      - 3s
      - 2015-01-01
      - !!timestamp "2015-01-01"
      - ""
      - plip
      - nil
b: c
`
	var ast Ast
	ast.Read([]byte(yml))
	node := node.BigUglySwitch(ast.Interface())
	err := node.IsCoherentWith(node) 
	if nil != err { 
		t.Errorf("Want coherency : %s\n%#v\n%v", err, node, node)
	}
}

func TestShallMatch(t *testing.T) {
	tables := []struct {s1 string; s2 string;} {
		{"a: 2", "a: 1"},
		{"a: 2", "a: \n  Not: 3"},
		{"{a: 2, b: 2}", "{a: 1, b: 2, c: 3}"},
		{`
- a: 2
  b: 3
  c: 4
- a:
   - 2
   - 2
  b: [2 , 2]
  c: [2 , 2]
`,`
- c: [2 , 2]
  a: [2 , 2]
  b: [2 , 2]
- a: 2
  b: 3
  c: 4
`},
		{`
a: 
 a:
  a:
   -
    - 
     - 2 
`,`
a: 
 a:
  a:
   -
    - 
     - 2 
`},
	}

	for i, yml := range tables {
		var ast1 Ast
		ast1.Read([]byte(yml.s1))
		var ast2 Ast
		ast2.Read([]byte(yml.s2))
		node1 := node.BigUglySwitch(ast1.Interface())
		node2 := node.BigUglySwitch(ast2.Interface())
		err := node1.IsCoherentWith(node2) 
		if nil != err {
			t.Errorf("Want coherency %v : %s\n%#v\n%#v\n", i, err, node1, node2)
			fmt.Printf(" %s\n %s\n %v\n %v\n", yml.s1, yml.s2,ast1.Interface(), ast2.Interface())

		}
	}
}
func TestShallNotMatch(t *testing.T) {
	tables := []struct {s1 string; s2 string;} {
		{"a: 2", "a: 3"},
		{"a: 2", "a: \"toto\""},
		{"a: 2", "a: 2.0"},
		{"- 2", "- toto"},
	}

	for _, yml := range tables {
		var ast1 Ast
		ast1.Read([]byte(yml.s1))
		var ast2 Ast
		ast2.Read([]byte(yml.s2))
		node1 := node.BigUglySwitch(ast1.Interface())
		node2 := node.BigUglySwitch(ast2.Interface())
		err := node1.IsCoherentWith(node2) 
		if nil == err { 
			t.Errorf("Want uncoherency :\n%v\n%v", node1, node2)
		}
	}
}

var possible_set []string = []string{
		"a: 2",
		"a: toto",
//s s s 
`	
a:
  a: 
   a: 2
`,
// s s a
`
a:
  a: 
  - 2
  - "toto"  
`,
// s a s
`
a: 
  - a: 2
    b: "toto"
  - c: 2
    d: "toto"  
`,
// a s a
` 
- a: 
  - 2
  - 3
  b: 
  - 4
  - 5
- c: 
  - 6
  - 7
  d: 
  - 8
  - 9
`,
// a a s 
`
- 
  - a: 2
    b: 3
  - c: 4
    d: 5
- 
  - e: 2
    f: 3
  - g: 4
    h: 5
`,
}

// from https://fr.wikipedia.org/wiki/Calcul_des_propositions

func TestCalculDePropositionTheorem(t *testing.T) {
	theorem := []func(node.Node) node.Node{
		//notTrue,
		identity,
		excludedmiddle,
		doubleNegation,
		classicalDoubleNegation,
		noncontradictionsLaw,
	}
	for _, A := range possible_set {
		var ast Ast
		ast.Read([]byte(A))
		nodeA := node.BigUglySwitch(ast.Interface())
		
		for _,f :=  range theorem {
			node := f(nodeA)
			err := node.IsCoherent()
			//fmt.Printf("%v -> %v : %v\n", nodeA, node, err)
			if (err != nil) {
				t.Errorf("Want coherency in %s : %s", node, err)
			}
		}
		break
	}	
}
func TestCalculDeProposition2(t *testing.T) {
	relation := []func(node.Node, node.Node) node.Node{
		PeircesLaw,
		DeMorgansLaws1,
		DeMorgansLaws2,
		Contraposition,
		ModusPonens,
		ModusTollens,
	}
	for _, A := range possible_set {
		var ast Ast
		ast.Read([]byte(A))
		nodeA := node.BigUglySwitch(ast.Interface())
		
		for _, B := range possible_set {

			ast.Read([]byte(B))
			nodeB := node.BigUglySwitch(ast.Interface())

			for i,f :=  range relation {
				node := f(nodeA, nodeB)
				err := node.IsCoherent()
				if (err != nil) {
					t.Fatalf("Want coherency in %v %v %v : %s",
						i, nodeA, nodeB, err)
				}
			}
		}
	}
}

func TestCalculDeProposition(t *testing.T) {
	relation := []func(node.Node, node.Node) node.Node{
		PeircesLaw,
		DeMorgansLaws1,
		DeMorgansLaws2,
		Contraposition,
		ModusPonens,
		ModusTollens,
	}
	relationString := []string{
		"PeircesLaw",
		"DeMorgansLaws1",
		"DeMorgansLaws2",
		"Contraposition",
		"ModusPonens",
		"ModusTollens",
	}
	relation3 := []func(node.Node, node.Node, node.Node) node.Node{
		ModusBarbara,
		ModusBarbaraImplicatif,
		DistributiveProperty1,
		DistributiveProperty2,
	}
	
	for _, A := range possible_set {
		var ast Ast
		ast.Read([]byte(A))
		nodeA := node.BigUglySwitch(ast.Interface())
		
		for i, B := range possible_set {

			ast.Read([]byte(B))
			nodeB := node.BigUglySwitch(ast.Interface())

			for _,f :=  range relation {
				node := f(nodeA, nodeB)
				err := node.IsCoherent()
				if (err != nil) {
					t.Errorf("Want coherency in %s %v %v : %s", relationString[i], nodeA, nodeB, err)
				}
			}
			break
			for _, C := range possible_set {

				ast.Read([]byte(C))
				nodeC := node.BigUglySwitch(ast.Interface())
				
				for _,f :=  range relation3 {
					node := f(nodeA, nodeB, nodeC)
					err := node.IsCoherent()
					if (err != nil) {
						t.Errorf("Want coherency in %s : %s", node, err)
					}
				}
			}
		}
	}
}

func TestModusTollens(t *testing.T) {
	var ast Ast
	ast.Read([]byte("a: 2"))
	nodeA := node.BigUglySwitch(ast.Interface())
	ast.Read([]byte("a: 2"))
	nodeB := node.BigUglySwitch(ast.Interface())
	node := ModusTollens(nodeA, nodeB)
	err := node.IsCoherent()
	if (err != nil) {
		fmt.Printf("modusTollens :\n %v\n", node)
		t.Errorf("Want coherency : %s\n", err)
	}
}

func yor(a node.Node, b node.Node) node.Node {
	return &node.OR{&node.NArray{[]node.Node{a,b}}}
}
func yand(a node.Node, b node.Node) node.Node {
	return &node.Coherent{&node.NArray{[]node.Node{a,b}}}
}
func ynot(a node.Node) node.Node {
	return &node.Not{a}
}
// (~a & b) or (a & ~b)
func yxor(a node.Node, b node.Node) node.Node {
	return yor(yand(ynot(a),b),yand(a,ynot(b)))
}
// (a -> b) & (b -> a)
func equivalence(a node.Node, b node.Node) node.Node {
	return yand(implication(a,b),implication(b,a))
}

// (a -> b) = (~a or b)
func implication(a node.Node, b node.Node) node.Node {
	return yor(ynot(a), b)
	
}

//
// theorème : must be always true :
//

// (A -> A)
func identity(a node.Node) node.Node {
	return implication(a, a)
}

func notTrue(a node.Node) node.Node {
	return ynot(&node.Leaf{reflect.ValueOf(1)})
}

// (A or ~A)
func excludedmiddle(a node.Node) node.Node {
	return yor(a,ynot(a))
}

// A -> ~~A
func doubleNegation(a node.Node) node.Node {
	return implication(a, ynot(ynot(a)))
}

// ~~A -> A
func classicalDoubleNegation(a node.Node) node.Node {
	return implication(ynot(ynot(a)),a)
}

// ((A -> B) -> A) -> A
func PeircesLaw (a node.Node, b node.Node) node.Node {
	return implication(implication(implication(a,b),a),a)
}

// ~ (A & ~A)
func noncontradictionsLaw(a node.Node) node.Node {
	return ynot(yand(a,ynot(a)))
}

//~(A & B) <-> (~A or ~B)
func DeMorgansLaws1 (a node.Node, b node.Node) node.Node {
	return equivalence(ynot(yand(a,b)), yor(ynot(a),ynot(b)))
}

//~(A or B) <-> (~A and ~B)
func DeMorgansLaws2 (a node.Node, b node.Node) node.Node {
	return equivalence(ynot(yor(a,b)), yand(ynot(a),ynot(b)))
}

//(a -> b) -> (~b -> ~a)
func Contraposition(a node.Node, b node.Node) node.Node {
	return implication(implication(a,b),implication(ynot(b),ynot(a)))
}

// ((A -> B) & A ) -> B
func ModusPonens(a node.Node, b node.Node) node.Node {
	return implication(yand(implication(a,b),a),b)
}

// ((A -> B) & ~B) -> ~A
func ModusTollens(a node.Node, b node.Node) node.Node {
	return implication(yand(implication(a,b),ynot(b)),ynot(a))
}

// ((A -> B) & (B -> C)) -> (A -> C)
func ModusBarbara(a node.Node, b node.Node, c node.Node) node.Node {
	return implication(yand(implication(a,b),implication(b,c)),implication(a,c))
}

// (a -> b) -> ((b -> c) -> (a -> c))
func ModusBarbaraImplicatif(a node.Node, b node.Node, c node.Node) node.Node {
	return implication(implication(a,b) , implication( implication(b,c) , implication(a,c)))
}

//(A & (B or C)) <-> ((A & B) or (A & C))
func DistributiveProperty1(a node.Node, b node.Node, c node.Node) node.Node {
	return equivalence(yand(a,yor(b,c)),yor(yand(a,b),yand(a,c)))
}

//(A or (B & C)) <-> ((A or B) & (A or C))
func DistributiveProperty2(a node.Node, b node.Node, c node.Node) node.Node {
	return equivalence(yor(a,yand(b,c)),yand(yor(a,b),yor(a,c)))
}
