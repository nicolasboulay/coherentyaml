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

func TestCoherentGlobalConstraints(t *testing.T) {
	a, b, c := MakeString("a"), MakeString("b"), MakeString("c")
	or := func(nodes ...Node) Node { return &OR{&NArray{nodes}} }
	and := func(nodes ...Node) Node { return &Coherent{&NArray{nodes}} }
	tests := []struct {
		name         string
		nodes        [3]Node
		wantCoherent bool
	}{
		// Each pair has a witness, but the complete conjunction has none.
		{"both alternatives excluded", [3]Node{or(a, b), &Not{a}, &Not{b}}, false},
		// No Not involved: pairwise intersections are b, a, and c respectively.
		{"three overlapping disjunctions", [3]Node{or(a, b), or(b, c), or(a, c)}, false},
		// Positive controls: c remains a witness for the whole conjunction.
		{"one alternative remains", [3]Node{or(a, b, c), &Not{a}, &Not{b}}, true},
		{"common alternative", [3]Node{or(a, c), or(b, c), or(a, b, c)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// All six orders must agree, including when constraints are grouped
			// or supplied through IsCoherentWith instead of a flat AND.
			for _, order := range []struct {
				name    string
				i, j, k int
			}{{"012", 0, 1, 2}, {"021", 0, 2, 1}, {"102", 1, 0, 2}, {"120", 1, 2, 0}, {"201", 2, 0, 1}, {"210", 2, 1, 0}} {
				t.Run(order.name, func(t *testing.T) {
					x, y, z := tt.nodes[order.i], tt.nodes[order.j], tt.nodes[order.k]
					checks := []struct {
						name string
						run  func() error
					}{
						{"pair", func() error { return x.IsCoherentWith(y) }},
						{"flat", and(x, y, z).IsCoherent},
						{"nested", and(and(x, y), z).IsCoherent},
						{"with", func() error { return and(x, y).IsCoherentWith(z) }},
						{"with reversed", func() error { return z.IsCoherentWith(and(x, y)) }},
					}
					for _, check := range checks {
						t.Run(check.name, func(t *testing.T) {
							want := tt.wantCoherent || check.name == "pair"
							if err := check.run(); (err == nil) != want {
								t.Errorf("error = %v, wantCoherent = %v", err, want)
							}
						})
					}
				})
			}
		})
	}
}

func TestNot(t *testing.T) {
	s1 := &Leaf{reflect.ValueOf("s1")}
	s2 := &Leaf{reflect.ValueOf("s2")}
	s3 := &Leaf{reflect.ValueOf("s3")}
	// Coherent is logical AND; coherence means that a satisfying value exists.
	// Not complements the allowed values, not the boolean coherence result.
	// wantCoherent=true means satisfiable (nil error).
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
		{"not s1", &Not{s1}, true},
		{"not not s1", &Not{&Not{s1}}, true},
		{"not not not s1", &Not{&Not{&Not{s1}}}, true},
		{"not (true or true)", &Not{&OR{&NArray{[]Node{yes, yes}}}}, false},
		{"not (true or false)", &Not{&OR{&NArray{[]Node{yes, no}}}}, false},
		{"not (false or true)", &Not{&OR{&NArray{[]Node{no, yes}}}}, false},
		{"not (false or false)", &Not{&OR{&NArray{[]Node{no, no}}}}, true},
		{"not (true and true)", &Not{&Coherent{&NArray{[]Node{yes, yes}}}}, false},
		{"not (true and false)", &Not{&Coherent{&NArray{[]Node{yes, no}}}}, true},
		{"not (false and true)", &Not{&Coherent{&NArray{[]Node{no, yes}}}}, true},
		{"not (false and false)", &Not{&Coherent{&NArray{[]Node{no, no}}}}, true},
		{"not (s1 or s2)", &Not{&OR{&NArray{[]Node{s1, s2}}}}, true},
		{"not (s1 and s2)", &Not{&Coherent{&NArray{[]Node{s1, s2}}}}, true},
		{"not (s1 or not s1)", &Not{&OR{&NArray{[]Node{s1, &Not{s1}}}}}, false},
		{"not (s1 and not s1)", &Not{&Coherent{&NArray{[]Node{s1, &Not{s1}}}}}, true},
		// Every pair is satisfiable, but no single value satisfies all three.
		{"or with both alternatives excluded", &Coherent{&NArray{[]Node{
			&OR{&NArray{[]Node{s1, s2}}}, &Not{s1}, &Not{s2},
		}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.IsCoherent()
			if (err == nil) != tt.wantCoherent {
				t.Errorf("IsCoherent() error = %v, wantCoherent = %v", err, tt.wantCoherent)
			}
		})
	}

	// Not(a) is coherent with b when (NOT a) AND b has a solution.
	// Both argument orders and the Coherent operator must agree.
	// Expectations are explicit; they are not computed using the implementation.
	t.Run("with", func(t *testing.T) {
		tests := []struct {
			name         string
			child        Node
			other        Node
			wantCoherent bool
		}{
			{"not true with true", yes, yes, false},
			{"not true with false", yes, no, false},
			{"not false with true", no, yes, true},
			{"not false with false", no, no, false},
			{"not not true with true", &Not{yes}, yes, true},
			{"not not true with false", &Not{yes}, no, false},
			{"not not false with true", &Not{no}, yes, false},
			{"not not false with false", &Not{no}, no, false},
			{"not s1 with s1", s1, s1, false},
			{"not s1 with s2", s1, s2, true},
			{"not s1 with yes", s1, yes, true},
			{"not s1 with no", s1, no, false},
			{"not not s1 with s1", &Not{s1}, s1, true},
			{"not not s1 with s2", &Not{s1}, s2, false},
			{"not not s1 with yes", &Not{s1}, yes, true},
			{"not not s1 with no", &Not{s1}, no, false},
			// One OR alternative outside s1 is sufficient; s2 is a witness.
			{"not s1 with or first match", s1, &OR{&NArray{[]Node{s1, s2}}}, true},
			{"not s1 with or second match", s1, &OR{&NArray{[]Node{s2, s1}}}, true},
			{"not s1 with or both match", s1, &OR{&NArray{[]Node{s1, s1}}}, false},
			{"not s1 with or neither match", s1, &OR{&NArray{[]Node{s2, s3}}}, true},
			{"not s1 with or both false", s1, &OR{&NArray{[]Node{no, no}}}, false},
			{"not not s1 with or one match", &Not{s1}, &OR{&NArray{[]Node{s1, s2}}}, true},
			{"not not s1 with or neither match", &Not{s1}, &OR{&NArray{[]Node{s2, s3}}}, false},
			// De Morgan also applies when the negated child is an OR.
			{"not or with first match", &OR{&NArray{[]Node{s1, s2}}}, s1, false},
			{"not or with second match", &OR{&NArray{[]Node{s2, s1}}}, s1, false},
			{"not or with both match", &OR{&NArray{[]Node{s1, s1}}}, s1, false},
			{"not or with neither match", &OR{&NArray{[]Node{s2, s3}}}, s1, true},
			{"not or with false", &OR{&NArray{[]Node{s1, s2}}}, no, false},
			{"not or with partially overlapping or", &OR{&NArray{[]Node{s1, s2}}}, &OR{&NArray{[]Node{s2, s3}}}, true},
			{"not or with identical or", &OR{&NArray{[]Node{s1, s2}}}, &OR{&NArray{[]Node{s1, s2}}}, false},
			{"not and with s1", &Coherent{&NArray{[]Node{s1, s2}}}, s1, true},
			{"not repeated and with s1", &Coherent{&NArray{[]Node{s1, s1}}}, s1, false},
			{"not repeated and with s2", &Coherent{&NArray{[]Node{s1, s1}}}, s2, true},
			{"not s1 with not s1", s1, &Not{s1}, true},
			{"not s1 with not s2", s1, &Not{s2}, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				n := &Not{tt.child}
				checks := []struct {
					name string
					run  func() error
				}{
					{"forward", func() error { return n.IsCoherentWith(tt.other) }},
					{"reverse", func() error { return tt.other.IsCoherentWith(n) }},
					{"and", (&Coherent{&NArray{[]Node{n, tt.other}}}).IsCoherent},
					{"and reversed", (&Coherent{&NArray{[]Node{tt.other, n}}}).IsCoherent},
				}
				for _, check := range checks {
					t.Run(check.name, func(t *testing.T) {
						err := check.run()
						if (err == nil) != tt.wantCoherent {
							t.Errorf("error = %v, wantCoherent = %v", err, tt.wantCoherent)
						}
					})
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
