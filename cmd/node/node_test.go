package node

import (
	"testing"
	"reflect"
	"encoding/json"
)

func TestNode(t *testing.T) {
	s1 := MakeString("s1")
	c:= &Coherent{&NArray{[]Node{s1,StrZero}}}
	root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}
	
	err := root.IsCoherent()
	if (err != nil) {
		t.Errorf("Want coherency %v : %s", root, err)
	}
	s2 := MakeString("s2")
	or:= &OR{&NArray{[]Node{s1,s2}}}
	c = &Coherent{&NArray{[]Node{or,StrZero}}}
	err = c.IsCoherent()
	if (err != nil) {
		t.Errorf("Want coherency : %s", err)
	}
}

func TestCoherent(t *testing.T) {
	s1 := MakeString("s1")
	s2 := MakeString("s2")
	c:= &Coherent{&NArray{[]Node{s1,StrZero}}}
	root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}
	or := &OR{&NArray{[]Node{s1, s2}}}
	not := &Not{s1}
	root2 := &Coherent{&NArray{[]Node{not,s2}}}
	root3 := &Coherent{&NArray{[]Node{s1,or}}}
	root4 := &Coherent{&NArray{[]Node{not,or}}} // (non A) && (A || B)
        neutralInt := &Leaf{reflect.ValueOf(-1)}
	coherentInt := &Coherent{&NArray{[]Node{neutralInt,&Leaf{reflect.ValueOf(0)}}}}
	tables := []struct{ name string; n Node;}{
		{"root",root},
		{"c",c},
		{"or",or},
		{"root2",root2},
		{"root3",root3},
		{"root4",root4},
		{"intLiteral",neutralInt},
		{"coherentInt",coherentInt},
	}

	for _, node := range tables {
		err := node.n.IsCoherent()
		if (err != nil) {
			t.Errorf("Want coherency in %s : %s",node.name, err)
		}
	}
}

func TestNotCoherent(t *testing.T) {
	s1 := &Leaf{reflect.ValueOf("s1")}
	s2 := &Leaf{reflect.ValueOf("s2")}
	c:= &Coherent{&NArray{[]Node{s2,StrZero}}}
	root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}
	not := &Not{s1}
	root2 := &Coherent{&NArray{[]Node{not,s1}}}
	intLiteral := &Leaf{reflect.ValueOf(3)}
	incoherentInt := &Coherent{&NArray{[]Node{intLiteral,&Leaf{reflect.ValueOf(2)}}}}
	tables := []struct{ name string; n Node;}{
		{"not s1",not},
		{"root",root},
		{"root2",root2},
		{"incoherentInt",incoherentInt},
	}

	for _, node := range tables {
		err := node.n.IsCoherent()
		if (err == nil) {
			t.Errorf("Want error in %s : %v",node.name, node.n)
		}
	}
}
func TestIsNeutral(t *testing.T) {
	tables := []struct{ name string; l Leaf; expected bool;}{
		{"true",Leaf{reflect.ValueOf(true)}, false},
		{"false",Leaf{reflect.ValueOf(false)}, false},
		{"1",Leaf{reflect.ValueOf(uint(1))}, true},
		{"-1",Leaf{reflect.ValueOf(-1)}, true}, 
		{"1.0",Leaf{reflect.ValueOf(1.0)}, true},
		{"''",Leaf{reflect.ValueOf("")}, true},
		{"2",Leaf{reflect.ValueOf(uint(2))}, false},
		{"-2",Leaf{reflect.ValueOf(-2)}, false}, 
		{"2.0",Leaf{reflect.ValueOf(2.0)}, false},
		{"Plop",Leaf{reflect.ValueOf("Plop")}, false},
	}

	for _, line := range tables {
		ok := line.l.isNeutral()
		if (ok != line.expected) {
			t.Errorf("%s: isNeutral (%v:%v) = %v should be %v",line.name,line.l.Value.Kind(), line.l.String(), ok, line.expected)
		}
	}
}

func TestNStruct(t *testing.T) {
	m := NStruct{}
	//m := make(map[interface{}]struct{n Node; key Node})
	k := MakeString("1")
	
	//m[k.AsKey()] = struct{n Node; key Node}{k,k}
	m.set(k,k)
	k22 := MakeString("2")
	//m[k22.AsKey()] = struct{n Node; key Node}{k22,k22}
	m.set(k22,k22)
	//n := nStruct{m}

	k1 := MakeString("1")
	if m.get(k1).AsKey() != k1.AsKey() {
		t.Errorf("get() error %#v %v %#v\n", k1, m, m.get(k1))
	}
}

func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}
