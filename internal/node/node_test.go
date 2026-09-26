package node

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNode(t *testing.T) {
	s1 := MakeString("s1")
	c := &Coherent{&NArray{[]Node{s1, StrZero}}}
	root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}

	err := root.IsCoherent()
	if err != nil {
		t.Errorf("Want coherency %v : %s", root, err)
	}
	s2 := MakeString("s2")
	or := &OR{&NArray{[]Node{s1, s2}}}
	c = &Coherent{&NArray{[]Node{or, StrZero}}}
	err = c.IsCoherent()
	if err != nil {
		t.Errorf("Want coherency : %s", err)
	}
}

func TestCoherent(t *testing.T) {
	s1 := MakeString("s1")
	s2 := MakeString("s2")
	c := &Coherent{&NArray{[]Node{s1, StrZero}}}
	root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}
	or := &OR{&NArray{[]Node{s1, s2}}}
	not := &Not{s1}
	root2 := &Coherent{&NArray{[]Node{not, s2}}}
	root3 := &Coherent{&NArray{[]Node{s1, or}}}
	root4 := &Coherent{&NArray{[]Node{not, or}}} // (non A) && (A || B)
	neutralInt := &Leaf{reflect.ValueOf(-1)}
	coherentInt := &Coherent{&NArray{[]Node{neutralInt, &Leaf{reflect.ValueOf(0)}}}}
	tables := []struct {
		name string
		n    Node
	}{
		{"root", root},
		{"c", c},
		{"or", or},
		{"root2", root2},
		{"root3", root3},
		{"root4", root4},
		{"intLiteral", neutralInt},
		{"coherentInt", coherentInt},
	}

	for _, n := range tables {
		err := n.n.IsCoherent()
		if err != nil {
			t.Errorf("Want coherency in %s : %s\n%v\n", n.name, err, ToYAMLString(n.n))
		}
	}
}

func TestRoot2(t *testing.T) {
	s1 := MakeString("s1")
	s2 := MakeString("s2")
	//c:= &Coherent{&NArray{[]Node{s1,StrZero}}}
	//root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}
	//or := &OR{&NArray{[]Node{s1, s2}}}
	not := &Not{s1}
	root2 := &Coherent{&NArray{[]Node{not, s2}}}
	//root3 := &Coherent{&NArray{[]Node{s1,or}}}
	//root4 := &Coherent{&NArray{[]Node{not,or}}} // (non A) && (A || B)
	//neutralInt := &Leaf{reflect.ValueOf(-1)}
	//coherentInt := &Coherent{&NArray{[]Node{neutralInt,&Leaf{reflect.ValueOf(0)}}}}
	tables := []struct {
		name string
		n    Node
	}{
		{"root2", root2},
	}

	for _, n := range tables {
		err := n.n.IsCoherent()
		if err != nil {
			t.Errorf("Want coherency in %s : %s\n%v\n", n.name, err, ToYAMLString(n.n))
		}
	}
}

func TestNotCoherent(t *testing.T) {
	s1 := &Leaf{reflect.ValueOf("s1")}
	s2 := &Leaf{reflect.ValueOf("s2")}
	c := &Coherent{&NArray{[]Node{s2, StrZero}}}
	root := &Coherent{&NArray{[]Node{s1, StrZero, c}}}
	not := &Not{s1}
	root2 := &Coherent{&NArray{[]Node{not, s1}}}
	intLiteral := &Leaf{reflect.ValueOf(3)}
	incoherentInt := &Coherent{&NArray{[]Node{intLiteral, &Leaf{reflect.ValueOf(2)}}}}
	tables := []struct {
		name string
		n    Node
	}{
		{"root", root},
		{"root2", root2},
		{"incoherentInt", incoherentInt},
	}

	for _, node := range tables {
		err := node.n.IsCoherent()
		if err == nil {
			t.Errorf("Want error in %s : %v", node.name, ToYAMLString(node.n))
		}
	}
}

func TestNot(t *testing.T) {
	s1 := &Leaf{reflect.ValueOf("s1")}
	s2 := &Leaf{reflect.ValueOf("s2")}
	s3 := &Leaf{reflect.ValueOf("s3")}
	// wantCoherent is the expected logical result: true means coherent (nil error).
	// IsCoherent compares with yes. Since s1 is coherent with yes, Not(s1) is not.
	tests := []struct {
		name         string
		n            Node
		wantCoherent bool
	}{
		{"not true", &Not{yes}, false},
		{"not false", &Not{no}, true},
		{"not not true", &Not{&Not{yes}}, true},
		{"not not false", &Not{&Not{no}}, false},
		{"not not not true", &Not{&Not{&Not{yes}}}, false},
		{"not not not false", &Not{&Not{&Not{no}}}, true},
		{"not not not not true", &Not{&Not{&Not{&Not{yes}}}}, true},
		{"not not not not false", &Not{&Not{&Not{&Not{no}}}}, false},
		{"not s1", &Not{s1}, false},
		{"not not s1", &Not{&Not{s1}}, true},
		{"not not not s1", &Not{&Not{&Not{s1}}}, false},
		{"not (true or true)", &Not{&OR{&NArray{[]Node{yes, yes}}}}, false},
		{"not (true or false)", &Not{&OR{&NArray{[]Node{yes, no}}}}, false},
		{"not (false or true)", &Not{&OR{&NArray{[]Node{no, yes}}}}, false},
		{"not (false or false)", &Not{&OR{&NArray{[]Node{no, no}}}}, true},
		{"not (true and true)", &Not{&Coherent{&NArray{[]Node{yes, yes}}}}, false},
		{"not (true and false)", &Not{&Coherent{&NArray{[]Node{yes, no}}}}, true},
		{"not (false and true)", &Not{&Coherent{&NArray{[]Node{no, yes}}}}, true},
		{"not (false and false)", &Not{&Coherent{&NArray{[]Node{no, no}}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.IsCoherent()
			if (err == nil) != tt.wantCoherent {
				t.Errorf("IsCoherent() error = %v, wantCoherent = %v", err, tt.wantCoherent)
			}
		})
	}

	// Strict negation: coherent(Not(a), b) = !coherent(a, b).
	// Invert the entire comparison, even when b is incoherent or an OR.
	// Expectations are explicit; they are not computed using the implementation.
	t.Run("with", func(t *testing.T) {
		tests := []struct {
			name         string
			child        Node
			other        Node
			wantCoherent bool
		}{
			{"not true with true", yes, yes, false},
			{"not true with false", yes, no, true},
			{"not false with true", no, yes, true},
			{"not false with false", no, no, true},
			{"not not true with true", &Not{yes}, yes, true},
			{"not not true with false", &Not{yes}, no, false},
			{"not not false with true", &Not{no}, yes, false},
			{"not not false with false", &Not{no}, no, false},
			{"not s1 with s1", s1, s1, false},
			{"not s1 with s2", s1, s2, true},
			{"not s1 with yes", s1, yes, false},
			{"not s1 with no", s1, no, true},
			{"not not s1 with s1", &Not{s1}, s1, true},
			{"not not s1 with s2", &Not{s1}, s2, false},
			{"not not s1 with yes", &Not{s1}, yes, true},
			{"not not s1 with no", &Not{s1}, no, false},
			// NOT (match first OR match second): neither alternative may match.
			{"not s1 with or first match", s1, &OR{&NArray{[]Node{s1, s2}}}, false},
			{"not s1 with or second match", s1, &OR{&NArray{[]Node{s2, s1}}}, false},
			{"not s1 with or both match", s1, &OR{&NArray{[]Node{s1, s1}}}, false},
			{"not s1 with or neither match", s1, &OR{&NArray{[]Node{s2, s3}}}, true},
			{"not s1 with or both false", s1, &OR{&NArray{[]Node{no, no}}}, true},
			{"not not s1 with or one match", &Not{s1}, &OR{&NArray{[]Node{s1, s2}}}, true},
			{"not not s1 with or neither match", &Not{s1}, &OR{&NArray{[]Node{s2, s3}}}, false},
			// De Morgan also applies when the negated child is an OR.
			{"not or with first match", &OR{&NArray{[]Node{s1, s2}}}, s1, false},
			{"not or with second match", &OR{&NArray{[]Node{s2, s1}}}, s1, false},
			{"not or with both match", &OR{&NArray{[]Node{s1, s1}}}, s1, false},
			{"not or with neither match", &OR{&NArray{[]Node{s2, s3}}}, s1, true},
			{"not or with false", &OR{&NArray{[]Node{s1, s2}}}, no, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				n := &Not{tt.child}
				err := n.IsCoherentWith(tt.other)
				if (err == nil) != tt.wantCoherent {
					t.Errorf("IsCoherentWith() error = %v, wantCoherent = %v", err, tt.wantCoherent)
				}
			})
		}
	})
}
func TestIsNeutral(t *testing.T) {
	tables := []struct {
		name     string
		l        Leaf
		expected bool
	}{
		{"true", Leaf{reflect.ValueOf(true)}, false},
		{"false", Leaf{reflect.ValueOf(false)}, false},
		{"1", Leaf{reflect.ValueOf(uint(1))}, true},
		{"-1", Leaf{reflect.ValueOf(-1)}, true},
		{"1.0", Leaf{reflect.ValueOf(1.0)}, true},
		{"''", Leaf{reflect.ValueOf("")}, true},
		{"2", Leaf{reflect.ValueOf(uint(2))}, false},
		{"-2", Leaf{reflect.ValueOf(-2)}, false},
		{"2.0", Leaf{reflect.ValueOf(2.0)}, false},
		{"Plop", Leaf{reflect.ValueOf("Plop")}, false},
	}

	for _, line := range tables {
		ok := line.l.isNeutral()
		if ok != line.expected {
			t.Errorf("%s: isNeutral (%v:%v) = %v should be %v", line.name, line.l.Value.Kind(), line.l.String(), ok, line.expected)
		}
	}
}

func TestNStruct(t *testing.T) {
	m := NStruct{}
	//m := make(map[interface{}]struct{n Node; key Node})
	k := MakeString("1")

	//m[k.AsKey()] = struct{n Node; key Node}{k,k}
	m.set(k, k)
	k22 := MakeString("2")
	//m[k22.AsKey()] = struct{n Node; key Node}{k22,k22}
	m.set(k22, k22)
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
