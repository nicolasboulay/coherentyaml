package main

import (
	"testing"
	"fmt"
	"reflect"
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
		node1 := BigUglySwitch(ast1.Interface())
		node2 := BigUglySwitch(ast2.Interface())
		err := node1.IsCoherentWith(node2) 
		if nil != err {
			t.Errorf("Want coherency %v : %s\n%#v\n%#v\n", i, err, node1, node2)
			//fmt.Printf(" %s\n %s\n %v\n %v\n", yml.s1, yml.s2,ast1.Interface(), ast2.Interface())

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

func TestByList(t *testing.T) {
	tables := []struct {s string; expected bool; com string} {
		{"a: 2", true, "partial proposal"},
		{"Not: {a: 2}", true, "partial proposal"}, 
		{`Coherent:
- {a: 2}
- {a: 3}
`, false, "basic type check"},
{`Coherent:
- a: 2
- Not: 
    a: 2
`, false, "'Not' inside"},
{`Coherent:
- a: 2
- Not: 
    Coherent: 
      - a: 2
      - a: 1
`, false, "2 level checks"},
{`Coherent:
- a: 3
- Not: 
    Coherent: 
      - a: 2
      - a: 1
`, true, "2 levels checks"},
		{`Coherent:
- a: 2
- Not: 
    Coherent: 
      - a: 2
      - Not:
           a: 3
`, true, "2 levels check with Not"},
				{`Coherent:
- a: 3
- Not: 
    Coherent: 
      - a: 2
      - Not:
           a: 3
`, false, "2 levels check with not"},
						{`Coherent:
- a: 3
  b: 
     c: "plip"
     d: 2 
     e: "ploup"
- a: 1
  b:
     c: ""
     d: 2
     e: 
        OR: 
        - "plop"
        - "ploup"
`, true, "data + simple schema example"},
		{`Coherent:
- a: 3
  b: 
     c: "plip"
     d: 2 
     e: "plup"
- a: 1
  b:
     c: ""
     d: 2
     e: 
        OR: 
        - "plop"
        - "ploup"
`, false, "data + simple schema example"},
		{`Coherent:
- a: 3
- a: 3
- a: 1
`, true, "data + schema + a type in the data"},
								{`Coherent:
- a: "ploup"
- a: 
   Not: "ploup"
`, false, "not oddities"},
								{`Coherent:
- a: "plop"
- a: 
   Not: "ploup"
`, true, "not oddities"},
		{`Coherent:
- a: "ploup"
- a: 
    OR: 
     - Not: "ploup"
     - Not: "Plip"
`, false, "not/or oddities"},
		{`Coherent:
- a: "plop"
- a: 
    OR: 
     - Not: "ploup"
     - Not: "Plip"
`, true, "not/or oddities"},
				{`Coherent:
- a: "plop"
- a: 
    Not:
      OR: 
       - "ploup"
       - "Plip"
`, true, "not/or oddities"},
						{`Coherent:
- a: "ploup"
- a: 
    Not:
      OR: 
       - "ploup"
       - "Plip"
`, false, "not or oddities"},
		{`Coherent:
- a: "plop"
- a: 
    Coherent: 
     - Not: "ploup"
     - Not: "Plip"
`, true, "not/and oddities"},
		{`Coherent:
- a: "ploup"
- a: 
    Coherent: 
     - Not: "ploup"
     - Not: "Plip"
`, false, "not/and oddities"},
		{`Coherent:
- a: "ploup"
- a: 
    Not: 
     Or:
       - "ploup"
       - Not: "Plip"
`, false, "not/and oddities"},
				{`Coherent:
- a: "ploup"
- a: 
    Not: 
     Coherent:
       - "ploup"
       - Not: "Plip"
`, false, "not/and oddities"},		
		
	}

	for _, yml := range tables {
		var ast1 Ast
		ast1.Read([]byte(yml.s))
		node1 := BigUglySwitch(ast1.Interface())
		err := node1.IsCoherent()
		if (nil == err) != yml.expected { 
			t.Fatalf("%s :\nShould be %v (%v):\n%v\n", yml.com, yml.expected, err, toYAMLString(node1))
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
	theorem := []func(node) node{
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
		nodeA := BigUglySwitch(ast.Interface())
		
		for _,f :=  range theorem {
			node := f(nodeA)
			err := node.IsCoherent()
			//fmt.Printf("%v -> %v : %v\n", nodeA, node, err)
			if (err != nil) {
				t.Errorf("Want coherency in %s : %s", toYAMLString(node), err)
			}
		}
		break
	}	
}
func TestCalculDeProposition2(t *testing.T) {
	relation := []func(node, node) node{
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
		nodeA := BigUglySwitch(ast.Interface())
		
		for _, B := range possible_set {

			ast.Read([]byte(B))
			nodeB := BigUglySwitch(ast.Interface())

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
	relation := []func(node, node) node{
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
	relation3 := []func(node, node, node) node{
		ModusBarbara,
		ModusBarbaraImplicatif,
		DistributiveProperty1,
		DistributiveProperty2,
	}
	
	for _, A := range possible_set {
		var ast Ast
		ast.Read([]byte(A))
		nodeA := BigUglySwitch(ast.Interface())
		
		for _, B := range possible_set {

			ast.Read([]byte(B))
			nodeB := BigUglySwitch(ast.Interface())

			for i,f :=  range relation {
				node := f(nodeA, nodeB)
				err := node.IsCoherent()
				if (err != nil) {
					t.Errorf("Want coherency in %s %v %v : %s\n", relationString[i], nodeA, nodeB, err)
				}
			}
			break
			for _, C := range possible_set {

				ast.Read([]byte(C))
				nodeC := BigUglySwitch(ast.Interface())
				
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
	nodeA := BigUglySwitch(ast.Interface())
	ast.Read([]byte("a: 2"))
	nodeB := BigUglySwitch(ast.Interface())
	node := ModusTollens(nodeA, nodeB)
	err := node.IsCoherent()
	if (err != nil) {
		//fmt.Printf("modusTollens :\n %v\n", node)
		t.Errorf("Want coherency : %s\n %v\n", err, toYAMLString(node))
		//yamlString,_ := yaml.Marshal(node)
		//fmt.Printf("yaml :\n %v\n", yamlString)
	}
}

// Not IsCoherentWith  ~[{a:2 }]   ( ~[{a:2 }]  | {a:2 } ): false
func TestModusTollensPart(t *testing.T) {
	var ast Ast
	ast.Read([]byte("a: 2"))
	nodeA := BigUglySwitch(ast.Interface())

	node := yand(yor(ynot(nodeA),nodeA),ynot(nodeA))
	err := node.IsCoherent()
	if (err != nil) {
		//fmt.Printf("modusTollensPart :\n %v\n", node)
		t.Errorf("Want coherency : %s\n %v\n", err, node)
	}
}

func TestIncompleteProposal(t * testing.T) {
	yamlString := `

 a: 
  Not:
    b:
      Coherent:
      - 2
      - 3 
    c: 4
`
	var ast Ast
	ast.Read([]byte(yamlString))
	yamlNode := BigUglySwitch(ast.Interface())
	err := yamlNode.IsCoherent()
	if (err != nil) {
		t.Errorf("Want coherency : %v\n%s\n", err,toYAMLString(yamlNode))
	}
	//fmt.Printf("incompleteProposal :\n%s\n", toYAMLString(yamlNode))
}

func TestIncompleteProposal2(t * testing.T) {
	yamlString := `
 a: 
  Not:
    b:
      Coherent:
      - 3
      - 3 
    cccc: 4
`
	var ast Ast
	ast.Read([]byte(yamlString))
	yamlNode := BigUglySwitch(ast.Interface())
	err := yamlNode.IsCoherent()
	if (err == nil) {
		t.Errorf("Want INcoherency : \n%s\n", toYAMLString(yamlNode))
	}
	//fmt.Printf("incompleteProposal :\n%s\n", toYAMLString(yamlNode))
}

func TestStructure(t * testing.T) {
	yamlString := `
 Coherent: 
  - Not:
     b: 3
     c: 4
  - b: 3 
    c: 4
`
	var ast Ast
	ast.Read([]byte(yamlString))
	yamlNode := BigUglySwitch(ast.Interface())
	err := yamlNode.IsCoherent()
	if (err == nil) {
		t.Errorf("Want INcoherency : \n%s\n", toYAMLString(yamlNode))
	}
	fmt.Printf("incompleteProposal :\n%s\n", toYAMLString(yamlNode))
}

func TestStructure2(t * testing.T) {
	yamlString := `
 Coherent: 
  - Not:
     b: 3
  - b: 3 
`
	var ast Ast
	ast.Read([]byte(yamlString))
	yamlNode := BigUglySwitch(ast.Interface())
	err := yamlNode.IsCoherent()
	if (err == nil) {
		t.Errorf("Want INcoherency : \n%s\n", toYAMLString(yamlNode))
	}
	fmt.Printf("incompleteProposal :\n%s\n", toYAMLString(yamlNode))
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

func notTrue(a node) node {
	return ynot(&leaf{reflect.ValueOf(1)})
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
	return implication(yand(implication(a,b),ynot(b)),ynot(a))
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
