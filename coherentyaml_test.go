package main

import (
	"testing"
	//"log"
	//"reflect"
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
		node := BigUglySwitch(ast.Interface())
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
	node := BigUglySwitch(ast.Interface())
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
- a: [2 , 2] 			
  b: [2 , 2] 			
  c: [2 , 2] 			
			`,`
- c: [2 , 2] 			
  a: [2 , 2] 			
  c: [2 , 2] 						
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

	for _, yml := range tables {
		var ast1 Ast
		ast1.Read([]byte(yml.s1))
		var ast2 Ast
		ast2.Read([]byte(yml.s2))
		node1 := BigUglySwitch(ast1.Interface())
		node2 := BigUglySwitch(ast2.Interface())
		err := node1.IsCoherentWith(node2) 
		if nil != err { 
			t.Errorf("Want coherency : %s\n%#v\n%v", err, node1, node2)
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
		node1 := BigUglySwitch(ast1.Interface())
		node2 := BigUglySwitch(ast2.Interface())
		err := node1.IsCoherentWith(node2) 
		if nil == err { 
			t.Errorf("Want uncoherency :\n%v\n%v", node1, node2)
		}
	}
}

// from https://fr.wikipedia.org/wiki/Calcul_des_propositions
func TestCalculDeProposition(t *testing.T) {
	possible_set := []string{
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

	for _, A := range possible_set {
		var ast Ast
		ast.Read([]byte(A))
		nodeA := BigUglySwitch(ast.Interface())
		node := identity(nodeA)
		err := node.IsCoherent()
		if (err != nil) {
			t.Errorf("Want coherency in %s : %s", node, err)
		}
		
//		for _, B := range possible_set {
//
//			ast.Read([]byte(B))
//			nodeB := BigUglySwitch(ast.Interface())
//
//		
//			
//			err := node.n.IsCoherent()
//			if (err != nil) {
//				t.Errorf("Want coherency in %s : %s",node.name, err)
//			}
//		}
	}
}

func yor(a node, b node) node {
	return &OR{&nArray{[]node{a,b}}}
}
func yand(a node, b node) node {
	return &Coherent{&nArray{[]node{a,b}}}
}
func ynot(a node) node {
	return &Not{a}
}
// (~a & b) or (a & ~b)
func yxor(a node, b node) node {
	return yor(yand(ynot(a),b),yand(a,ynot(b)))
}
// (a -> b) & (b -> a)
func equivalence(a node, b node) node {
	return yand(implication(a,b),implication(b,a))
}

// (a -> b) = (~a or b)
func implication(a node, b node) node {
	return yor(ynot(a), b)
	
}

//
// theorème : must be always true :
//

// (A -> A)
func identity(a node) node {
	return implication(a, a)
}

// (A or ~A)
func excludedmiddle(a node) node {
	return yor(a,ynot(a))
}

// A -> ~~A
func doubleNegation(a node) node {
	return implication(a, ynot(ynot(a)))
}

// ~~A -> A
func classicalDoubleNegation(a node) node {
	return implication(ynot(ynot(a)),a)
}

// ((A -> B) -> A) -> A
func PeircesLaw (a node, b node) node {
	return implication(implication(implication(a,b),a),a)
}

// ~ (A & ~A)
func noncontradictionsLaw(a node) node {
	return ynot(yand(a,ynot(a)))
}

//~(A & B) <-> (~A or ~B)
func DeMorgansLaws1 (a node, b node) node {
	return equivalence(ynot(yand(a,b)), yor(ynot(a),ynot(b)))
}

//~(A or B) <-> (~A and ~B)
func DeMorgansLaws2 (a node, b node) node {
	return equivalence(ynot(yor(a,b)), yand(ynot(a),ynot(b)))
}

//(a -> b) -> (~b -> ~a)
func Contraposition(a node, b node) node {
	return implication(implication(a,b),implication(ynot(b),ynot(a)))
}

// ((A -> B) & A ) -> B
func ModusPonens(a node, b node) node {
	return implication(yand(implication(a,b),a),b)
}

// ((A -> B) & ~B) -> ~A
func ModusTollens(a node, b node) node {
	return implication(yand(implication(a,b),ynot(b)),ynot(b))
}

// ((A -> B) & (B -> C)) -> (A -> C)
func ModusBarbara(a node, b node, c node) node {
	return implication(yand(implication(a,b),implication(b,c)),implication(a,c))
}

// (a -> b) -> ((b -> c) -> (a -> c))
func ModusBarbaraImplicatif(a node, b node, c node) node {
	return implication(implication(a,b) , implication( implication(b,c) , implication(a,c)))
}

//(A & (B or C)) <-> ((A & B) or (A & C))
func DistributiveProperty1(a node, b node, c node) node {
	return equivalence(yand(a,yor(b,c)),yor(yand(a,b),yand(a,c)))
}

//(A or (B & C)) <-> ((A or B) & (A or C))
func DistributiveProperty2(a node, b node, c node) node {
	return equivalence(yor(a,yand(b,c)),yand(yor(a,b),yor(a,c)))
}

