package main

import (
	//"log"
	"fmt"
	"reflect"
	"strings"
	"github.com/goccy/go-yaml"
)

func toYAMLString(root node) string {
	yamlString,_ := yaml.Marshal(root)
	return string(yamlString)
}
var debug = true

type node interface {
	IsCoherent() error
	IsCoherentWith(n node) error
	String() string
	IsOperator() bool
	AsKey() interface{}
	MarshalYAML() (interface{}, error)
}

type Yes struct {
}

func (*Yes) IsCoherent() error {
	return nil
}

func (*Yes) IsCoherentWith(n node) error {
	return nil
}
func (*Yes) String() string {
	return "true"
}
func (*Yes) IsOperator() bool {
	return true
}
func (y *Yes) AsKey() interface{} {
	return y
}
func (*Yes) MarshalYAML() (interface{}, error) {
	return "true", nil
}

var yes = &Yes{}

type No struct {
}

func (*No) IsCoherent() error {
	return fmt.Errorf("no is never coherent")
}

func (*No) IsCoherentWith(n node) error {
	return nil
}
func (*No) String() string {
	return "true"
}
func (*No) IsOperator() bool {
	return false
}
func (n *No) AsKey() interface{} {
	return n
}
func (*No) MarshalYAML() (interface{}, error) {
	return "false", nil
}

var no = &No{}

type OR struct {
	child node
}

func (or *OR) GetChild() []node {
	a, ok := or.child.(*nArray);
	if ok {
		return a.child
	}	
	return []node{or.child} 
}

func (or *OR) IsCoherent() error {
	debugPrintfStart("OR IsCoherent %v : true\n", or)
//	children := or.GetChild()
//	for _, child := range children {
//		err := child.IsCoherent()
//		if (err == nil) {
//			debugPrintfEnd("OR IsCoherent %v : true\n", or)
//			return nil
//		}
//
//	}	
//	debugPrintfEnd("OR IsCoherent %v : false\n", or)
//	return fmt.Errorf("OR %v is not coherent", children)
	ret := or.IsCoherentWith(yes)
	debugPrintfEnd("OR IsCoherent : %s\n", ret)
	return ret
}

func (or *OR) IsCoherentWith(n node) error {
	debugPrintfStart("OR IsCoherentWith %v    %v : true\n", or,n)
	children := or.GetChild()
	var err error
	for _, child := range children {
		err = child.IsCoherentWith(n)
		if (err == nil) {
			debugPrintf( "OR IsCoherentWith %v    %v : true\n", or,n)
			return nil
		}
	}
	debugPrintfEnd("OR IsCoherentWith %v    %v : false\n", or, n)
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

func bold(s string ) string {
	return "\033[1m*" + s + "\033[0m"
}

func keyS(n node, s string) string {
	d :=debug
	debug = false
	if (n.IsCoherent() == nil){
		s = bold(s)
	}
	debug = d
	return s
}

func (o *OR) MarshalYAML() (interface{}, error) {
	return yaml.MapSlice{
		{Key: keyS(o, "Or"), Value:o.child},
	}, nil
}

type Coherent struct {
	child node
}

func (c *Coherent) GetChild() []node {
	a, ok := c.child.(*nArray);
	if ok {
		return a.child
	}	
	return []node{c.child} 
}

func (c *Coherent) IsCoherent() error {
	debugPrintfStart("Coherent IsCoherent %v :\n", c)
//	children := c.GetChild()
//	for _, child := range children {
//		for _, child2 := range children {
//			err := child.IsCoherentWith(child2) 
//			if (err != nil) {
//				ret := child2.IsCoherentWith(child)
//				debugPrintfEnd("Coherent IsCoherent %v : %v\n", c, ret)
//				return ret
//			}
//		}
	//	}
	ret := c.IsCoherentWith(yes)
	debugPrintfEnd("Coherent IsCoherent : %v\n", ret)
	return ret
}

func (c *Coherent) IsCoherentWith(n node) error {
	debugPrintfStart("Coherent IsCoherentWith %v  %v :\n", c, n)
	children := c.GetChild()
	var err error 
	for _, child := range children {
		err = child.IsCoherentWith(n) 
		if (err != nil) {
			debugPrintfEnd("Coherent IsCoherentWith %v  %v : false\n", c, n)
			return fmt.Errorf("%v is not coherent with %v : %v",c,n,err)
		}
	}
	debugPrintfEnd("Coherent IsCoherentWith %v  %v : true\n", c, n)
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

func (c *Coherent) MarshalYAML() (interface{}, error) {
	return yaml.MapSlice{
		{Key: keyS(c,"Coherent"), Value:c.child},
	}, nil
}

type Not struct {
	child node
}

func (n *Not) GetChild() []node {
	a, ok := n.child.(*nArray);
	if ok {
		return a.child
	}	
	return []node{n.child} 
}

func (n *Not) IsCoherent() error {
	debugPrintfStart("Not IsCoherent %v :\n", n)
//	if n.child.IsOperator() { // proposal
//		err := n.child.IsCoherent()
//		if (err != nil) {
//			debugPrintfEnd("Not IsCoherent %v : true\n", n)
//			return nil
//		}
//		debugPrintfEnd("Not IsCoherent %v : false\n", n)
//		return fmt.Errorf("Not is not coherent '%v'",n.child)
//	}
//
//	// incomplete proposal : always true
//	debugPrintfEnd("Not IsCoherent %v : true\n", n)
//	return nil
	ret := n.IsCoherentWith(yes)
	debugPrintfEnd("not IsCoherent %v : %v\n", n ,ret)
	return ret
}

// This is a very special case.
// I mix classical logic and proposal construction.
// To use the same object we should split the 2 use case :
//  When the child are a part of a proposal (leaf, struct, array), we behave like "we don't want this to be that"
// (/A.B) became /(A.B)
// "The sky is not blue" // proposal
// vs "The sky is blue and not the weather is nice"

func (n *Not) IsCoherentWith(o node) error {
	debugPrintfStart("Not IsCoherentWith %v  %v:\n", n, o)
//	if n.child.IsOperator() { // proposal Not: Coherent:, Not: Or:
//		err1 := n.IsCoherent()
//		err2 := o.IsCoherent()	
//		if (err1 == nil && err2 == nil) {
//			debugPrintfEnd("Not IsCoherentWith %v  %v: true (proposal)\n", n, o)
//			return nil
//		}
//	} else { // incomplete proposal is always true
//		err := n.child.IsCoherentWith(o)
//		if (err != nil) {
//			debugPrintfEnd("Not IsCoherentWith %v  %v : true (part) %v\n", n, o , err)
//			return nil
//		}
//	}
//	debugPrintfEnd("Not IsCoherentWith %v  %v: false\n", n, o)
//	return fmt.Errorf("Not, Both node should be different %v vs %v", n, o) 

	b := o.IsCoherent()
	if b == nil {
		//b true
		err := n.child.IsCoherentWith(o)
		if (err != nil) {
			//a false
			debugPrintfEnd("Not IsCoherentWith %v  %v : true %v\n", n, o , err)
			return nil
		} else {
			//a true
			debugPrintfEnd("Not IsCoherentWith %v  %v : false %v\n", n, o , err)
			return fmt.Errorf("Not, %v is not coherent with %v\n",n,o)
		}		
	}
	// b false
	debugPrintfEnd("Not IsCoherentWith %v  %v : false\n", n, o)
	return fmt.Errorf("Not, %v is false : %v",o, b) 
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

func (n *Not) MarshalYAML() (interface{}, error) {
	return yaml.MapSlice{
		{Key: keyS(n,"Not"), Value:n.child},
	}, nil
}

// it could be a leaf but string are common, special case could be handle
// 
//type Str string

var StrZero  node = &leaf{reflect.ValueOf("")}

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

type leaf struct {
	value reflect.Value
}

func (l *leaf) IsCoherent() error {
	return nil
}

func (l *leaf) IsCoherentWith(n node) error {	
	l2, ok := n.(*leaf);
	if (!ok) {
		//case with OR/Not/Coherency/.. in between
		if n.IsOperator() {			
			return n.IsCoherentWith(l)
		} else {
			return fmt.Errorf("Incoherent leaf %v vs %v", l, n)
		}
	}

	if (l2.value.Interface() == l.value.Interface()) {
		return nil
	}
	equalKind := l2.value.Kind() == l.value.Kind() 
	if (equalKind) {
		if(l.isNeutral() || l2.isNeutral()) {
			return nil
		}
	}
	
	return fmt.Errorf("Incoherent leaf %v vs %v (%s vs %s)", l, l2, l.value.Kind(), l2.value.Kind())
}

func (l *leaf) isNeutral() bool {
	i := l.value
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

func (l *leaf) String() string {
	return fmt.Sprintf("%v", l.value)
}

func (*leaf) IsOperator() bool {
	return false
}

func (l *leaf) AsKey() interface{} {
	return l.value.Interface()
}

func (l *leaf) MarshalYAML() (interface{}, error) {
	return l.value.Interface(), nil
}

type nStruct struct {
	child map[interface{}] nStructValue
}

type nStructValue struct {
	n node
	key node
}

func (n *nStructValue) String() string {
	return fmt.Sprintf("%v:%v ", n.key, n.n)
}


//todo : à complexifier par regexp possible
// un type string se retrouve dans un leaf
// un leaf contient un object reflect qui n'est "comparrable" qu'avec Interface()
// d'ou l'usage du AsKey
func (n *nStruct) get(k node) node {
	key := k.AsKey()
	value := n.child[key]
	debugPrintfIn("Struct get[%v] %v\n", k, value.n)
	return value.n
}

func (n *nStruct) set(k node, v node) {
	key := k.AsKey()
	if (n.child == nil) {
		n.child = make(map[interface{}]nStructValue)
	}
	n.child[key] = struct{n node; key node}{v,k}
}

func (n *nStruct) IsCoherent() error {
	debugPrintfStart("Struct IsCoherent %v :\n", n)
	for _,node := range n.child {
		err := node.key.IsCoherent()
		if (err != nil) {
			debugPrintfEnd("Struct IsCoherent %v : false\n", n)
			return fmt.Errorf("Struct is not coherent, key is %v : %v",node.key, err)
		}
		err = node.n.IsCoherent()
		if (err != nil) {
			debugPrintfEnd("Struct IsCoherent %v : false\n", n)
			return fmt.Errorf("Struct is not coherent, [%v] valuse is %v: %v ",node.key, node.n, err)
		}
	}
	debugPrintfEnd("Struct IsCoherent %v : true\n", n)
	return nil
}
// une struct est coherent avec une autre si les champs présent sont coherent entre eux.
// coherent doit rester commutatif, on ne peut donc pas faire un coté inclus dans l'autre
// L'ordre des champs n'est pas important.
// il faudrait trouver un moyen d'inclusion strict si besoin.
// (ex: avec regexp, si plusieurs match les faire du plus particluiers au plus général, coherent: {{...},{ ..., "*": nil})
// si une clef est absente d'un coté, ce n'est pas grave.
func (n *nStruct) IsCoherentWith(n2 node) error {
	debugPrintfStart("Struct IsCoherentWith %v  %v :\n", n, n2)
	s2, ok := n2.(*nStruct)
	if !ok {
		//return fmt.Errorf("Structure needed, %v vs %v\n", n, n2)
		if n2.IsOperator() {
			ret := n2.IsCoherentWith(n)
			debugPrintfEnd("Struct IsCoherentWith %v  %v : %v\n", n, n2, ret)
			return ret
		} else {
			debugPrintfEnd("Struct IsCoherentWith %v  %v : false\n", n, n2)
			return fmt.Errorf("Struct %v is not coherent with %v ",n, n2)
		}
	}
	for _, element := range n.child {
		v2 := s2.get(element.key)
		if v2 == nil {
			continue
		}
		err := v2.IsCoherentWith(element.n)
		if err != nil {
			debugPrintfEnd("Struct IsCoherentWith %v  %v : false\n", n, n2)
			return fmt.Errorf("Struct %v is not coherent with %v : %v",v2, element.n, err)
		}
	}
	debugPrintfEnd("Struct IsCoherentWith %v  %v : true\n", n, n2)
	return nil
}

func (n *nStruct) String() string {
	var ret string
	for _, element := range n.child {
		ret+= fmt.Sprintf("%v:%v ", element.key, element.n)
	}
	return "{"+ret+"}"
}

func (* nStruct) IsOperator() bool {
	return false
}

func (n *nStruct) AsKey() interface{} {
	return n
}

func (n *nStruct) MarshalYAML() (interface{}, error) {
	var ret yaml.MapSlice
	for _, element := range n.child {
		ret = append(ret,yaml.MapItem{Key: keyS(n,element.key.String()), Value: element.n})
	}
	return ret, nil
}

type nArray struct {
	child []node
}

func (a *nArray) IsCoherent() error {
	debugPrintfStart("Array IsCoherent %v :\n", a)
	for _,node := range a.child {
		err := node.IsCoherent()
		if (err != nil) {
			debugPrintfEnd("Array IsCoherent %v : false\n", a)
			return fmt.Errorf("Array is not coherent, %v : %v",node, err)
			
		}
	}
	debugPrintfEnd("Array IsCoherent %v : true\n", a)
	return nil
}

// array is coherent with an other array
// Array are not ordered : so an element must be coherent with an other element in the other array, symmetricaly
// multiplicity are not defined, 
// 
func (a *nArray) IsCoherentWith(n2 node) error {
	debugPrintfStart("Array IsCoherentWith %v  %v : true\n",a,n2)
	a2, ok := n2.(*nArray)
	if !ok {
		return n2.IsCoherentWith(a)
	}
	c  :=  a.child
	c2 := a2.child
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
			debugPrintfEnd("Array IsCoherentWith %v  %v : false\n",a,n2)
			return fmt.Errorf("'Array' value should match without order :\n%v\n%v",a,n2)
		}
	}
	debugPrintfEnd("Array IsCoherentWith %v  %v : true\n",a,n2)
	return nil
}

func (n *nArray) String() string {
	var ret string
	for _,e := range n.child {
		ret += fmt.Sprintf("%v,", e)
	}
	return "["+ret+"]"
}

func (* nArray) IsOperator() bool {
	return false
}

func (a *nArray) AsKey() interface{} {
	return a
}

func (a *nArray) MarshalYAML() (interface{}, error) {
	return a.child, nil
}

func debugPrintf(format string, a ...interface{}) (n int, err error) {
	if debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

var debugSpaceNum = 0
func debugPrintfStart(format string, a ...interface{}) (n int, err error) {
	n,err = debugPrintfIn(format, a...)
	debugSpaceNum++
	return n, err
}
func debugPrintfIn(format string, a ...interface{}) (n int, err error) {
	n,err = debugPrintf(strings.Repeat(" ", debugSpaceNum) + format, a...)
	return n, err
}
func debugPrintfEnd(format string, a ...interface{}) (n int, err error) {
	debugSpaceNum--
	n,err = debugPrintfIn("/" + format, a...)
	return n, err
}



