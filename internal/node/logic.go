package node

import (
	"fmt"
	"reflect"
)

type logicalLiteral struct {
	node    Node
	negated bool
}

// logicalCoherence reports whether one value satisfies every supplied
// expression. OR branches are explored with the rest of the conjunction.
func logicalCoherence(nodes ...Node) error {
	queue := make([]logicalLiteral, 0, len(nodes))
	for _, node := range nodes {
		queue = append(queue, logicalLiteral{node: node})
	}
	if logicalSearch(queue, nil, nil) {
		return nil
	}
	return fmt.Errorf("logical expressions have no common value")
}

func logicalSearch(queue []logicalLiteral, positive, negative []Node) bool {
	if len(queue) == 0 {
		return atomicConstraintsSatisfiable(positive, negative)
	}
	literal := queue[0]
	rest := queue[1:]
	switch n := literal.node.(type) {
	case *Yes:
		if literal.negated {
			return false
		}
		return logicalSearch(rest, positive, negative)
	case *No:
		if literal.negated {
			return logicalSearch(rest, positive, negative)
		}
		return false
	case *Not:
		return logicalSearch(append([]logicalLiteral{{node: n.Child, negated: !literal.negated}}, rest...), positive, negative)
	case *Coherent:
		children := n.GetChild()
		if literal.negated {
			// NOT (A AND B) = NOT A OR NOT B.
			for _, child := range children {
				branch := append([]logicalLiteral{{node: child, negated: true}}, rest...)
				if logicalSearch(branch, positive, negative) {
					return true
				}
			}
			return false
		}
		return logicalSearch(prependLiterals(rest, children, false), positive, negative)
	case *OR:
		children := n.GetChild()
		if literal.negated {
			// NOT (A OR B) = NOT A AND NOT B.
			return logicalSearch(prependLiterals(rest, children, true), positive, negative)
		}
		for _, child := range children {
			branch := append([]logicalLiteral{{node: child}}, rest...)
			if logicalSearch(branch, positive, negative) {
				return true
			}
		}
		return false
	default:
		if literal.node == nil {
			return false
		}
		if literal.negated {
			return logicalSearch(rest, positive, append(append([]Node(nil), negative...), literal.node))
		}
		return logicalSearch(rest, append(append([]Node(nil), positive...), literal.node), negative)
	}
}

func prependLiterals(rest []logicalLiteral, children []Node, negated bool) []logicalLiteral {
	result := make([]logicalLiteral, 0, len(children)+len(rest))
	for _, child := range children {
		result = append(result, logicalLiteral{node: child, negated: negated})
	}
	return append(result, rest...)
}

func atomicConstraintsSatisfiable(positive, negative []Node) bool {
	for i, left := range positive {
		for _, right := range positive[i+1:] {
			if left.IsCoherentWith(right) != nil {
				return false
			}
		}
	}
	if len(positive) == 0 {
		// A finite list of exact exclusions cannot exhaust the shared universe.
		return true
	}
	for _, excluded := range negative {
		if excludesPositiveDomain(positive, excluded) {
			return false
		}
	}
	return true
}

func excludesPositiveDomain(positive []Node, excluded Node) bool {
	if exact, ok := excluded.(*Leaf); ok && !exact.isNeutral() {
		for _, candidate := range positive {
			leaf, ok := candidate.(*Leaf)
			if !ok || leaf.Value.Kind() != exact.Value.Kind() {
				continue
			}
			if leaf.isNeutral() {
				// The positive leaf denotes all values of its kind; excluding one
				// exact value leaves another value in these infinite domains.
				continue
			}
			if reflect.DeepEqual(leaf.Value.Interface(), exact.Value.Interface()) {
				return true
			}
		}
		return false
	}
	for _, candidate := range positive {
		if reflect.DeepEqual(candidate, excluded) {
			return true
		}
		leaf, ok := candidate.(*Leaf)
		neutral, isNeutralLeaf := excluded.(*Leaf)
		if ok && isNeutralLeaf && leaf.Value.Kind() == neutral.Value.Kind() && neutral.isNeutral() {
			return true
		}
	}
	return false
}
