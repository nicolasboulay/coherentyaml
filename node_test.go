package main

import (
	"testing"
	"reflect"
)


func TestNode(t *testing.T) {
	s1 := Str("s1")
	c:= Coherent{&nArray{[]node{&s1,&StrZero}}}
	root := Coherent{&nArray{[]node{&s1, &StrZero, &c}}}
	
	err := root.IsCoherent()
	if (err != nil) {
		t.Errorf("Want coherency : %s", err)
	}
	s2 := Str("s2")
	or:= OR{&nArray{[]node{&s1,&s2}}}
	c = Coherent{&nArray{[]node{&or,&StrZero}}}
	err = c.IsCoherent()
	if (err != nil) {
		t.Errorf("Want coherency : %s", err)
	}
}


func TestCoherent(t *testing.T) {
	s1 := Str("s1")
	s2 := Str("s2")
	c:= Coherent{&nArray{[]node{&s1,&StrZero}}}
	root := Coherent{&nArray{[]node{&s1, &StrZero, &c}}}
	or := OR{&nArray{[]node{&s1, &s2}}}
	not := Not{&nArray{[]node{&s1}}}
	root2 := Coherent{&nArray{[]node{&not,&s2}}}
	root3 := Coherent{&nArray{[]node{&s1,&or}}}
	root4 := Coherent{&nArray{[]node{&not,&or}}} // (non A) && (A || B)
        neutralInt := leaf{reflect.ValueOf(-1)}
	coherentInt := Coherent{&nArray{[]node{&neutralInt,&leaf{reflect.ValueOf(0)}}}}
	tables := []struct{ name string; n node;}{
		{"root",&root},
		{"c",&c},
		{"or",&or},
		{"root2",&root2},
		{"root3",&root3},
		{"root4",&root4},
		{"intLiteral",&neutralInt},
		{"coherentInt",&coherentInt},
	}

	for _, node := range tables {
		err := node.n.IsCoherent()
		if (err != nil) {
			t.Errorf("Want coherency in %s : %s",node.name, err)
		}
	}
}

func TestNotCoherent(t *testing.T) {
	s1 := Str("s1")
	s2 := Str("s2")
	c:= Coherent{&nArray{[]node{&s2,&StrZero}}}
	root := Coherent{&nArray{[]node{&s1, &StrZero, &c}}}
	not := Not{&s1}
	root2 := Coherent{&nArray{[]node{&not,&s1}}}
	fakenot := Not{&nArray{[]node{&s1,&s2}}}
	//root3 := Coherent{[]node{&OR{[]node{&fakenot,&s1}},&s1}}
	intLiteral := leaf{reflect.ValueOf(3)}
	incoherentInt := Coherent{&nArray{[]node{&intLiteral,&leaf{reflect.ValueOf(2)}}}}
	tables := []struct{ name string; n node;}{
		{"root",&root},
		{"root2",&root2},
		{"fakenot",&fakenot},
		{"incoherentInt",&incoherentInt},
		//{"root3",&root3}, // using isCoherent in iscoherentwith is not decided

	}

	for _, node := range tables {
		//log.Print(" % " + node.name + "\n")
		err := node.n.IsCoherent()
		if (err == nil) {
			t.Errorf("Want error in %s : %v",node.name, node.n)
		}
	}
}
func TestIsNeutral(t *testing.T) {
	tables := []struct{ name string; l leaf; expected bool;}{
		{"true",leaf{reflect.ValueOf(true)}, false},
		{"false",leaf{reflect.ValueOf(false)}, false},
		{"1",leaf{reflect.ValueOf(uint(1))}, true},
		{"-1",leaf{reflect.ValueOf(-1)}, true}, 
		{"1.0",leaf{reflect.ValueOf(1.0)}, true},
		{"''",leaf{reflect.ValueOf("")}, true},
		{"2",leaf{reflect.ValueOf(uint(2))}, false},
		{"-2",leaf{reflect.ValueOf(-2)}, false}, 
		{"2.0",leaf{reflect.ValueOf(2.0)}, false},
		{"Plop",leaf{reflect.ValueOf("Plop")}, false},
	}

	for _, line := range tables {
		ok := line.l.isNeutral()
		if (ok != line.expected) {
			t.Errorf("%s: isNeutral (%v:%v) = %v should be %v",line.name,line.l.value.Kind(), line.l.String(), ok, line.expected)
		}
	}
}

