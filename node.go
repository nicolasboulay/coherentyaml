package main

import (
	//"errors"
	//"log"
	"fmt"
	"reflect"
	"github.com/pkg/errors"
)

type node interface {
	//GetChild() []node
	IsCoherent() error
	IsCoherentWith(n node) error
	String() string
	//New(interface{}) node
}

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
	children := or.GetChild()
	for _, child := range children {
		err := child.IsCoherent() 
		if (err != nil) {
			return errors.Wrap(err, "OR is not coherent")
		}
	}
	return nil
}

func (or *OR) IsCoherentWith(n node) error {
	children := or.GetChild()
	var err error 
	for _, child := range children {
		err = child.IsCoherentWith(n) 
		if (err == nil) {
			return nil
		}
	}
	for _, child := range children {
		err = n.IsCoherentWith(child) 
		if (err == nil) {
			return nil
		}
	}
	return errors.Wrap(err, "OR is not coherent with")
}

func (o *OR) String() string {
	return fmt.Sprintf("OR{%v}", o.GetChild())
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
	children := c.GetChild()
	for _, child := range children {
		for _, child2 := range children {
			err := child.IsCoherentWith(child2) 
			if (err != nil) {
				return child2.IsCoherentWith(child)
			}
		}
	}
	return nil
}

func (c *Coherent) IsCoherentWith(n node) error {
	//log.Printf("Coherent %v %v ?",c,n)
	children := c.GetChild()
	var err error 
	for _, child := range children {
		err = child.IsCoherentWith(n) 
		if (err != nil) {
			return errors.Wrap(err, "Coherent is not coherent")
		}
	}
	return nil
}

func (c *Coherent) String() string {
	return fmt.Sprintf("Coherent{%v}", c.GetChild())
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
	err := n.child.IsCoherent() 
	if (err != nil) {
		return nil
	}
	
	return errors.Wrap(err, "Not is not coherent")
}

func (n *Not) IsCoherentWith(o node) error {
	
	var err error
	//log.Printf("%v",n.child[0])
	//log.Printf("%v",o)
	err = n.child.IsCoherentWith(o)
	if (err != nil) {
		//log.Print(nil)
		return nil
	}
	//log.Print("not : Both node should be different")
	return fmt.Errorf("Not, Both node should be different %v vs %v", n, o) 
}

func (n *Not) String() string {
	return fmt.Sprintf("Not{%v}", n.GetChild())
}

type Str string

var StrZero Str = Str("")

//func (s Str) GetChild() []node {
//	return []node{}
//} 

func (s Str) IsCoherent() error {
	return nil
}

func (s Str) IsCoherentWith(n node) error {
	s2, ok := n.(Str);
	if (!ok) {
		//case with OR/Not/Coherency/.. in between
		return n.IsCoherentWith(s)
	}
	if (s2 == s) {
		return nil
	} else if s == StrZero || s2 == StrZero {
		return nil
	}
	
	return errors.New("String shall be coherent")
}


func (s Str) String() string {
	return fmt.Sprintf("Str{%s}",string(s))
}

type leaf struct {
	value reflect.Value
}

//func (l * leaf) GetChild() []node {
//	return []node{}
//}

func (l *leaf) IsCoherent() error {
	return nil
}

func (l *leaf) IsCoherentWith(n node) error {
	l2, ok := n.(*leaf);
	if (!ok) {
		//case with OR/Not/Coherency/.. in between
		return n.IsCoherentWith(l)
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

type nStruct struct {
	child map[node]node
}

//func (n *nStruct) GetChild() []node {
//	var v []node
//	for k := range n.value {
//		v = append(v, k)
//	}
//	return v
//}

//todo : à complexifier par regexp possible  
func (n *nStruct) get(k node) node {
	return n.child[k]
}

func (n *nStruct) IsCoherent() error {
	for k,node := range n.child {
		err := k.IsCoherent()
		if (err != nil) {
			return errors.Wrapf(err, "Struct is not coherent, key is %v",k)
		}
		err = node.IsCoherent()
		if (err != nil) {
			return errors.Wrapf(err, "Struct is not coherent, [%v] value is %v ", k, node)
		}
	}
	return nil
}
// une struct est coherent avec une autre si les champs présent sont coherent entre eux.
// coherent doit rester commutatif, on ne peut donc pas faire un coté inclus dans l'autre
// L'ordre des champs n'est pas important.
// il faudrait trouver un moyen d'inclusion strict si besoin.
// (ex: avec regexp, si plusieurs match les faire du plus particluiers au plus général, coherent: {{...},{ ..., "*": nil})
// si une clef est absente d'un coté, ce n'est pas grave.
func (n *nStruct) IsCoherentWith(n2 node) error {
	//fmt.Printf("Structure :\n %v vs\n %v\n", n, n2)
	s2, ok := n2.(*nStruct)
	if !ok {
		return fmt.Errorf("Structure needed :\n %v vs\n %v", n, n2)
	}
	for k, element := range n.child {
		v2 := s2.get(k)
		if v2 == nil {
			continue
		}
		err := v2.IsCoherentWith(element)
		if err != nil {
			return errors.Wrapf(err, "Struct %v is not coherent with %v",v2, element)
		}
	}
	return nil
}

func (n *nStruct) String() string {
	return fmt.Sprintf("nStruct{%v}", n.child)
}

type nArray struct {
	child []node
}

//func (a *nArray) GetChild() []node {
//	return a.child
//}

func (a *nArray) IsCoherent() error {
	//c := a.GetChild()
	for _,node := range a.child {
		err := node.IsCoherent()
		if (err != nil) {
			return errors.Wrapf(err, "Array is not coherent, %v",node)
		}
	}
	return nil
}

// array is coherent with an other array
// Array are not ordered : so an element must be coherent with an other element in the other array, symmetricaly
// multiplicity are not defined, 
// 
func (a *nArray) IsCoherentWith(n2 node) error {
	a2, ok := n2.(*nArray)
	if !ok {
		return fmt.Errorf("Array needed : %v vs %v", a, n2)
	}
	c  :=  a.child
	c2 := a2.child
	for _,k := range c {
		ok := false
		for _,k2 := range c2 {
			err := k2.IsCoherentWith(k)
			//log.Print(err)
			if (err == nil) {
				ok = true
				break
			}
		}
		if (!ok) {
			return fmt.Errorf("'Array' value should match without order :\n%v\n%v",a,n2)
		}

	}
	return nil
}

func (n *nArray) String() string {
	return fmt.Sprintf("[%v]", n.child)
}


