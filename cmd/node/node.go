package node

import (
	//"log"
	"fmt"
	"reflect"
	"strings"
)

type Node interface {
	IsCoherent() error
	IsCoherentWith(n Node) error
	String() string
	IsOperator() bool
	AsKey() interface{}
}

type OR struct {
	Child Node
}

func (or *OR) GetChild() []Node {
	a, ok := or.Child.(*NArray);
	if ok {
		return a.Child
	}	
	return []Node{or.Child} 
}

func (or *OR) IsCoherent() error {
	children := or.GetChild()
	for _, child := range children {
		err := child.IsCoherent()
		if (err == nil) {
			//debugPrintf("OR IsCoherent %v : true\n", or)
			return nil
		}

	}
	//debugPrintf("OR IsCoherent %v : false\n", or)
	return fmt.Errorf("OR %v is not coherent", children)
}

func (or *OR) IsCoherentWith(n Node) error {
	children := or.GetChild()
	var err error
	for _, child := range children {
		err = child.IsCoherentWith(n)
		if (err == nil) {
			debugPrintf("OR IsCoherentWith %v    %v : true\n", or,n)
			return nil
		}
	}
	debugPrintf("OR IsCoherentWith %v    %v : false\n", or, n)
	return fmt.Errorf("OR is not coherent with : %v", err)
}

func (o *OR) String() string {
	var ret string
	for _, child := range o.GetChild() {
		ret += fmt.Sprintf("%v | ", child)
	}
	ret = strings.TrimSuffix(ret,"| ")
	return "(" + ret + ")"
}

func (*OR) IsOperator() bool {
	return true
}

func (o *OR) AsKey() interface{} {
	return o
}

type Coherent struct {
	Child Node
}

func (c *Coherent) GetChild() []Node {
	a, ok := c.Child.(*NArray);
	if ok {
		return a.Child
	}	
	return []Node{c.Child} 
}

func (c *Coherent) IsCoherent() error {
	children := c.GetChild()
	for _, child := range children {
		for _, child2 := range children {
			err := child.IsCoherentWith(child2) 
			if (err != nil) {
				return child2.IsCoherentWith(child)
			}
		}
	}	
	debugPrintf("Coherent IsCoherent %v : true\n", c)

	return nil
}

func (c *Coherent) IsCoherentWith(n Node) error {
	children := c.GetChild()
	var err error 
	for _, child := range children {
		err = child.IsCoherentWith(n) 
		if (err != nil) {
			return fmt.Errorf("%v is not coherent with %v : %v",c,n,err)
		}
	}
	debugPrintf("Coherent IsCoherentWith %v %v : true\n", c, n)
	return nil
}

func (c *Coherent) String() string {
	var ret string
	for _, child := range c.GetChild() {
		ret += fmt.Sprintf("%v&", child)
	}
	ret = strings.TrimSuffix(ret,"&")
	return "(" + ret + ")"
}

func (*Coherent) IsOperator() bool {
	return true
}

func (c *Coherent) AsKey() interface{} {
	return c
}

type Not struct {
	Child Node
}

func (n *Not) GetChild() []Node {
	a, ok := n.Child.(*NArray);
	if ok {
		return a.Child
	}	
	return []Node{n.Child} 
}

func (n *Not) IsCoherent() error {
	err := n.Child.IsCoherent()
	if (err != nil) {
		return nil
	}
	debugPrintf("Not IsCoherent %v : false\n", n)
	return fmt.Errorf("Not is not coherent '%v'",n.Child)
}

// This is a very special case.
// I mix classical logic and proposal construction.
// To use the same object we should split the 2 use case :
//  When the child are a part of a proposal (leaf, struct, array), we behave like "we don't want this to be that"
// (/A.B) became /(A.B)
// "The sky is not blue" // proposal
// vs "The sky is blue and not the weather is nice"

func (n *Not) IsCoherentWith(o Node) error {
	if n.Child.IsOperator() { // proposal
		err1 := n.IsCoherent()
		err2 := o.IsCoherent()	
		if (err1 == nil && err2 == nil) {
			debugPrintf("Not %v vs %v: true (proposal)\n", n, o)
			return nil
		}
	} else { // incomplete proposal
		err := n.Child.IsCoherentWith(o)
		if (err != nil) {
			debugPrintf("Not %v vs %v : true (part) %v\n", n, o , err)
			return nil
		}
	}
	debugPrintf("Not %v vs %v: false\n", n, o)
	return fmt.Errorf("Not, Both node should be different %v vs %v", n, o) 
}

func (n *Not) String() string {
	return fmt.Sprintf(" ~%v ", n.GetChild())
}

func (*Not) IsOperator() bool {
	return true
}

func (n *Not) AsKey() interface{} {
	return n
}

// it could be a leaf but string are common, special case could be handle
// 
//type Str string

var StrZero  Node = &Leaf{reflect.ValueOf("")}

//func (s Str) IsCoherent() error {
//	return nil
//}
//
//func (s Str) IsCoherentWith(n node) error {
//	s2, ok := n.(Str);
//	if (!ok) {
//		//case with OR/Not/Coherency/..
//		if n.IsOperator() {
//			return n.IsCoherentWith(s)
//		} else {
//			return fmt.Errorf("String shall be coherent")
//		}
//	}
//	if (s2 == s) {
//		return nil
//	} else if s == StrZero || s2 == StrZero {
//		return nil
//	}
//	
//	return fmt.Errorf("String shall be coherent")
//}
//
//
//func (s Str) String() string {
//	return fmt.Sprintf("'%s'",string(s))
//}
//
//func (Str) IsOperator() bool {
//	return false
//}

type Leaf struct {
	Value reflect.Value
}

func (l *Leaf) IsCoherent() error {
	return nil
}

func (l *Leaf) IsCoherentWith(n Node) error {
	l2, ok := n.(*Leaf);
	if (!ok) {
		//case with OR/Not/Coherency/.. in between
		if n.IsOperator() {			
			return n.IsCoherentWith(l)
		} else {
			return fmt.Errorf("Incoherent leaf %v vs %v", l, n)
		}
	}

	if (l2.Value.Interface() == l.Value.Interface()) {
		return nil
	}
	equalKind := l2.Value.Kind() == l.Value.Kind() 
	if (equalKind) {
		if(l.isNeutral() || l2.isNeutral()) {
			return nil
		}
	}
	
	return fmt.Errorf("Incoherent leaf %v vs %v (%s vs %s)", l, l2, l.Value.Kind(), l2.Value.Kind())
}

func (l *Leaf) isNeutral() bool {
	i := l.Value
	switch i.Kind() {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return i.Int() == -1 
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return i.Uint() == 1 
	case reflect.Float32, reflect.Float64:
		return i.Float() == 1.0 
	case reflect.String:
		return i.String() == "" 
	default: {
		fmt.Printf("isNeutral unknwon: %v\n",i)
		return false
	}
	}
	return false
}

func (l *Leaf) String() string {
	return fmt.Sprintf("%v", l.Value)
}

func (*Leaf) IsOperator() bool {
	return false
}

func (l *Leaf) AsKey() interface{} {
	return l.Value.Interface()
}

type NStruct struct {
	Child map[interface{}] NStructValue
}

type NStructValue struct {
	n Node
	key Node
}

func (n *NStructValue) String() string {
	return fmt.Sprintf("%v:%v ", n.key, n.n)
}


//todo : à complexifier par regexp possible
// un type string se retrouve dans un leaf
// un leaf contient un object reflect qui n'est "comparrable" qu'avec Interface()
// d'ou l'usage du AsKey
func (n *NStruct) get(k Node) Node {
	key := k.AsKey()
	value := n.Child[key]
	debugPrintf("Struct get[%v] %v\n", k, value.n)
	return value.n
}

func (n *NStruct) set(k Node, v Node) {
	key := k.AsKey()
	if (n.Child == nil) {
		n.Child = make(map[interface{}]NStructValue)
	}
	n.Child[key] = struct{n Node; key Node}{v,k}
}



func (n *NStruct) IsCoherent() error {
	for _,node := range n.Child {
		err := node.key.IsCoherent()
		if (err != nil) {
			return fmt.Errorf("Struct is not coherent, key is %v : %v",node.key, err)
		}
		err = node.n.IsCoherent()
		if (err != nil) {
			return fmt.Errorf("Struct is not coherent, [%v] valuse is %v: %v ",node.key, node.n, err)
		}
	}
	//debugPrintf("Struct IsCoherent %v : true\n", n)
	return nil
}
// une struct est coherent avec une autre si les champs présent sont coherent entre eux.
// coherent doit rester commutatif, on ne peut donc pas faire un coté inclus dans l'autre
// L'ordre des champs n'est pas important.
// il faudrait trouver un moyen d'inclusion strict si besoin.
// (ex: avec regexp, si plusieurs match les faire du plus particluiers au plus général, coherent: {{...},{ ..., "*": nil})
// si une clef est absente d'un coté, ce n'est pas grave.
func (n *NStruct) IsCoherentWith(n2 Node) error {
	s2, ok := n2.(*NStruct)
	if !ok {
		//return fmt.Errorf("Structure needed, %v vs %v\n", n, n2)
		if n2.IsOperator() {
			return n2.IsCoherentWith(n)
		} else {
			return fmt.Errorf("Struct %v is not coherent with %v ",n, n2)
		}
		debugPrintf("Struct")
	}
	for _, element := range n.Child {
		v2 := s2.get(element.key)
		if v2 == nil {
			continue
		}
		err := v2.IsCoherentWith(element.n)
		if err != nil {
			return fmt.Errorf("Struct %v is not coherent with %v : %v",v2, element.n, err)
		}
		debugPrintf("Struct")
	}
	debugPrintf("Struct IsCoherentWith %v %v : true\n", n, n2)
	return nil
}

func (n *NStruct) String() string {
	var ret string
	for _, element := range n.Child {
		ret+= fmt.Sprintf("%v:%v ", element.key, element.n)
	}
	return "{"+ret+"}"
}

func (* NStruct) IsOperator() bool {
	return false
}

func (n *NStruct) AsKey() interface{} {
	return n
}

type NArray struct {
	Child []Node
}

func (a *NArray) IsCoherent() error {
	for _,node := range a.Child {
		err := node.IsCoherent()
		if (err != nil) {
			return fmt.Errorf("Array is not coherent, %v : %v",node, err)
			
		}
	}
	debugPrintf("Array IsCoherent %v : true\n", a)
	return nil
}

// array is coherent with an other array
// Array are not ordered : so an element must be coherent with an other element in the other array, symmetricaly
// multiplicity are not defined, 
// 
func (a *NArray) IsCoherentWith(n2 Node) error {
	a2, ok := n2.(*NArray)
	if !ok {
		return n2.IsCoherentWith(a)
	}
	c  :=  a.Child
	c2 := a2.Child
	for _,k := range c {
		ok := false
		for _,k2 := range c2 {
			err := k2.IsCoherentWith(k)
			if (err == nil) {
				ok = true
				break
			}
		}
		if (!ok) {
			return fmt.Errorf("'Array' value should match without order :\n%v\n%v",a,n2)
		}

	}
	debugPrintf("Array IsCoherentWith %v %v : true\n",a,n2)
	return nil
}

func (n *NArray) String() string {
	var ret string
	for _,e := range n.Child {
		ret += fmt.Sprintf("%v,", e)
	}
	return "["+ret+"]"
}

func (* NArray) IsOperator() bool {
	return false
}

func (a *NArray) AsKey() interface{} {
	return a
}

var debug = false
func debugPrintf(format string, a ...interface{}) (n int, err error) {
	if debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}


