package node

import (
	"reflect"
	"testing"
)

func TestLogicalCoherenceGlobal(t *testing.T) {
	a, b, c := MakeString("a"), MakeString("b"), MakeString("c")
	or := func(nodes ...Node) Node { return &OR{&NArray{nodes}} }
	and := func(nodes ...Node) Node { return &Coherent{&NArray{nodes}} }
	tests := []struct {
		name  string
		nodes []Node
		want  bool
	}{
		{"three pairwise-overlapping unions", []Node{or(a, b), or(b, c), or(a, c)}, false},
		{"one shared witness", []Node{or(a, c), or(b, c), or(a, b, c)}, true},
		{"excluded alternatives", []Node{or(a, b), &Not{a}, &Not{b}}, false},
		{"one alternative remains", []Node{or(a, b, c), &Not{a}, &Not{b}}, true},
	}
	permutations := [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range permutations {
				nodes := []Node{tt.nodes[p[0]], tt.nodes[p[1]], tt.nodes[p[2]]}
				for _, expression := range []Node{and(nodes...), and(and(nodes[0], nodes[1]), nodes[2])} {
					if got := logicalCoherence(expression) == nil; got != tt.want {
						t.Errorf("logicalCoherence(%v) = %v, want %v", expression, got, tt.want)
					}
				}
			}
		})
	}
}

func TestLogicalCoherenceIdentities(t *testing.T) {
	a, b := MakeString("a"), MakeString("b")
	and := func(nodes ...Node) Node { return &Coherent{&NArray{nodes}} }
	or := func(nodes ...Node) Node { return &OR{&NArray{nodes}} }
	tests := []struct {
		name       string
		expression Node
		want       bool
	}{
		{"empty AND", and(), true}, {"empty OR", or(), false},
		{"AND with yes", and(a, yes), true}, {"AND with no", and(a, no), false},
		{"OR with no", or(a, no), true}, {"a is reflexive", and(a, a), true},
		{"distinct exacts contradict", and(a, b), false},
		{"not yes", &Not{yes}, false}, {"not no", &Not{no}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logicalCoherence(tt.expression) == nil; got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicalCoherenceLaws(t *testing.T) {
	a, b, c := MakeString("a"), MakeString("b"), MakeString("c")
	and := func(nodes ...Node) Node { return &Coherent{&NArray{nodes}} }
	or := func(nodes ...Node) Node { return &OR{&NArray{nodes}} }
	distributed := or(and(a, a), and(a, c))
	forms := []Node{and(or(a, b), or(a, c)), and(or(a, c), or(a, b)), and(and(or(a, b), or(a, b)), or(a, c)), distributed}
	for _, witness := range []struct {
		node Node
		want bool
	}{{a, true}, {b, false}, {c, false}} {
		for i, expression := range forms {
			t.Run(reflect.ValueOf(i).String(), func(t *testing.T) {
				got := logicalCoherence(expression, witness.node) == nil
				if got != witness.want {
					t.Errorf("%v with %v = %v, want %v", expression, witness.node, got, witness.want)
				}
			})
		}
	}
}

func TestLogicalCoherenceBooleanValues(t *testing.T) {
	trueValue := &Leaf{reflect.ValueOf(true)}
	falseValue := &Leaf{reflect.ValueOf(false)}
	text := MakeString("bonjour")
	and := func(nodes ...Node) Node { return &Coherent{&NArray{nodes}} }
	or := func(nodes ...Node) Node { return &OR{&NArray{nodes}} }
	tests := []struct {
		name       string
		expression Node
		want       bool
	}{
		{"distinct bools", and(trueValue, falseValue), false},
		{"not true with false", and(&Not{trueValue}, falseValue), true},
		{"not true with string", and(&Not{trueValue}, text), true},
		{"both bools excluded", and(&Not{trueValue}, &Not{falseValue}), true},
		{"bool domain exhausted", and(or(trueValue, falseValue), &Not{trueValue}, &Not{falseValue}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logicalCoherence(tt.expression) == nil; got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicalCoherenceDeMorgan(t *testing.T) {
	a, b, c := MakeString("a"), MakeString("b"), MakeString("c")
	and := func(nodes ...Node) Node { return &Coherent{&NArray{nodes}} }
	or := func(nodes ...Node) Node { return &OR{&NArray{nodes}} }
	forms := []struct {
		name       string
		expression Node
		want       []bool
	}{
		{"not OR", &Not{or(a, b)}, []bool{false, false, true}},
		{"AND of NOTs", and(&Not{a}, &Not{b}), []bool{false, false, true}},
		{"not AND", &Not{and(a, b)}, []bool{true, true, true}},
		{"OR of NOTs", or(&Not{a}, &Not{b}), []bool{true, true, true}},
	}
	for _, form := range forms {
		t.Run(form.name, func(t *testing.T) {
			for i, witness := range []Node{a, b, c} {
				if got := logicalCoherence(form.expression, witness) == nil; got != form.want[i] {
					t.Errorf("with %v got %v, want %v", witness, got, form.want[i])
				}
			}
		})
	}
}
